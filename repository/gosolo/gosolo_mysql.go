package gosolo

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

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
			(name, phone, hotel_name, hotel_address, hotel_language_prompt, location_permission,
			 emergency_contact_name, emergency_contact_phone, hotel_whatsapp_number)
		VALUES
			(?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := m.Conn.ExecContext(ctx, query,
		input.Name,
		input.Phone,
		input.HotelName,
		input.HotelAddress,
		input.HotelLanguagePrompt,
		input.LocationPermission,
		input.EmergencyContactName,
		input.EmergencyContactPhone,
		input.HotelWhatsAppNumber,
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
		       location_permission, last_lat, last_lng, last_location_source, last_location_at, network_lost,
		       nudge_frequency_minutes, nudges_paused_until, nudges_disabled, last_nudge_at,
		       emergency_contact_name, emergency_contact_phone, hotel_whatsapp_number
		FROM travellers WHERE id=? 
	`

	row := m.Conn.QueryRowContext(ctx, query, travellerID)
	traveller := &models.Traveller{}
	var (
		locPerm               sql.NullInt64
		lastLat               sql.NullFloat64
		lastLng               sql.NullFloat64
		lastSource            sql.NullString
		lastLocTime           sql.NullTime
		networkLost           sql.NullInt64
		nudgeFreq             sql.NullInt64
		nudgesPausedUntil     sql.NullTime
		nudgesDisabled        sql.NullInt64
		lastNudgeAt           sql.NullTime
		emergencyContactName  sql.NullString
		emergencyContactPhone sql.NullString
		hotelWhatsApp         sql.NullString
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
		&nudgeFreq,
		&nudgesPausedUntil,
		&nudgesDisabled,
		&lastNudgeAt,
		&emergencyContactName,
		&emergencyContactPhone,
		&hotelWhatsApp,
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
	if nudgeFreq.Valid {
		val := int(nudgeFreq.Int64)
		traveller.NudgeFrequencyMinutes = &val
	}
	if nudgesPausedUntil.Valid {
		val := nudgesPausedUntil.Time
		traveller.NudgesPausedUntil = &val
	}
	if nudgesDisabled.Valid {
		traveller.NudgesDisabled = nudgesDisabled.Int64 == 1
	}
	if lastNudgeAt.Valid {
		val := lastNudgeAt.Time
		traveller.LastNudgeAt = &val
	}
	if emergencyContactName.Valid {
		val := emergencyContactName.String
		traveller.EmergencyContactName = &val
	}
	if emergencyContactPhone.Valid {
		val := emergencyContactPhone.String
		traveller.EmergencyContactPhone = &val
	}
	if hotelWhatsApp.Valid {
		val := hotelWhatsApp.String
		traveller.HotelWhatsAppNumber = &val
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

// UpdateNudgeSettings updates the nudge frequency preference
func (m *mysqlGoSoloRepo) UpdateNudgeSettings(ctx context.Context, travellerID int64, settings *models.NudgeSettingsInput) error {
	query := `UPDATE travellers SET nudge_frequency_minutes=? WHERE id=?`
	_, err := m.Conn.ExecContext(ctx, query, settings.FrequencyMinutes, travellerID)
	return err
}

// PauseNudges temporarily stops nudges for 1 hour or indefinitely
func (m *mysqlGoSoloRepo) PauseNudges(ctx context.Context, travellerID int64, duration string) error {
	var query string
	if duration == "1hour" {
		// Pause for 1 hour from now
		query = `UPDATE travellers SET nudges_paused_until=DATE_ADD(UTC_TIMESTAMP(), INTERVAL 1 HOUR), nudges_disabled=0 WHERE id=?`
	} else if duration == "indefinite" {
		// Pause indefinitely
		query = `UPDATE travellers SET nudges_disabled=1, nudges_paused_until=NULL WHERE id=?`
	} else {
		return errors.New("invalid duration: must be '1hour' or 'indefinite'")
	}
	_, err := m.Conn.ExecContext(ctx, query, travellerID)
	return err
}

// ResumeNudges re-enables nudges
func (m *mysqlGoSoloRepo) ResumeNudges(ctx context.Context, travellerID int64) error {
	query := `UPDATE travellers SET nudges_disabled=0, nudges_paused_until=NULL WHERE id=?`
	_, err := m.Conn.ExecContext(ctx, query, travellerID)
	return err
}

// ShouldSendSafetyNudge checks if a safety check nudge should be sent based on frequency
func (m *mysqlGoSoloRepo) ShouldSendSafetyNudge(ctx context.Context, travellerID int64) (bool, error) {
	traveller, err := m.GetTraveller(ctx, travellerID)
	if err != nil {
		return false, err
	}

	// If nudges are disabled or paused, don't send
	now := time.Now()
	if traveller.NudgesDisabled {
		return false, nil
	}
	if traveller.NudgesPausedUntil != nil && traveller.NudgesPausedUntil.After(now) {
		return false, nil
	}

	// If no frequency set (null), use conditional nudges (don't send safety check nudges)
	if traveller.NudgeFrequencyMinutes == nil {
		return false, nil
	}

	// If never sent a nudge, send one
	if traveller.LastNudgeAt == nil {
		return true, nil
	}

	// Check if enough time has passed
	frequencyMinutes := *traveller.NudgeFrequencyMinutes
	nextNudgeTime := traveller.LastNudgeAt.Add(time.Duration(frequencyMinutes) * time.Minute)
	return now.After(nextNudgeTime), nil
}

// MarkSafetyNudgeSent updates the last_nudge_at timestamp
func (m *mysqlGoSoloRepo) MarkSafetyNudgeSent(ctx context.Context, travellerID int64, nudgeID string) error {
	query := `UPDATE travellers SET last_nudge_at=UTC_TIMESTAMP() WHERE id=?`
	_, err := m.Conn.ExecContext(ctx, query, travellerID)
	return err
}

// RecordSafetyResponse stores a safety check response
func (m *mysqlGoSoloRepo) RecordSafetyResponse(ctx context.Context, travellerID int64, response *models.SafetyCheckResponseInput, lat, lng *float64) (*models.SafetyCheckResponse, error) {
	query := `
		INSERT INTO safety_check_responses 
			(traveller_id, nudge_id, response, response_time, response_lat, response_lng)
		VALUES (?, ?, ?, UTC_TIMESTAMP(), ?, ?)
	`
	result, err := m.Conn.ExecContext(ctx, query, travellerID, response.NudgeID, response.Response, lat, lng)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return &models.SafetyCheckResponse{
		ID:           id,
		TravellerID:  travellerID,
		NudgeID:      response.NudgeID,
		Response:     response.Response,
		ResponseTime: time.Now(),
		ResponseLat:  lat,
		ResponseLng:  lng,
	}, nil
}

// CreateEmergencyEscalation creates a new emergency escalation record
func (m *mysqlGoSoloRepo) CreateEmergencyEscalation(ctx context.Context, travellerID int64, responseID *int64, lat, lng float64, address string) (*models.EmergencyEscalation, error) {
	query := `
		INSERT INTO emergency_escalations 
			(traveller_id, response_id, escalated_at, location_lat, location_lng, location_address, status)
		VALUES (?, ?, UTC_TIMESTAMP(), ?, ?, ?, 'active')
	`
	result, err := m.Conn.ExecContext(ctx, query, travellerID, responseID, lat, lng, address)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return &models.EmergencyEscalation{
		ID:              id,
		TravellerID:     travellerID,
		ResponseID:      responseID,
		EscalatedAt:     time.Now(),
		LocationLat:     lat,
		LocationLng:     lng,
		LocationAddress: address,
		Status:          "active",
	}, nil
}

// UpdateEscalationStatus updates escalation notification status
func (m *mysqlGoSoloRepo) UpdateEscalationStatus(ctx context.Context, escalationID int64, status string, notifications map[string]bool) error {
	query := `
		UPDATE emergency_escalations 
		SET status=?, hotel_notified=?, emergency_contact_notified=?, police_notified=?
		WHERE id=?
	`
	hotelNotified := 0
	emergencyNotified := 0
	policeNotified := 0
	if notifications["hotel"] {
		hotelNotified = 1
	}
	if notifications["emergency_contact"] {
		emergencyNotified = 1
	}
	if notifications["police"] {
		policeNotified = 1
	}
	_, err := m.Conn.ExecContext(ctx, query, status, hotelNotified, emergencyNotified, policeNotified, escalationID)
	return err
}

// GetActiveEscalation gets the active escalation for a traveller
func (m *mysqlGoSoloRepo) GetActiveEscalation(ctx context.Context, travellerID int64) (*models.EmergencyEscalation, error) {
	query := `
		SELECT id, traveller_id, response_id, escalated_at, hotel_notified, emergency_contact_notified, 
		       police_notified, location_lat, location_lng, location_address, status, notes
		FROM emergency_escalations
		WHERE traveller_id=? AND status='active'
		ORDER BY escalated_at DESC
		LIMIT 1
	`
	row := m.Conn.QueryRowContext(ctx, query, travellerID)
	escalation := &models.EmergencyEscalation{}
	var responseID sql.NullInt64
	var locationAddress sql.NullString
	var notes sql.NullString

	err := row.Scan(
		&escalation.ID,
		&escalation.TravellerID,
		&responseID,
		&escalation.EscalatedAt,
		&escalation.HotelNotified,
		&escalation.EmergencyContactNotified,
		&escalation.PoliceNotified,
		&escalation.LocationLat,
		&escalation.LocationLng,
		&locationAddress,
		&escalation.Status,
		&notes,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if responseID.Valid {
		val := responseID.Int64
		escalation.ResponseID = &val
	}
	if locationAddress.Valid {
		escalation.LocationAddress = locationAddress.String
	}
	if notes.Valid {
		escalation.Notes = notes.String
	}

	return escalation, nil
}

// UpdateEmergencyContacts updates emergency contact information
func (m *mysqlGoSoloRepo) UpdateEmergencyContacts(ctx context.Context, travellerID int64, contacts *models.EmergencyContactUpdateInput) error {
	query := `
		UPDATE travellers 
		SET emergency_contact_name=COALESCE(?, emergency_contact_name),
		    emergency_contact_phone=COALESCE(?, emergency_contact_phone),
		    hotel_whatsapp_number=COALESCE(?, hotel_whatsapp_number)
		WHERE id=?
	`
	_, err := m.Conn.ExecContext(ctx, query,
		contacts.EmergencyContactName,
		contacts.EmergencyContactPhone,
		contacts.HotelWhatsAppNumber,
		travellerID,
	)
	return err
}
