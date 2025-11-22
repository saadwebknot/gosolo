# Update Emergency Contacts API

## Endpoint
**PUT** `/gosolo/travellers/{travellerID}/emergency-contacts`

## Purpose
Update emergency contact information for a registered traveller (name, phone, hotel WhatsApp number).

## Request Body
All fields are optional - you can update one, two, or all three:

```json
{
  "emergency_contact_name": "John Doe",
  "emergency_contact_phone": "+1-202-555-0123",
  "hotel_whatsapp_number": "+91-98765-43210"
}
```

## Response
**Success (200 OK):**
```json
{
  "message": "Emergency contacts updated successfully",
  "emergency_contact_name": "John Doe",
  "emergency_contact_phone": "+1-202-555-0123",
  "hotel_whatsapp_number": "+91-98765-43210"
}
```

## cURL Examples

### Update All Contacts
```bash
curl -X PUT http://localhost:8005/gosolo/travellers/1/emergency-contacts \
  -H "Content-Type: application/json" \
  -d '{
    "emergency_contact_name": "John Doe",
    "emergency_contact_phone": "+1-202-555-0123",
    "hotel_whatsapp_number": "+91-98765-43210"
  }'
```

### Update Only Emergency Contact Name
```bash
curl -X PUT http://localhost:8005/gosolo/travellers/1/emergency-contacts \
  -H "Content-Type: application/json" \
  -d '{
    "emergency_contact_name": "Jane Smith"
  }'
```

### Update Only Emergency Contact Phone
```bash
curl -X PUT http://localhost:8005/gosolo/travellers/1/emergency-contacts \
  -H "Content-Type: application/json" \
  -d '{
    "emergency_contact_phone": "+1-202-555-9999"
  }'
```

### Update Only Hotel WhatsApp Number
```bash
curl -X PUT http://localhost:8005/gosolo/travellers/1/emergency-contacts \
  -H "Content-Type: application/json" \
  -d '{
    "hotel_whatsapp_number": "+91-99999-88888"
  }'
```

### Update Multiple Fields
```bash
curl -X PUT http://localhost:8005/gosolo/travellers/1/emergency-contacts \
  -H "Content-Type: application/json" \
  -d '{
    "emergency_contact_name": "Alice Johnson",
    "hotel_whatsapp_number": "+91-11111-22222"
  }'
```

## Error Responses

**Invalid traveller ID (400):**
```json
{
  "message": "Invalid traveller id"
}
```

**No fields provided (400):**
```json
{
  "message": "At least one field must be provided"
}
```

**Traveller not found (404):**
```json
{
  "message": "Traveller not found"
}
```

## Notes

- Fields use `COALESCE` - only updates provided fields, leaves others unchanged
- Empty strings are treated as updates (to clear a field, send empty string)
- To set a field to NULL, send `null` in JSON

