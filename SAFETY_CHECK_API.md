# Safety Check Polling & Emergency Escalation API

## Overview
This MVP implements safety check polling where the app polls every 10-20 seconds, and if the user's set frequency is hit, they get a "Are you alright?" nudge. If they respond "yes, I'm safe" or "no, I'm in trouble", the system handles emergency escalation.

## Step 1: Run Database Migration

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
  `nudge_id` VARCHAR(100) DEFAULT NULL,
  `response` VARCHAR(20) NOT NULL COMMENT 'yes_safe or no_trouble',
  `response_time` DATETIME NOT NULL,
  `response_lat` DECIMAL(10,6) DEFAULT NULL,
  `response_lng` DECIMAL(10,6) DEFAULT NULL,
  `confirmed_escalation` TINYINT(1) DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_traveller_response` (`traveller_id`, `response_time`),
  CONSTRAINT `fk_safety_traveller` FOREIGN KEY (`traveller_id`) REFERENCES `travellers` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create emergency_escalations table
CREATE TABLE IF NOT EXISTS `emergency_escalations` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `traveller_id` BIGINT UNSIGNED NOT NULL,
  `response_id` BIGINT UNSIGNED DEFAULT NULL,
  `escalated_at` DATETIME NOT NULL,
  `hotel_notified` TINYINT(1) DEFAULT 0,
  `emergency_contact_notified` TINYINT(1) DEFAULT 0,
  `police_notified` TINYINT(1) DEFAULT 0,
  `location_lat` DECIMAL(10,6) NOT NULL,
  `location_lng` DECIMAL(10,6) NOT NULL,
  `location_address` VARCHAR(500) DEFAULT NULL,
  `status` VARCHAR(20) NOT NULL DEFAULT 'active',
  `notes` TEXT,
  PRIMARY KEY (`id`),
  KEY `idx_escalation_traveller` (`traveller_id`, `status`),
  CONSTRAINT `fk_escalation_traveller` FOREIGN KEY (`traveller_id`) REFERENCES `travellers` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## Step 2: Set Environment Variables

Add these to your `.env` file (and Render dashboard):

```bash
# Twilio WhatsApp API (for emergency alerts)
TWILIO_ACCOUNT_SID=your_account_sid
TWILIO_AUTH_TOKEN=your_auth_token
TWILIO_WHATSAPP_FROM=whatsapp:+14155238886
```

Get these from: https://console.twilio.com/

## Step 3: Register Traveller with Emergency Contacts

```bash
curl -X POST http://localhost:8005/gosolo/travellers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ava Solo",
    "phone": "+1-202-555-0199",
    "hotel_name": "Aurora Suites",
    "hotel_address": "221B MG Road",
    "hotel_language_prompt": "आप ऑरोरा सूट्स में ठहर रहे हैं।",
    "location_permission": true,
    "emergency_contact_name": "John Doe",
    "emergency_contact_phone": "+1-202-555-0123",
    "hotel_whatsapp_number": "+91-98765-43210"
  }'
```

## Step 4: Set Nudge Frequency

```bash
# Set frequency to 5 minutes (app will poll every 10-20 secs, but nudge only every 5 mins)
curl -X PUT http://localhost:8005/gosolo/travellers/1/nudge-settings \
  -H "Content-Type: application/json" \
  -d '{"frequency_minutes": 5}'
```

## Step 5: App Polling Flow

### App polls every 10-20 seconds:

```bash
# App polls this endpoint repeatedly
curl http://localhost:8005/gosolo/travellers/1/poll-safety-nudges
```

**Response when NO nudge needed (frequency not hit):**
```json
[]
```

**Response when nudge SHOULD be sent:**
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

## Step 6: User Responds to Safety Check

### User responds "Yes, I'm safe":

```bash
curl -X POST http://localhost:8005/gosolo/travellers/1/respond-safety-check \
  -H "Content-Type: application/json" \
  -d '{
    "nudge_id": "safety_check_1_1699123456",
    "response": "yes_safe"
  }'
```

**Response:**
```json
{
  "message": "Glad to know you're safe!",
  "response": "yes_safe",
  "response_id": 1
}
```

### User responds "No, I'm in trouble":

```bash
curl -X POST http://localhost:8005/gosolo/travellers/1/respond-safety-check \
  -H "Content-Type: application/json" \
  -d '{
    "nudge_id": "safety_check_1_1699123456",
    "response": "no_trouble"
  }'
```

**Response (asks for confirmation):**
```json
{
  "message": "We'll alert hotel staff, your emergency contact, and local police with your last recorded location. Do you want to continue?",
  "escalation_id": 1,
  "location": {
    "lat": 12.914599,
    "lng": 77.595036,
    "address": "221B MG Road"
  },
  "contacts_to_alert": {
    "hotel_whatsapp": true,
    "emergency_contact": true,
    "police_station": true
  },
  "requires_confirmation": true
}
```

## Step 7: Confirm Emergency Escalation

### User confirms "Yes, proceed":

```bash
curl -X POST http://localhost:8005/gosolo/travellers/1/confirm-escalation \
  -H "Content-Type: application/json" \
  -d '{
    "escalation_id": 1,
    "confirm": true
  }'
```

**Response:**
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
- Hotel WhatsApp number
- Emergency contact phone
- Police station (TODO: needs police API integration)

### User cancels (confirm: false):

```bash
curl -X POST http://localhost:8005/gosolo/travellers/1/confirm-escalation \
  -H "Content-Type: application/json" \
  -d '{
    "escalation_id": 1,
    "confirm": false
  }'
```

**Response:**
```json
{
  "message": "Emergency escalation cancelled"
}
```

## Complete Flow Summary

1. **App polls** `GET /gosolo/travellers/{id}/poll-safety-nudges` every 10-20 seconds
2. When frequency is hit, returns safety check nudge with `nudge_id`
3. User responds via `POST /gosolo/travellers/{id}/respond-safety-check`
   - `"yes_safe"` → Acknowledgement, done
   - `"no_trouble"` → Creates escalation, asks for confirmation
4. User confirms via `POST /gosolo/travellers/{id}/confirm-escalation`
   - `confirm: true` → Sends WhatsApp alerts to hotel, emergency contact, police
   - `confirm: false` → Cancels escalation

## Rule Engine Logic

- **Frequency check**: Only sends nudge if `nudge_frequency_minutes` has passed since `last_nudge_at`
- **Response timing**: Response must include the `nudge_id` from the poll response
- **Confirmation**: Emergency escalation requires explicit confirmation before sending alerts
- **Location**: Uses traveller's `last_lat`/`last_lng` for emergency alerts

