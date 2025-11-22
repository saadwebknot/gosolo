# Quick Reference - All cURL Commands

## 1. Update Emergency Contacts

**PUT** `/gosolo/travellers/{id}/emergency-contacts`

```bash
# Update all contacts
curl -X PUT http://localhost:8005/gosolo/travellers/1/emergency-contacts \
  -H "Content-Type: application/json" \
  -d '{
    "emergency_contact_name": "John Doe",
    "emergency_contact_phone": "+1-202-555-0123",
    "hotel_whatsapp_number": "+91-98765-43210"
  }'

# Update only name
curl -X PUT http://localhost:8005/gosolo/travellers/1/emergency-contacts \
  -H "Content-Type: application/json" \
  -d '{"emergency_contact_name": "Jane Smith"}'

# Update only phone
curl -X PUT http://localhost:8005/gosolo/travellers/1/emergency-contacts \
  -H "Content-Type: application/json" \
  -d '{"emergency_contact_phone": "+1-202-555-9999"}'

# Update only hotel WhatsApp
curl -X PUT http://localhost:8005/gosolo/travellers/1/emergency-contacts \
  -H "Content-Type: application/json" \
  -d '{"hotel_whatsapp_number": "+91-99999-88888"}'
```

## 2. Pause Nudges (Safety Checks Stop, Conditional Nudges Continue)

```bash
# Pause for 1 hour
curl -X POST http://localhost:8005/gosolo/travellers/1/pause-nudges \
  -H "Content-Type: application/json" \
  -d '{"duration": "1hour"}'

# Pause indefinitely
curl -X POST http://localhost:8005/gosolo/travellers/1/pause-nudges \
  -H "Content-Type: application/json" \
  -d '{"duration": "indefinite"}'
```

## 3. Resume Nudges

```bash
curl -X POST http://localhost:8005/gosolo/travellers/1/resume-nudges \
  -H "Content-Type: application/json"
```

## 4. Verify Pause Behavior

### Check Safety Check Nudges (Returns [] when paused)
```bash
curl http://localhost:8005/gosolo/travellers/1/poll-safety-nudges
```

### Check Conditional Nudges (Continue even when paused)
```bash
curl http://localhost:8005/gosolo/travellers/1/nudges
```

## 5. SOS Button (Immediate Emergency Alert)

```bash
curl -X POST http://localhost:8005/gosolo/travellers/1/sos-button \
  -H "Content-Type: application/json" \
  -d '{"notes":"Emergency SOS"}'
```

## Complete Test Flow

```bash
# 1. Update emergency contacts
curl -X PUT http://localhost:8005/gosolo/travellers/1/emergency-contacts \
  -H "Content-Type: application/json" \
  -d '{"emergency_contact_phone":"+1-202-555-0123"}'

# 2. Pause nudges for 1 hour
curl -X POST http://localhost:8005/gosolo/travellers/1/pause-nudges \
  -H "Content-Type: application/json" \
  -d '{"duration":"1hour"}'

# 3. Verify safety checks are paused (should return [])
curl http://localhost:8005/gosolo/travellers/1/poll-safety-nudges

# 4. Verify conditional nudges still work (should return area/weather alerts)
curl http://localhost:8005/gosolo/travellers/1/nudges

# 5. Resume nudges
curl -X POST http://localhost:8005/gosolo/travellers/1/resume-nudges \
  -H "Content-Type: application/json"
```

## Important Notes

✅ **Conditional nudges** (weather, crime, area, SOS) **ALWAYS work** - not affected by pause
✅ **Safety check nudges** ("Are you alright?") **STOP when paused**
✅ Update emergency contacts anytime with PUT endpoint
✅ Emergency contacts can be updated partially (only provide fields you want to change)

