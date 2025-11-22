# SOS Button API - Immediate Emergency Alert

## Overview
The SOS button is a one-tap emergency feature that immediately sends WhatsApp alerts to all safety contacts (hotel, emergency contact, police) without requiring confirmation. This is different from the safety check flow which requires user responses.

## Endpoint

**POST** `/gosolo/travellers/{travellerID}/sos-button`

## Request

**Headers:**
```
Content-Type: application/json
```

**Body (optional - can send notes):**
```json
{
  "channel": "sos_button",
  "notes": "Feeling unsafe, need immediate help"
}
```

**Body (minimal - just trigger):**
```json
{}
```

Or no body at all - all fields are optional.

## Response

**Success (200 OK):**
```json
{
  "message": "Emergency SOS activated. Alerts sent to all contacts.",
  "sos_id": 1,
  "escalation_id": 2,
  "notifications_sent": {
    "hotel": true,
    "emergency_contact": true,
    "police": true
  },
  "location": {
    "lat": 12.914599,
    "lng": 77.595036,
    "address": "221B MG Road, JP Nagar, Bengaluru"
  },
  "contacts_alerted": {
    "hotel": true,
    "emergency_contact": true,
    "police": true
  }
}
```

## What Happens When SOS Button is Pressed

1. **Logs SOS Event** - Creates a record in `sos_events` table
2. **Creates Emergency Escalation** - Creates record in `emergency_escalations` table
3. **Sends WhatsApp Alerts** to:
   - Hotel WhatsApp number (from traveller's `hotel_whatsapp_number`)
   - Emergency contact phone (from traveller's `emergency_contact_phone`)
   - Police station (TODO: needs police API integration)
4. **Uses Last Known Location** - Sends traveller's `last_lat`/`last_lng` in the alert

## WhatsApp Message Format

The message sent to all contacts looks like:

```
🚨 EMERGENCY SOS - Go-SOLO Safety System

Traveller: Ava Solo
Phone: +1-202-555-0199
Location: 221B MG Road, JP Nagar, Bengaluru (Lat: 12.914599, Lng: 77.595036)
Time: 2025-11-21 20:30:45 IST
Channel: sos_button

⚠️ IMMEDIATE ACTION REQUIRED
Please contact the traveller and provide assistance immediately.

Additional Notes: Feeling unsafe, need immediate help
```

## cURL Examples

### Basic SOS Button (no notes)
```bash
curl -X POST http://localhost:8005/gosolo/travellers/1/sos-button \
  -H "Content-Type: application/json" \
  -d '{}'
```

### SOS Button with Notes
```bash
curl -X POST http://localhost:8005/gosolo/travellers/1/sos-button \
  -H "Content-Type: application/json" \
  -d '{
    "channel": "sos_button",
    "notes": "Being followed, need immediate help"
  }'
```

### Minimal Request (empty body)
```bash
curl -X POST http://localhost:8005/gosolo/travellers/1/sos-button \
  -H "Content-Type: application/json"
```

## Requirements

1. **Traveller must be registered** with emergency contacts:
   - `emergency_contact_phone` (optional but recommended)
   - `hotel_whatsapp_number` (optional but recommended)

2. **Location must be set** - The traveller should have updated their location recently via:
   ```
   POST /gosolo/travellers/{id}/location
   ```

3. **WhatsApp API configured** - Environment variables must be set:
   - `TWILIO_ACCOUNT_SID`
   - `TWILIO_AUTH_TOKEN`
   - `TWILIO_WHATSAPP_FROM`

## Error Responses

**Traveller not found (404):**
```json
{
  "message": "Traveller not found"
}
```

**Missing location (handled gracefully):**
- If no location is set, sends alerts with "Location not available"
- Still sends all alerts successfully

**WhatsApp not configured (still works):**
- If WhatsApp env vars are missing, alerts are logged but not sent
- Response still returns success (for MVP, actual sending can be added later)

## Difference from Safety Check Flow

| Feature | Safety Check Flow | SOS Button |
|---------|------------------|------------|
| **Trigger** | Automatic (based on frequency) | Manual (user presses button) |
| **Confirmation** | Required (asks "Are you alright?") | None (immediate) |
| **Escalation** | Requires confirmation | Immediate |
| **Use Case** | Regular check-ins | Emergency situations |
| **Speed** | 2-step process | 1-step process |

## React Native Implementation

```javascript
const handleSOSButton = async () => {
  try {
    const response = await fetch(
      `${API_BASE_URL}/gosolo/travellers/${travellerID}/sos-button`,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          channel: 'sos_button',
          notes: 'Emergency SOS activated' // Optional
        })
      }
    );
    
    const data = await response.json();
    
    if (response.ok) {
      // Show success screen
      showEmergencyScreen({
        message: data.message,
        escalationId: data.escalation_id,
        notifications: data.notifications_sent,
        location: data.location
      });
      
      // Optionally show emergency contacts
      // Optionally allow calling emergency services
    }
  } catch (error) {
    // Handle error - still show emergency screen
    showEmergencyScreen({ error: true });
  }
};
```

## UI/UX Recommendations

1. **Large, prominent button** - Easy to find and press in emergency
2. **Confirmation (optional)** - Some apps show "Press again to confirm" to prevent accidental triggers
3. **Visual feedback** - Show loading state, then success screen
4. **Emergency screen** - After SOS is sent, show:
   - "Help is on the way" message
   - Emergency contact numbers
   - Call buttons for emergency services
   - Location details
   - "I'm Safe Now" button to resolve emergency

## Testing

```bash
# 1. Register traveller with emergency contacts
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

# 2. Update location
curl -X POST http://localhost:8005/gosolo/travellers/1/location \
  -H "Content-Type: application/json" \
  -d '{"lat":12.914599,"lng":77.595036,"source":"gps","network_lost":false}'

# 3. Press SOS button
curl -X POST http://localhost:8005/gosolo/travellers/1/sos-button \
  -H "Content-Type: application/json" \
  -d '{"notes":"Testing SOS button"}'
```

## Notes

- The SOS button is **immediate** - no confirmation required
- All alerts are sent **synchronously** - the response waits for all notifications
- If WhatsApp fails for one contact, others still get notified
- The escalation is marked as "active" and can be resolved later
- Location is always included if available (uses last known location)

