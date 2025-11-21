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
}
