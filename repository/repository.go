package repsitory

import (
	"context"

	"github.com/s1s1ty/go-mysql-crud/models"
)

// PostRepo explain...
type PostRepo interface {
	Fetch(ctx context.Context, num int64) ([]*models.Post, error)
	GetByID(ctx context.Context, id int64) (*models.Post, error)
	Create(ctx context.Context, p *models.Post) (int64, error)
	Update(ctx context.Context, p *models.Post) (*models.Post, error)
	Delete(ctx context.Context, id int64) (bool, error)
}

// TestRepo exposes read-only access to test records.
type TestRepo interface {
	FetchAll(ctx context.Context) ([]*models.Test, error)
}

// GoSoloRepo abstracts storage for the Go-SOLO module.
type GoSoloRepo interface {
	CreateTraveller(ctx context.Context, input *models.TravellerInput) (int64, error)
	UpdateLocation(ctx context.Context, travellerID int64, loc *models.LocationUpdate) error
	GetTraveller(ctx context.Context, travellerID int64) (*models.Traveller, error)
	LogSOS(ctx context.Context, travellerID int64, req *models.SOSRequest) (*models.SOSEvent, error)
	GetActiveSOSEvent(ctx context.Context, travellerID int64) (*models.SOSEvent, error)
	ListGuides(ctx context.Context) ([]*models.Guide, error)
	ListHazards(ctx context.Context) ([]*models.HazardZone, error)
	ListWeatherAlerts(ctx context.Context) ([]*models.WeatherAlert, error)
	// Nudge preferences management
	UpdateNudgeSettings(ctx context.Context, travellerID int64, settings *models.NudgeSettingsInput) error
	PauseNudges(ctx context.Context, travellerID int64, duration string) error
	ResumeNudges(ctx context.Context, travellerID int64) error
	// Safety check polling
	ShouldSendSafetyNudge(ctx context.Context, travellerID int64) (bool, error)
	MarkSafetyNudgeSent(ctx context.Context, travellerID int64, nudgeID string) error
	// Safety check responses
	RecordSafetyResponse(ctx context.Context, travellerID int64, response *models.SafetyCheckResponseInput, lat, lng *float64) (*models.SafetyCheckResponse, error)
	// Emergency escalations
	CreateEmergencyEscalation(ctx context.Context, travellerID int64, responseID *int64, lat, lng float64, address string) (*models.EmergencyEscalation, error)
	UpdateEscalationStatus(ctx context.Context, escalationID int64, status string, notifications map[string]bool) error
	GetActiveEscalation(ctx context.Context, travellerID int64) (*models.EmergencyEscalation, error)
	// Update emergency contacts
	UpdateEmergencyContacts(ctx context.Context, travellerID int64, contacts *models.EmergencyContactUpdateInput) error
}
