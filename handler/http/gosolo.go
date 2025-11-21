package handler

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/s1s1ty/go-mysql-crud/driver"
	"github.com/s1s1ty/go-mysql-crud/models"
	repository "github.com/s1s1ty/go-mysql-crud/repository"
	gosolorepo "github.com/s1s1ty/go-mysql-crud/repository/gosolo"
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
