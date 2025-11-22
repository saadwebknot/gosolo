# Prompt for React Native Development Tab

Copy and paste this entire prompt into your React Native Cursor tab:

---

**I'm building a React Native app for Go-SOLO (solo traveler safety app). The backend API is ready and I need you to help me implement the safety check polling, emergency escalation, and nudge management features. Here are the API endpoints:**

## Safety Check & Emergency APIs:

### 1. Poll for Safety Check (Poll every 10-20 seconds)
**GET** `/gosolo/travellers/{travellerID}/poll-safety-nudges`

Returns empty array `[]` if no nudge needed, or:
```json
[{
  "type": "safety_check",
  "message": "Are you alright?",
  "nudge_id": "safety_check_1_1699123456",
  "severity": "medium"
}]
```

### 2. Respond to Safety Check
**POST** `/gosolo/travellers/{travellerID}/respond-safety-check`
Body: `{"nudge_id": "safety_check_1_1699123456", "response": "yes_safe"}` or `"no_trouble"`

If "yes_safe": Returns `{"message": "Glad to know you're safe!", "response": "yes_safe", "response_id": 1}`

If "no_trouble": Returns escalation details asking for confirmation:
```json
{
  "message": "We'll alert hotel staff, your emergency contact, and local police... Do you want to continue?",
  "escalation_id": 1,
  "location": {"lat": 12.914599, "lng": 77.595036, "address": "..."},
  "contacts_to_alert": {"hotel_whatsapp": true, "emergency_contact": true, "police_station": true},
  "requires_confirmation": true
}
```

### 3. Confirm Emergency Escalation
**POST** `/gosolo/travellers/{travellerID}/confirm-escalation`
Body: `{"escalation_id": 1, "confirm": true}` (or `false` to cancel)

If confirm=true: Returns `{"message": "Emergency alerts sent successfully", "escalation_id": 1, "notifications_sent": {...}}`
If confirm=false: Returns `{"message": "Emergency escalation cancelled"}`

## Nudge Management APIs:

### 4. Set Nudge Frequency (Modal Window Setting)
**PUT** `/gosolo/travellers/{travellerID}/nudge-settings`
Body: `{"frequency_minutes": 5}` or `{"frequency_minutes": null}` for conditional nudges only

Returns:
```json
{
  "message": "Nudge settings updated",
  "frequency_minutes": 5,
  "note": "null or omitted = conditional nudges (default)"
}
```

**Implementation Notes:**
- Show this in a settings modal where user selects frequency (e.g., 5, 10, 15, 30 minutes)
- If user selects "Conditional Only" (no fixed frequency), send `{"frequency_minutes": null}`
- This controls how often safety check nudges appear (polling still happens every 10-20 seconds, but nudges only when frequency threshold is hit)

### 5. Pause Nudges (Snooze/Stop Indefinitely)
**POST** `/gosolo/travellers/{travellerID}/pause-nudges`
Body: `{"duration": "1hour"}` or `{"duration": "indefinite"}`

Returns for "1hour":
```json
{
  "message": "Nudges paused for 1 hour",
  "duration": "1hour"
}
```

Returns for "indefinite":
```json
{
  "message": "Nudges paused indefinitely. Use resume endpoint to re-enable.",
  "duration": "indefinite"
}
```

**Implementation Notes:**
- Show options: "Pause for 1 hour" and "Stop until I turn it back on"
- When paused, stop polling or continue polling but ignore nudges (backend handles this)
- Show status indicator when nudges are paused
- Auto-resume after 1 hour if duration was "1hour" (polling will automatically detect this)

### 6. Resume Nudges
**POST** `/gosolo/travellers/{travellerID}/resume-nudges`
No body needed

Returns:
```json
{
  "message": "Nudges resumed successfully"
}
```

**Implementation Notes:**
- Show this button when nudges are paused/disabled
- After resuming, immediately resume polling for safety check nudges
- Update UI to reflect that nudges are active again

### 7. Get All Nudges (Area/Weather/SOS conditional nudges)
**GET** `/gosolo/travellers/{travellerID}/nudges`

Returns array of nudges based on location and conditions:
```json
[
  {
    "type": "area",
    "message": "Avoid the poorly lit backlane near the market after 9pm...",
    "severity": "high"
  },
  {
    "type": "weather",
    "message": "Heavy rain predicted in 45 mins...",
    "severity": "medium"
  },
  {
    "type": "sos",
    "message": "SOS monitoring active. Tap 'I am safe' in the app to pause nudges.",
    "severity": "high"
  }
]
```

**Implementation Notes:**
- Call this separately from safety check polling (different endpoint)
- Show these as contextual notifications/banners in the app
- These are conditional (based on location, weather, SOS status), not time-based

## Complete User Flow:
1. App polls endpoint #1 every 10-20 seconds
2. When nudge appears, show modal: "Are you alright?" with buttons ["Yes, I'm Safe", "No, I'm in Trouble"]
3. User taps "Yes" → Call endpoint #2 with "yes_safe" → Show success → Dismiss modal → Resume polling
4. User taps "No" → Call endpoint #2 with "no_trouble" → Show confirmation modal with message and location → Buttons ["Yes, Send Alerts", "Cancel"]
5. User taps "Yes, Send Alerts" → Call endpoint #3 with confirm=true → Show success → Show emergency contacts
6. User taps "Cancel" → Call endpoint #3 with confirm=false → Dismiss modal → Resume polling

## Nudge Management Flow:

1. **User opens Settings** → Shows current frequency and pause status
2. **User sets frequency** → Call endpoint #4 → Update UI → Continue polling
3. **User pauses for 1 hour** → Call endpoint #5 with "1hour" → Stop polling OR continue polling but ignore responses → Show countdown timer → Auto-resume after 1 hour
4. **User stops indefinitely** → Call endpoint #5 with "indefinite" → Stop polling OR continue polling but ignore responses → Show "Nudges paused" indicator → User must manually resume
5. **User resumes** → Call endpoint #6 → Resume polling immediately → Update UI

**Important:** 
- If nudges are paused, the polling endpoint (#1) will still work but may return empty array or respect pause status
- Conditional nudges (endpoint #7) work independently - these show area/weather/SOS warnings regardless of pause status
- Safety check nudges (endpoint #1) respect pause settings - if paused, they won't appear

**Base URL**: `http://localhost:8005` (dev) or `https://gosolo-4.onrender.com` (prod)

**Please implement:**
- Polling mechanism (setInterval or polling library) for safety check nudges
- Modal components:
  - Safety check modal ("Are you alright?")
  - Emergency confirmation modal (after "no_trouble" response)
  - Nudge frequency settings modal (for user to set frequency)
  - Pause/Resume controls (1 hour or indefinite)
- API integration functions for all endpoints above
- State management for:
  - Nudge polling status
  - Current nudge (if any)
  - Escalation details
  - Nudge settings (frequency, pause status)
- Settings screen/component with:
  - Nudge frequency picker (5, 10, 15, 30 mins or "Conditional Only")
  - Pause buttons: "Pause for 1 hour" and "Stop until I turn back on"
  - Resume button (shown when paused)
- Error handling for all API calls
- Handle pause/resume states properly (don't poll when paused)

**UI/UX Considerations:**
- Show indicator when nudges are paused
- Show countdown timer when paused for 1 hour
- Make it easy to access settings from anywhere in the app
- Show a badge/notification count for conditional nudges (endpoint #7)
- Handle network errors gracefully (queue actions, retry when online)

Make it clean, handle edge cases, and follow React Native best practices.

---

