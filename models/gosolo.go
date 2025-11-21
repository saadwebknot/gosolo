package models

import "time"

// Traveller represents a solo traveller enrolled in Go-SOLO.
type Traveller struct {
	ID                  int64      `json:"id"`
	Name                string     `json:"name"`
	Phone               string     `json:"phone"`
	HotelName           string     `json:"hotel_name"`
	HotelAddress        string     `json:"hotel_address"`
	HotelLanguagePrompt string     `json:"hotel_language_prompt"`
	LocationPermission  bool       `json:"location_permission"`
	LastLat             *float64   `json:"last_lat,omitempty"`
	LastLng             *float64   `json:"last_lng,omitempty"`
	LastLocationSource  *string    `json:"last_location_source,omitempty"`
	LastLocationAt      *time.Time `json:"last_location_at,omitempty"`
	NetworkLost         bool       `json:"network_lost"`
	// Nudge preferences
	NudgeFrequencyMinutes *int       `json:"nudge_frequency_minutes,omitempty"` // null = use default/conditional
	NudgesPausedUntil     *time.Time `json:"nudges_paused_until,omitempty"`     // null = not paused
	NudgesDisabled        bool       `json:"nudges_disabled"`                   // true = paused indefinitely
	LastNudgeAt           *time.Time `json:"last_nudge_at,omitempty"`           // last safety check nudge sent
	// Emergency contacts
	EmergencyContactName  *string `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone *string `json:"emergency_contact_phone,omitempty"`
	HotelWhatsAppNumber   *string `json:"hotel_whatsapp_number,omitempty"`
}

// TravellerInput captures onboarding payload.
type TravellerInput struct {
	Name                  string  `json:"name"`
	Phone                 string  `json:"phone"`
	HotelName             string  `json:"hotel_name"`
	HotelAddress          string  `json:"hotel_address"`
	HotelLanguagePrompt   string  `json:"hotel_language_prompt"`
	LocationPermission    bool    `json:"location_permission"`
	EmergencyContactName  *string `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone *string `json:"emergency_contact_phone,omitempty"`
	HotelWhatsAppNumber   *string `json:"hotel_whatsapp_number,omitempty"`
}

// LocationUpdate registers the latest coordinates.
type LocationUpdate struct {
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Source      string  `json:"source"`
	NetworkLost bool    `json:"network_lost"`
}

// SOSRequest captures SOS trigger payload.
type SOSRequest struct {
	Channel string `json:"channel"`
	Notes   string `json:"notes"`
}

// SOSEvent stores SOS lifecycle.
type SOSEvent struct {
	ID          int64      `json:"id"`
	TravellerID int64      `json:"traveller_id"`
	Channel     string     `json:"channel"`
	Notes       string     `json:"notes"`
	Status      string     `json:"status"`
	TriggeredAt time.Time  `json:"triggered_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
}

// Nudge contains contextual notifications.
type Nudge struct {
	Type      string `json:"type"`
	Message   string `json:"message"`
	Severity  string `json:"severity"`
	ActionURL string `json:"action_url,omitempty"`
}

// Guide represents SOLO-BUDDY listing.
type Guide struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Language  string  `json:"language"`
	Contact   string  `json:"contact"`
	Rating    float64 `json:"rating"`
	RatePerHr float64 `json:"rate_per_hr"`
	Specialty string  `json:"specialty"`
}

// HazardZone describes unsafe area metadata.
type HazardZone struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	RadiusKM    float64 `json:"radius_km"`
	Severity    string  `json:"severity"`
	SafeRating  float64 `json:"safe_rating"`
	Description string  `json:"description"`
}

// WeatherAlert describes weather-based nudges.
type WeatherAlert struct {
	ID          int64   `json:"id"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	RadiusKM    float64 `json:"radius_km"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Severity    string  `json:"severity"`
}

// NudgeSettingsInput for updating nudge preferences
type NudgeSettingsInput struct {
	FrequencyMinutes *int `json:"frequency_minutes,omitempty"` // null or omit = use conditional nudges
}

// PauseNudgesInput for temporarily stopping nudges
type PauseNudgesInput struct {
	Duration string `json:"duration"` // "1hour" or "indefinite"
}

// SafetyCheckNudge represents a safety check poll response
type SafetyCheckNudge struct {
	Type      string `json:"type"`     // "safety_check"
	Message   string `json:"message"`  // "Are you alright?"
	NudgeID   string `json:"nudge_id"` // Unique identifier for this nudge
	Severity  string `json:"severity"` // "medium"
	ActionURL string `json:"action_url,omitempty"`
}

// SafetyCheckResponseInput for user responding to safety check
type SafetyCheckResponseInput struct {
	NudgeID  string `json:"nudge_id"` // ID of the nudge being responded to
	Response string `json:"response"` // "yes_safe" or "no_trouble"
}

// SafetyCheckResponse represents a stored response
type SafetyCheckResponse struct {
	ID                  int64     `json:"id"`
	TravellerID         int64     `json:"traveller_id"`
	NudgeID             string    `json:"nudge_id"`
	Response            string    `json:"response"`
	ResponseTime        time.Time `json:"response_time"`
	ResponseLat         *float64  `json:"response_lat,omitempty"`
	ResponseLng         *float64  `json:"response_lng,omitempty"`
	ConfirmedEscalation bool      `json:"confirmed_escalation"`
}

// EmergencyEscalation represents an emergency escalation event
type EmergencyEscalation struct {
	ID                       int64     `json:"id"`
	TravellerID              int64     `json:"traveller_id"`
	ResponseID               *int64    `json:"response_id,omitempty"`
	EscalatedAt              time.Time `json:"escalated_at"`
	HotelNotified            bool      `json:"hotel_notified"`
	EmergencyContactNotified bool      `json:"emergency_contact_notified"`
	PoliceNotified           bool      `json:"police_notified"`
	LocationLat              float64   `json:"location_lat"`
	LocationLng              float64   `json:"location_lng"`
	LocationAddress          string    `json:"location_address,omitempty"`
	Status                   string    `json:"status"` // "active", "resolved", "cancelled"
	Notes                    string    `json:"notes,omitempty"`
}

// EmergencyEscalationConfirmationInput for confirming emergency escalation
type EmergencyEscalationConfirmationInput struct {
	EscalationID int64 `json:"escalation_id"`
	Confirm      bool  `json:"confirm"` // true = yes, proceed with escalation
}
