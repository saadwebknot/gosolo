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
}

// TravellerInput captures onboarding payload.
type TravellerInput struct {
	Name                string `json:"name"`
	Phone               string `json:"phone"`
	HotelName           string `json:"hotel_name"`
	HotelAddress        string `json:"hotel_address"`
	HotelLanguagePrompt string `json:"hotel_language_prompt"`
	LocationPermission  bool   `json:"location_permission"`
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
