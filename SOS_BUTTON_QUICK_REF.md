# SOS Button API - Quick Reference

## Endpoint
**POST** `/gosolo/travellers/{travellerID}/sos-button`

## cURL Commands

### Basic SOS Button (Immediate Emergency Alert)
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

### For Production (Render)
```bash
curl -X POST https://gosolo-4.onrender.com/gosolo/travellers/1/sos-button \
  -H "Content-Type: application/json" \
  -d '{"notes":"Emergency SOS"}'
```

## Response
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

## What It Does
1. ✅ Logs SOS event
2. ✅ Creates emergency escalation
3. ✅ **Immediately sends WhatsApp alerts to:**
   - Hotel WhatsApp number
   - Emergency contact phone
   - Police station
4. ✅ Uses last known location
5. ✅ Returns success with all details

**No confirmation needed - alerts sent immediately!**

