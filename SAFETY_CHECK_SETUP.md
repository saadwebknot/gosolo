# Safety Check API - Setup & Testing Guide

## Step 1: Database Migration

Run this SQL in TiDB Cloud Chat2Query:

```sql
USE test;

-- Add emergency contact fields to travellers
ALTER TABLE `travellers` 
ADD COLUMN `emergency_contact_name` VARCHAR(255) DEFAULT NULL,
ADD COLUMN `emergency_contact_phone` VARCHAR(40) DEFAULT NULL,
ADD COLUMN `hotel_whatsapp_number` VARCHAR(40) DEFAULT NULL,
ADD COLUMN `last_nudge_at` DATETIME DEFAULT NULL COMMENT 'Last time a safety check nudge was sent';

-- Create safety_check_responses table
CREATE TABLE IF NOT EXISTS `safety_check_responses` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `traveller_id` BIGINT UNSIGNED NOT NULL,
  `nudge_id` VARCHAR(100) DEFAULT NULL COMMENT 'Identifier for the nudge being responded to',
  `response` VARCHAR(20) NOT NULL COMMENT 'yes_safe or no_trouble',
  `response_time` DATETIME NOT NULL,
  `response_lat` DECIMAL(10,6) DEFAULT NULL,
  `response_lng` DECIMAL(10,6) DEFAULT NULL,
  `confirmed_escalation` TINYINT(1) DEFAULT 0 COMMENT 'If user confirmed emergency escalation',
  PRIMARY KEY (`id`),
  KEY `idx_traveller_response` (`traveller_id`, `response_time`),
  CONSTRAINT `fk_safety_traveller` FOREIGN KEY (`traveller_id`) REFERENCES `travellers` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create emergency_escalations table
CREATE TABLE IF NOT EXISTS `emergency_escalations` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `traveller_id` BIGINT UNSIGNED NOT NULL,
  `response_id` BIGINT UNSIGNED DEFAULT NULL COMMENT 'Related safety_check_response',
  `escalated_at` DATETIME NOT NULL,
  `hotel_notified` TINYINT(1) DEFAULT 0,
  `emergency_contact_notified` TINYINT(1) DEFAULT 0,
  `police_notified` TINYINT(1) DEFAULT 0,
  `location_lat` DECIMAL(10,6) NOT NULL,
  `location_lng` DECIMAL(10,6) NOT NULL,
  `location_address` VARCHAR(500) DEFAULT NULL,
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT 'active, resolved, cancelled',
  `notes` TEXT,
  PRIMARY KEY (`id`),
  KEY `idx_escalation_traveller` (`traveller_id`, `status`),
  CONSTRAINT `fk_escalation_traveller` FOREIGN KEY (`traveller_id`) REFERENCES `travellers` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## Step 2: Set Environment Variables (for WhatsApp)

Add to `.env` file:
```bash
TWILIO_ACCOUNT_SID=your_account_sid_here
TWILIO_AUTH_TOKEN=your_auth_token_here
TWILIO_WHATSAPP_FROM=whatsapp:+14155238886
```

Or in Render dashboard, add:
- `TWILIO_ACCOUNT_SID`
- `TWILIO_AUTH_TOKEN`
- `TWILIO_WHATSAPP_FROM` (format: `whatsapp:+14155238886`)

## Step 3: Register Traveller with Emergency Contacts

```bash
curl -X POST http://localhost:8005/gosolo/travellers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ava Solo",
    "phone": "+1-202-555-0199",
    "hotel_name": "Aurora Suites",
    "hotel_address": "221B MG Road, JP Nagar, Bengaluru",
    "hotel_language_prompt": "आप ऑरोरा सूट्स में ठहर रहे हैं।",
    "location_permission": true,
    "emergency_contact_name": "John Doe",
    "emergency_contact_phone": "+1-202-555-0123",
    "hotel_whatsapp_number": "+91-98765-43210"
  }'
```

**Expected Response:**
```json
{
  "id": 1,
  "message": "Go-SOLO traveller registered",
  "go_solo": "active",
  "features": ["instant_sos", "conditional_nudges", "solo-buddy"]
}
```

## Step 4: Update Location (so we have last known location)

```bash
curl -X POST http://localhost:8005/gosolo/travellers/1/location \
  -H "Content-Type: application/json" \
  -d '{
    "lat": 12.914599,
    "lng": 77.595036,
    "source": "gps",
    "network_lost": false
  }'
```

**Expected Response:**
```json
{
  "message": "Location stored safely",
  "network_lost": false,
  "traveller_id": 1
}
```

## Step 5: Set Nudge Frequency (e.g., 5 minutes)

```bash
curl -X PUT http://localhost:8005/gosolo/travellers/1/nudge-settings \
  -H "Content-Type: application/json" \
  -d '{
    "frequency_minutes": 5
  }'
```

**Expected Response:**
```json
{
  "message": "Nudge settings updated",
  "frequency_minutes": 5,
  "note": "null or omitted = conditional nudges (default)"
}
```

## Step 6: Poll for Safety Check Nudges (App polls every 10-20 seconds)

```bash
# Poll 1: No nudge yet (frequency not hit or no frequency set)
curl http://localhost:8005/gosolo/travellers/1/poll-safety-nudges
```

**Expected Response (empty - no nudge needed yet):**
```json
[]
```

**After 5 minutes (if frequency was set to 5):**
```bash
# Poll 2: Nudge should appear
curl http://localhost:8005/gosolo/travellers/1/poll-safety-nudges
```

**Expected Response (nudge appears):**
```json
[
  {
    "type": "safety_check",
    "message": "Are you alright?",
    "nudge_id": "safety_check_1_1699123456",
    "severity": "medium"
  }
]
```

**Note:** App should poll this endpoint every 10-20 seconds. The server will return a nudge only when the frequency threshold is hit.

## Step 7: User Responds "Yes, I'm Safe"

```bash
curl -X POST http://localhost:8005/gosolo/travellers/1/respond-safety-check \
  -H "Content-Type: application/json" \
  -d '{
    "nudge_id": "safety_check_1_1699123456",
    "response": "yes_safe"
  }'
```

**Expected Response:**
```json
{
  "message": "Glad to know you're safe!",
  "response": "yes_safe",
  "response_id": 1
}
```

## Step 8: User Responds "No, I'm in Trouble" (Emergency Scenario)

```bash
curl -X POST http://localhost:8005/gosolo/travellers/1/respond-safety-check \
  -H "Content-Type: application/json" \
  -d '{
    "nudge_id": "safety_check_1_1699123456",
    "response": "no_trouble"
  }'
```

**Expected Response (asks for confirmation):**
```json
{
  "message": "We'll alert hotel staff, your emergency contact, and local police with your last recorded location. Do you want to continue?",
  "escalation_id": 1,
  "location": {
    "lat": 12.914599,
    "lng": 77.595036,
    "address": "221B MG Road, JP Nagar, Bengaluru"
  },
  "contacts_to_alert": {
    "hotel_whatsapp": true,
    "emergency_contact": true,
    "police_station": true
  },
  "requires_confirmation": true
}
```

## Step 9A: Confirm Emergency Escalation (Send Alerts)

```bash
curl -X POST http://localhost:8005/gosolo/travellers/1/confirm-escalation \
  -H "Content-Type: application/json" \
  -d '{
    "escalation_id": 1,
    "confirm": true
  }'
```

**Expected Response:**
```json
{
  "message": "Emergency alerts sent successfully",
  "escalation_id": 1,
  "notifications_sent": {
    "hotel": true,
    "emergency_contact": true,
    "police": true
  }
}
```

**WhatsApp messages are sent to:**
- Hotel WhatsApp: `+91-98765-43210`
- Emergency Contact: `+1-202-555-0123`
- Police Station: (TODO: needs police API integration)

## Step 9B: Cancel Emergency Escalation

```bash
curl -X POST http://localhost:8005/gosolo/travellers/1/confirm-escalation \
  -H "Content-Type: application/json" \
  -d '{
    "escalation_id": 1,
    "confirm": false
  }'
```

**Expected Response:**
```json
{
  "message": "Emergency escalation cancelled"
}
```

## Complete Test Flow Summary

1. **Register traveller** with emergency contacts
2. **Update location** so system knows where they are
3. **Set nudge frequency** (e.g., 5 minutes)
4. **Poll repeatedly** - App polls `GET /poll-safety-nudges` every 10-20 seconds
5. **Receive nudge** - After frequency threshold, get "Are you alright?" with `nudge_id`
6. **Respond**:
   - `"yes_safe"` → Done, marked safe
   - `"no_trouble"` → Creates escalation, asks for confirmation
7. **Confirm or Cancel**:
   - `confirm: true` → Sends WhatsApp alerts
   - `confirm: false` → Cancels escalation

## Testing Different Frequencies

```bash
# Test with 1 minute frequency (for quick testing)
curl -X PUT http://localhost:8005/gosolo/travellers/1/nudge-settings \
  -H "Content-Type: application/json" \
  -d '{"frequency_minutes": 1}'

# Then poll after 1 minute
curl http://localhost:8005/gosolo/travellers/1/poll-safety-nudges

# Test with conditional nudges (no frequency - won't send safety checks)
curl -X PUT http://localhost:8005/gosolo/travellers/1/nudge-settings \
  -H "Content-Type: application/json" \
  -d '{"frequency_minutes": null}'

# Poll will return empty array []
curl http://localhost:8005/gosolo/travellers/1/poll-safety-nudges
```

## Pause/Resume Nudges

```bash
# Pause for 1 hour
curl -X POST http://localhost:8005/gosolo/travellers/1/pause-nudges \
  -H "Content-Type: application/json" \
  -d '{"duration": "1hour"}'

# Pause indefinitely
curl -X POST http://localhost:8005/gosolo/travellers/1/pause-nudges \
  -H "Content-Type: application/json" \
  -d '{"duration": "indefinite"}'

# Resume nudges
curl -X POST http://localhost:8005/gosolo/travellers/1/resume-nudges \
  -H "Content-Type: application/json"
```

## For Production (Render.com)

Replace `http://localhost:8005` with your Render URL:
```bash
curl -X POST https://gosolo-4.onrender.com/gosolo/travellers/1/respond-safety-check \
  -H "Content-Type: application/json" \
  -d '{
    "nudge_id": "safety_check_1_1699123456",
    "response": "yes_safe"
  }'
```

