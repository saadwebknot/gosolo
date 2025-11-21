package gosolo

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/s1s1ty/go-mysql-crud/models"
	repository "github.com/s1s1ty/go-mysql-crud/repository"
)

func NewSQLGoSoloRepo(conn *sql.DB) repository.GoSoloRepo {
	return &mysqlGoSoloRepo{Conn: conn}
}

type mysqlGoSoloRepo struct {
	Conn *sql.DB
}

func (m *mysqlGoSoloRepo) CreateTraveller(ctx context.Context, input *models.TravellerInput) (int64, error) {
	query := `
		INSERT INTO travellers
			(name, phone, hotel_name, hotel_address, hotel_language_prompt, location_permission)
		VALUES
			(?, ?, ?, ?, ?, ?)
	`
	result, err := m.Conn.ExecContext(ctx, query,
		input.Name,
		input.Phone,
		input.HotelName,
		input.HotelAddress,
		input.HotelLanguagePrompt,
		input.LocationPermission,
	)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func (m *mysqlGoSoloRepo) UpdateLocation(ctx context.Context, travellerID int64, loc *models.LocationUpdate) error {
	query := `
		UPDATE travellers
			SET last_lat=?, last_lng=?, last_location_source=?, last_location_at=UTC_TIMESTAMP(), network_lost=?
		WHERE id=?
	`
	_, err := m.Conn.ExecContext(ctx, query, loc.Lat, loc.Lng, loc.Source, loc.NetworkLost, travellerID)
	return err
}

func (m *mysqlGoSoloRepo) GetTraveller(ctx context.Context, travellerID int64) (*models.Traveller, error) {
	query := `
		SELECT id, name, phone, hotel_name, hotel_address, hotel_language_prompt,
		       location_permission, last_lat, last_lng, last_location_source, last_location_at, network_lost
		FROM travellers WHERE id=? 
	`

	row := m.Conn.QueryRowContext(ctx, query, travellerID)
	traveller := &models.Traveller{}
	var (
		locPerm     sql.NullInt64
		lastLat     sql.NullFloat64
		lastLng     sql.NullFloat64
		lastSource  sql.NullString
		lastLocTime sql.NullTime
		networkLost sql.NullInt64
	)

	err := row.Scan(
		&traveller.ID,
		&traveller.Name,
		&traveller.Phone,
		&traveller.HotelName,
		&traveller.HotelAddress,
		&traveller.HotelLanguagePrompt,
		&locPerm,
		&lastLat,
		&lastLng,
		&lastSource,
		&lastLocTime,
		&networkLost,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		log.Printf("gosolo.GetTraveller scan error: %v", err)
		return nil, err
	}
	if locPerm.Valid {
		traveller.LocationPermission = locPerm.Int64 == 1
	}
	if lastLat.Valid {
		val := lastLat.Float64
		traveller.LastLat = &val
	}
	if lastLng.Valid {
		val := lastLng.Float64
		traveller.LastLng = &val
	}
	if lastSource.Valid {
		val := lastSource.String
		traveller.LastLocationSource = &val
	}
	if lastLocTime.Valid {
		val := lastLocTime.Time
		traveller.LastLocationAt = &val
	}
	if networkLost.Valid {
		traveller.NetworkLost = networkLost.Int64 == 1
	}
	return traveller, nil
}

func (m *mysqlGoSoloRepo) LogSOS(ctx context.Context, travellerID int64, req *models.SOSRequest) (*models.SOSEvent, error) {
	query := `
		INSERT INTO sos_events (traveller_id, channel, notes, status, triggered_at)
		VALUES (?, ?, ?, 'active', UTC_TIMESTAMP())
	`
	res, err := m.Conn.ExecContext(ctx, query, travellerID, req.Channel, req.Notes)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()

	return &models.SOSEvent{
		ID:          id,
		TravellerID: travellerID,
		Channel:     req.Channel,
		Notes:       req.Notes,
		Status:      "active",
	}, nil
}

func (m *mysqlGoSoloRepo) GetActiveSOSEvent(ctx context.Context, travellerID int64) (*models.SOSEvent, error) {
	query := `
		SELECT id, traveller_id, channel, notes, status, triggered_at, resolved_at
		FROM sos_events
		WHERE traveller_id=? AND status='active'
		ORDER BY triggered_at DESC
		LIMIT 1
	`
	row := m.Conn.QueryRowContext(ctx, query, travellerID)
	event := &models.SOSEvent{}

	err := row.Scan(
		&event.ID,
		&event.TravellerID,
		&event.Channel,
		&event.Notes,
		&event.Status,
		&event.TriggeredAt,
		&event.ResolvedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return event, nil
}

func (m *mysqlGoSoloRepo) ListGuides(ctx context.Context) ([]*models.Guide, error) {
	rows, err := m.Conn.QueryContext(ctx, `
		SELECT id, name, language, contact, rating, rate_per_hr, specialty
		FROM solo_guides ORDER BY rating DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var guides []*models.Guide
	for rows.Next() {
		g := new(models.Guide)
		if err := rows.Scan(&g.ID, &g.Name, &g.Language, &g.Contact, &g.Rating, &g.RatePerHr, &g.Specialty); err != nil {
			return nil, err
		}
		guides = append(guides, g)
	}
	return guides, rows.Err()
}

func (m *mysqlGoSoloRepo) ListHazards(ctx context.Context) ([]*models.HazardZone, error) {
	rows, err := m.Conn.QueryContext(ctx, `
		SELECT id, name, latitude, longitude, radius_km, severity, safe_rating, description
		FROM hazard_zones
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var zones []*models.HazardZone
	for rows.Next() {
		h := new(models.HazardZone)
		if err := rows.Scan(&h.ID, &h.Name, &h.Latitude, &h.Longitude, &h.RadiusKM, &h.Severity, &h.SafeRating, &h.Description); err != nil {
			return nil, err
		}
		zones = append(zones, h)
	}
	return zones, rows.Err()
}

func (m *mysqlGoSoloRepo) ListWeatherAlerts(ctx context.Context) ([]*models.WeatherAlert, error) {
	rows, err := m.Conn.QueryContext(ctx, `
		SELECT id, latitude, longitude, radius_km, type, description, severity
		FROM weather_alerts
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []*models.WeatherAlert
	for rows.Next() {
		w := new(models.WeatherAlert)
		if err := rows.Scan(&w.ID, &w.Latitude, &w.Longitude, &w.RadiusKM, &w.Type, &w.Description, &w.Severity); err != nil {
			return nil, err
		}
		alerts = append(alerts, w)
	}
	return alerts, rows.Err()
}
