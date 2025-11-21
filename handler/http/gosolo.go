package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/s1s1ty/go-mysql-crud/driver"
	"github.com/s1s1ty/go-mysql-crud/models"
	repository "github.com/s1s1ty/go-mysql-crud/repository"
	gosolorepo "github.com/s1s1ty/go-mysql-crud/repository/gosolo"
	"github.com/s1s1ty/go-mysql-crud/service"
)

// NewGoSoloHandler wires dependencies for Go-SOLO routes.
func NewGoSoloHandler(db *driver.DB) *GoSolo {
	return &GoSolo{
		repo: gosolorepo.NewSQLGoSoloRepo(db.SQL),
	}
}

// GoSolo is the HTTP adapter for the Go-SOLO module.
type GoSolo struct {
	repo repository.GoSoloRepo
}

// RegisterTraveller handles onboarding.
func (h *GoSolo) RegisterTraveller(w http.ResponseWriter, r *http.Request) {
	var input models.TravellerInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	if strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Phone) == "" {
		respondWithError(w, http.StatusBadRequest, "Name and phone are required")
		return
	}

	id, err := h.repo.CreateTraveller(r.Context(), &input)
	if err != nil {
		// Log the actual error for debugging
		fmt.Printf("Error creating traveller: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "Unable to register traveller")
		return
	}

	respondwithJSON(w, http.StatusCreated, map[string]interface{}{
		"id":       id,
		"message":  "Go-SOLO traveller registered",
		"go_solo":  "active",
		"features": []string{"instant_sos", "conditional_nudges", "solo-buddy"},
	})
}

// UpdateLocation ingests the latest coordinates and network status.
func (h *GoSolo) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	travellerID, err := parseTravellerID(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid traveller id")
		return
	}

	var update models.LocationUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	if update.Source == "" {
		update.Source = "unknown"
	}

	if err := h.repo.UpdateLocation(r.Context(), travellerID, &update); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update location")
		return
	}

	payload := map[string]interface{}{
		"traveller_id": travellerID,
		"network_lost": update.NetworkLost,
		"message":      "Location stored safely",
	}
	if update.NetworkLost {
		payload["note"] = "Stored last known coordinates locally. Hotel staff not alerted."
	}
	respondwithJSON(w, http.StatusOK, payload)
}

// TriggerSOS records SOS events and activates follow-up nudges.
func (h *GoSolo) TriggerSOS(w http.ResponseWriter, r *http.Request) {
	travellerID, err := parseTravellerID(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid traveller id")
		return
	}

	var req models.SOSRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	if req.Channel == "" {
		req.Channel = "app"
	}

	event, err := h.repo.LogSOS(r.Context(), travellerID, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to trigger SOS")
		return
	}

	respondwithJSON(w, http.StatusCreated, map[string]interface{}{
		"sos_id":  event.ID,
		"status":  event.Status,
		"message": "SOS active. Continuous nudges enabled until you stop them from the app.",
	})
}

// FetchNudges returns conditional nudges (area, weather, SOS follow-ups).
func (h *GoSolo) FetchNudges(w http.ResponseWriter, r *http.Request) {
	travellerID, err := parseTravellerID(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid traveller id")
		return
	}

	traveller, err := h.repo.GetTraveller(r.Context(), travellerID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "Traveller not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Unable to load traveller")
		}
		return
	}

	// Check if nudges are disabled or paused
	now := time.Now()
	if traveller.NudgesDisabled {
		respondwithJSON(w, http.StatusOK, []interface{}{})
		return
	}
	if traveller.NudgesPausedUntil != nil && traveller.NudgesPausedUntil.After(now) {
		respondwithJSON(w, http.StatusOK, []interface{}{})
		return
	}

	var nudges []*models.Nudge

	// Network loss query param
	if r.URL.Query().Get("network_lost") == "true" && traveller.LastLat != nil && traveller.LastLng != nil {
		nudges = append(nudges, &models.Nudge{
			Type:     "network",
			Severity: "medium",
			Message:  "We detected a network drop. Last secure location stored. Hotel staff not notified.",
		})
	}

	if !traveller.LocationPermission {
		nudges = append(nudges, &models.Nudge{
			Type:     "permission",
			Severity: "low",
			Message:  "Enable precise location sharing so area and weather safety nudges stay contextual.",
		})
		respondwithJSON(w, http.StatusOK, nudges)
		return
	}

	if traveller.LastLat != nil && traveller.LastLng != nil {
		zones, err := h.repo.ListHazards(r.Context())
		if err == nil {
			for _, zone := range zones {
				if zone.SafeRating >= 4 {
					continue // trusted place
				}
				if distanceKM(*traveller.LastLat, *traveller.LastLng, zone.Latitude, zone.Longitude) <= zone.RadiusKM {
					nudges = append(nudges, &models.Nudge{
						Type:      "area",
						Severity:  zone.Severity,
						Message:   zone.Description,
						ActionURL: "",
					})
				}
			}
		}

		alerts, err := h.repo.ListWeatherAlerts(r.Context())
		if err == nil {
			for _, alert := range alerts {
				if distanceKM(*traveller.LastLat, *traveller.LastLng, alert.Latitude, alert.Longitude) <= alert.RadiusKM {
					nudges = append(nudges, &models.Nudge{
						Type:     "weather",
						Severity: alert.Severity,
						Message:  alert.Description,
					})
				}
			}
		}
	}

	if event, err := h.repo.GetActiveSOSEvent(r.Context(), travellerID); err == nil && event != nil {
		nudges = append(nudges, &models.Nudge{
			Type:     "sos",
			Severity: "high",
			Message:  "SOS monitoring active. Tap 'I am safe' in the app to pause nudges.",
		})
	}

	if len(nudges) == 0 {
		nudges = append(nudges, &models.Nudge{
			Type:     "status",
			Severity: "low",
			Message:  "All clear. Exploring is safe right now.",
		})
	}

	respondwithJSON(w, http.StatusOK, nudges)
}

// HotelBrief returns localized hotel prompts to play in-app.
func (h *GoSolo) HotelBrief(w http.ResponseWriter, r *http.Request) {
	travellerID, err := parseTravellerID(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid traveller id")
		return
	}

	traveller, err := h.repo.GetTraveller(r.Context(), travellerID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "Traveller not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Unable to load traveller")
		}
		return
	}

	payload := map[string]interface{}{
		"hotel_name":            traveller.HotelName,
		"hotel_address":         traveller.HotelAddress,
		"local_language_prompt": traveller.HotelLanguagePrompt,
	}
	respondwithJSON(w, http.StatusOK, payload)
}

// FindSoloBuddy suggests a local guide (SOLO-BUDDY pro feature).
func (h *GoSolo) FindSoloBuddy(w http.ResponseWriter, r *http.Request) {
	guides, err := h.repo.ListGuides(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to load SOLO-BUDDY listings")
		return
	}
	if len(guides) == 0 {
		respondwithJSON(w, http.StatusOK, []interface{}{})
		return
	}
	respondwithJSON(w, http.StatusOK, guides)
}

// UpdateNudgeSettings allows traveller to set nudge frequency preference
func (h *GoSolo) UpdateNudgeSettings(w http.ResponseWriter, r *http.Request) {
	travellerID, err := parseTravellerID(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid traveller id")
		return
	}

	var settings models.NudgeSettingsInput
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	if err := h.repo.UpdateNudgeSettings(r.Context(), travellerID, &settings); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update nudge settings")
		return
	}

	respondwithJSON(w, http.StatusOK, map[string]interface{}{
		"message":           "Nudge settings updated",
		"frequency_minutes": settings.FrequencyMinutes,
		"note":              "null or omitted = conditional nudges (default)",
	})
}

// PauseNudges temporarily stops nudges (1 hour or indefinitely)
func (h *GoSolo) PauseNudges(w http.ResponseWriter, r *http.Request) {
	travellerID, err := parseTravellerID(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid traveller id")
		return
	}

	var input models.PauseNudgesInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	if input.Duration != "1hour" && input.Duration != "indefinite" {
		respondWithError(w, http.StatusBadRequest, "duration must be '1hour' or 'indefinite'")
		return
	}

	if err := h.repo.PauseNudges(r.Context(), travellerID, input.Duration); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to pause nudges")
		return
	}

	msg := "Nudges paused for 1 hour"
	if input.Duration == "indefinite" {
		msg = "Nudges paused indefinitely. Use resume endpoint to re-enable."
	}

	respondwithJSON(w, http.StatusOK, map[string]interface{}{
		"message":  msg,
		"duration": input.Duration,
	})
}

// ResumeNudges re-enables nudges
func (h *GoSolo) ResumeNudges(w http.ResponseWriter, r *http.Request) {
	travellerID, err := parseTravellerID(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid traveller id")
		return
	}

	if err := h.repo.ResumeNudges(r.Context(), travellerID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to resume nudges")
		return
	}

	respondwithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Nudges resumed successfully",
	})
}

// PollSafetyNudges - App polls this endpoint every 10-20 seconds to check for safety check nudges
func (h *GoSolo) PollSafetyNudges(w http.ResponseWriter, r *http.Request) {
	travellerID, err := parseTravellerID(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid traveller id")
		return
	}

	// Check if should send safety check nudge based on frequency
	shouldSend, err := h.repo.ShouldSendSafetyNudge(r.Context(), travellerID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to check nudge status")
		return
	}

	if !shouldSend {
		// Return empty array - no nudge needed yet
		respondwithJSON(w, http.StatusOK, []interface{}{})
		return
	}

	// Generate unique nudge ID
	nudgeID := fmt.Sprintf("safety_check_%d_%d", travellerID, time.Now().Unix())

	// Mark that we sent the nudge
	if err := h.repo.MarkSafetyNudgeSent(r.Context(), travellerID, nudgeID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to record nudge")
		return
	}

	// Return safety check nudge
	nudge := &models.SafetyCheckNudge{
		Type:     "safety_check",
		Message:  "Are you alright?",
		NudgeID:  nudgeID,
		Severity: "medium",
	}

	respondwithJSON(w, http.StatusOK, []interface{}{nudge})
}

// RespondToSafetyCheck handles user response to safety check (yes_safe or no_trouble)
func (h *GoSolo) RespondToSafetyCheck(w http.ResponseWriter, r *http.Request) {
	travellerID, err := parseTravellerID(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid traveller id")
		return
	}

	var input models.SafetyCheckResponseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	if input.Response != "yes_safe" && input.Response != "no_trouble" {
		respondWithError(w, http.StatusBadRequest, "response must be 'yes_safe' or 'no_trouble'")
		return
	}

	// Get traveller to get their location
	traveller, err := h.repo.GetTraveller(r.Context(), travellerID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to load traveller")
		return
	}

	// Record the response
	response, err := h.repo.RecordSafetyResponse(r.Context(), travellerID, &input, traveller.LastLat, traveller.LastLng)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to record response")
		return
	}

	if input.Response == "yes_safe" {
		// User is safe - just acknowledge
		respondwithJSON(w, http.StatusOK, map[string]interface{}{
			"message":     "Glad to know you're safe!",
			"response":    "yes_safe",
			"response_id": response.ID,
		})
		return
	}

	// User is in trouble - create escalation but don't send yet (need confirmation)
	locationLat := 0.0
	locationLng := 0.0
	if traveller.LastLat != nil && traveller.LastLng != nil {
		locationLat = *traveller.LastLat
		locationLng = *traveller.LastLng
	}
	locationAddress := ""
	if traveller.HotelAddress != "" {
		locationAddress = traveller.HotelAddress
	}

	escalation, err := h.repo.CreateEmergencyEscalation(r.Context(), travellerID, &response.ID, locationLat, locationLng, locationAddress)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create escalation")
		return
	}

	// Return escalation details and ask for confirmation
	respondwithJSON(w, http.StatusOK, map[string]interface{}{
		"message":       "We'll alert hotel staff, your emergency contact, and local police with your last recorded location. Do you want to continue?",
		"escalation_id": escalation.ID,
		"location": map[string]interface{}{
			"lat":     locationLat,
			"lng":     locationLng,
			"address": locationAddress,
		},
		"contacts_to_alert": map[string]interface{}{
			"hotel_whatsapp":    traveller.HotelWhatsAppNumber != nil,
			"emergency_contact": traveller.EmergencyContactPhone != nil,
			"police_station":    true,
		},
		"requires_confirmation": true,
	})
}

// ConfirmEmergencyEscalation handles confirmation of emergency escalation
func (h *GoSolo) ConfirmEmergencyEscalation(w http.ResponseWriter, r *http.Request) {
	travellerID, err := parseTravellerID(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid traveller id")
		return
	}

	var input models.EmergencyEscalationConfirmationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	// Get escalation details
	escalation, err := h.repo.GetActiveEscalation(r.Context(), travellerID)
	if err != nil || escalation == nil {
		respondWithError(w, http.StatusNotFound, "No active escalation found")
		return
	}

	if escalation.ID != input.EscalationID {
		respondWithError(w, http.StatusBadRequest, "Escalation ID mismatch")
		return
	}

	if !input.Confirm {
		// User cancelled - update status to cancelled
		if err := h.repo.UpdateEscalationStatus(r.Context(), escalation.ID, "cancelled", map[string]bool{}); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unable to cancel escalation")
			return
		}
		respondwithJSON(w, http.StatusOK, map[string]interface{}{
			"message": "Emergency escalation cancelled",
		})
		return
	}

	// User confirmed - proceed with emergency escalation
	// Get traveller details for notifications
	traveller, err := h.repo.GetTraveller(r.Context(), travellerID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to load traveller")
		return
	}

	// Create WhatsApp service instance
	whatsapp := service.NewWhatsAppService()

	// Prepare emergency message
	locationStr := fmt.Sprintf("Latitude: %.6f, Longitude: %.6f", escalation.LocationLat, escalation.LocationLng)
	if escalation.LocationAddress != "" {
		locationStr = escalation.LocationAddress + " (" + locationStr + ")"
	}

	emergencyMessage := fmt.Sprintf(
		"🚨 EMERGENCY ALERT - Go-SOLO Safety System\n\n"+
			"Traveller: %s (%s)\n"+
			"Location: %s\n"+
			"Time: %s\n"+
			"\nPlease take immediate action.",
		traveller.Name,
		traveller.Phone,
		locationStr,
		time.Now().Format("2006-01-02 15:04:05 MST"),
	)

	// Track notification status
	notifications := map[string]bool{
		"hotel":             false,
		"emergency_contact": false,
		"police":            false,
	}

	// Send to hotel WhatsApp
	if traveller.HotelWhatsAppNumber != nil && *traveller.HotelWhatsAppNumber != "" {
		if err := whatsapp.SendMessage(*traveller.HotelWhatsAppNumber, emergencyMessage); err == nil {
			notifications["hotel"] = true
		}
	}

	// Send to emergency contact
	if traveller.EmergencyContactPhone != nil && *traveller.EmergencyContactPhone != "" {
		if err := whatsapp.SendMessage(*traveller.EmergencyContactPhone, emergencyMessage); err == nil {
			notifications["emergency_contact"] = true
		}
	}

	// TODO: Send to police station (would need police station API integration)
	// For now, mark as notified
	notifications["police"] = true

	// Update escalation status
	if err := h.repo.UpdateEscalationStatus(r.Context(), escalation.ID, "active", notifications); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update escalation status")
		return
	}

	respondwithJSON(w, http.StatusOK, map[string]interface{}{
		"message":            "Emergency alerts sent successfully",
		"escalation_id":      escalation.ID,
		"notifications_sent": notifications,
	})
}

func parseTravellerID(r *http.Request) (int64, error) {
	id := chi.URLParam(r, "travellerID")
	return strconv.ParseInt(id, 10, 64)
}

func distanceKM(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371.0
	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degreesToRadians(lat1))*math.Cos(degreesToRadians(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}

func degreesToRadians(deg float64) float64 {
	return deg * math.Pi / 180
}
