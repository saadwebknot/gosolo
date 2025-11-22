# React Native App - Go-SOLO Safety Check API Integration Guide

## Overview
I'm building a React Native app for Go-SOLO, a solo traveler safety companion. The backend API is ready, and I need you to help me implement the frontend integration for the safety check polling and emergency escalation features.

## Base URL
- **Development**: `http://localhost:8005`
- **Production**: `https://gosolo-4.onrender.com`

## API Endpoints Overview

### 1. Poll for Safety Check Nudges (Main Polling Endpoint)
The app should poll this endpoint every 10-20 seconds to check if a safety check nudge should be displayed.

**Endpoint**: `GET /gosolo/travellers/{travellerID}/poll-safety-nudges`

**Request**:
```javascript
GET /gosolo/travellers/1/poll-safety-nudges
```

**Response - No Nudge (frequency not hit yet)**:
```json
[]
```

**Response - Nudge Should Be Shown**:
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

**Implementation Notes**:
- Poll every 10-20 seconds using `setInterval` or a polling library
- When response is empty array `[]`, do nothing
- When response contains a nudge object, show a modal/popup with:
  - Message: "Are you alright?"
  - Two buttons: "Yes, I'm Safe" and "No, I'm in Trouble"
  - Store the `nudge_id` from the response for the response API call

**cURL Example**:
```bash
curl http://localhost:8005/gosolo/travellers/1/poll-safety-nudges
```

---

### 2. Respond to Safety Check
When user taps "Yes, I'm Safe" or "No, I'm in Trouble", call this endpoint.

**Endpoint**: `POST /gosolo/travellers/{travellerID}/respond-safety-check`

**Request Body**:
```json
{
  "nudge_id": "safety_check_1_1699123456",
  "response": "yes_safe"  // or "no_trouble"
}
```

**Response for "yes_safe"**:
```json
{
  "message": "Glad to know you're safe!",
  "response": "yes_safe",
  "response_id": 1
}
```

**Response for "no_trouble"**:
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

**Implementation Notes**:
- If response is `"yes_safe"`: Show success message, dismiss modal, continue polling
- If response is `"no_trouble"`: Show confirmation modal with:
  - The message from `response.message`
  - Show location details
  - Show contacts that will be alerted
  - Two buttons: "Yes, Send Alerts" and "Cancel"
  - Store `escalation_id` for the confirmation API call

**cURL Examples**:

```bash
# Yes, I'm Safe
curl -X POST http://localhost:8005/gosolo/travellers/1/respond-safety-check \
  -H "Content-Type: application/json" \
  -d '{
    "nudge_id": "safety_check_1_1699123456",
    "response": "yes_safe"
  }'

# No, I'm in Trouble
curl -X POST http://localhost:8005/gosolo/travellers/1/respond-safety-check \
  -H "Content-Type: application/json" \
  -d '{
    "nudge_id": "safety_check_1_1699123456",
    "response": "no_trouble"
  }'
```

---

### 3. Confirm Emergency Escalation (Accept/Cancel)
After user responds "No, I'm in Trouble", they need to confirm or cancel the emergency escalation.

**Endpoint**: `POST /gosolo/travellers/{travellerID}/confirm-escalation`

**Request Body**:
```json
{
  "escalation_id": 1,
  "confirm": true  // true = send alerts, false = cancel
}
```

**Response for `confirm: true` (Send Alerts)**:
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

**Response for `confirm: false` (Cancel)**:
```json
{
  "message": "Emergency escalation cancelled"
}
```

**Implementation Notes**:
- When user taps "Yes, Send Alerts": Send `confirm: true`
- When user taps "Cancel": Send `confirm: false`
- Show appropriate success/cancellation message
- After successful escalation, you may want to:
  - Show emergency contacts on screen
  - Allow user to call emergency services
  - Show "I'm Safe Now" button to resolve the emergency

**cURL Examples**:

```bash
# Confirm - Send Emergency Alerts
curl -X POST http://localhost:8005/gosolo/travellers/1/confirm-escalation \
  -H "Content-Type: application/json" \
  -d '{
    "escalation_id": 1,
    "confirm": true
  }'

# Cancel - Don't Send Alerts
curl -X POST http://localhost:8005/gosolo/travellers/1/confirm-escalation \
  -H "Content-Type: application/json" \
  -d '{
    "escalation_id": 1,
    "confirm": false
  }'
```

---

## Additional Supporting Endpoints

### 4. Set Nudge Frequency
Set how often safety check nudges should appear (in minutes).

**Endpoint**: `PUT /gosolo/travellers/{travellerID}/nudge-settings`

**Request Body**:
```json
{
  "frequency_minutes": 5  // null = conditional nudges only (no safety checks)
}
```

**Response**:
```json
{
  "message": "Nudge settings updated",
  "frequency_minutes": 5,
  "note": "null or omitted = conditional nudges (default)"
}
```

**cURL**:
```bash
curl -X PUT http://localhost:8005/gosolo/travellers/1/nudge-settings \
  -H "Content-Type: application/json" \
  -d '{"frequency_minutes": 5}'
```

---

### 5. Pause Nudges
Temporarily stop safety check nudges.

**Endpoint**: `POST /gosolo/travellers/{travellerID}/pause-nudges`

**Request Body**:
```json
{
  "duration": "1hour"  // or "indefinite"
}
```

**Response**:
```json
{
  "message": "Nudges paused for 1 hour",
  "duration": "1hour"
}
```

**cURL**:
```bash
curl -X POST http://localhost:8005/gosolo/travellers/1/pause-nudges \
  -H "Content-Type: application/json" \
  -d '{"duration": "1hour"}'
```

---

### 6. Resume Nudges
Re-enable safety check nudges.

**Endpoint**: `POST /gosolo/travellers/{travellerID}/resume-nudges`

**Request**: No body needed

**Response**:
```json
{
  "message": "Nudges resumed successfully"
}
```

**cURL**:
```bash
curl -X POST http://localhost:8005/gosolo/travellers/1/resume-nudges \
  -H "Content-Type: application/json"
```

---

## Complete User Flow

1. **App starts** → Start polling `GET /poll-safety-nudges` every 10-20 seconds

2. **Poll returns nudge** → Show modal:
   ```
   "Are you alright?"
   [Yes, I'm Safe]  [No, I'm in Trouble]
   ```

3. **User taps "Yes, I'm Safe"**:
   - Call `POST /respond-safety-check` with `response: "yes_safe"`
   - Show success message
   - Dismiss modal
   - Continue polling

4. **User taps "No, I'm in Trouble"**:
   - Call `POST /respond-safety-check` with `response: "no_trouble"`
   - Show confirmation modal:
     ```
     "We'll alert hotel staff, your emergency contact, and local police 
      with your last recorded location. Do you want to continue?"
     
     Location: 221B MG Road, JP Nagar, Bengaluru
     
     [Yes, Send Alerts]  [Cancel]
     ```

5. **User taps "Yes, Send Alerts"**:
   - Call `POST /confirm-escalation` with `confirm: true`
   - Show: "Emergency alerts sent successfully"
   - Show emergency contacts, call buttons

6. **User taps "Cancel"**:
   - Call `POST /confirm-escalation` with `confirm: false`
   - Show: "Emergency escalation cancelled"
   - Dismiss modal
   - Continue polling

## React Native Implementation Suggestions

### State Management
```javascript
const [nudgePolling, setNudgePolling] = useState(false);
const [currentNudge, setCurrentNudge] = useState(null);
const [escalationDetails, setEscalationDetails] = useState(null);
const [travellerID, setTravellerID] = useState(1); // Get from user session
```

### Polling Function
```javascript
useEffect(() => {
  if (!nudgePolling) return;
  
  const pollInterval = setInterval(async () => {
    try {
      const response = await fetch(
        `${API_BASE_URL}/gosolo/travellers/${travellerID}/poll-safety-nudges`
      );
      const data = await response.json();
      
      if (data.length > 0) {
        // Nudge received - show modal
        setCurrentNudge(data[0]);
        setNudgePolling(false); // Stop polling until user responds
      }
    } catch (error) {
      console.error('Polling error:', error);
    }
  }, 15000); // Poll every 15 seconds
  
  return () => clearInterval(pollInterval);
}, [nudgePolling, travellerID]);
```

### Handle Response
```javascript
const handleSafetyResponse = async (response) => {
  try {
    const apiResponse = await fetch(
      `${API_BASE_URL}/gosolo/travellers/${travellerID}/respond-safety-check`,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          nudge_id: currentNudge.nudge_id,
          response: response // "yes_safe" or "no_trouble"
        })
      }
    );
    
    const data = await apiResponse.json();
    
    if (response === "yes_safe") {
      // Show success, dismiss modal, resume polling
      showSuccessMessage(data.message);
      setCurrentNudge(null);
      setNudgePolling(true);
    } else {
      // Show confirmation modal with escalation details
      setEscalationDetails(data);
      showConfirmationModal(data);
    }
  } catch (error) {
    console.error('Response error:', error);
  }
};
```

### Confirm Escalation
```javascript
const handleEscalationConfirmation = async (confirm) => {
  try {
    const response = await fetch(
      `${API_BASE_URL}/gosolo/travellers/${travellerID}/confirm-escalation`,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          escalation_id: escalationDetails.escalation_id,
          confirm: confirm
        })
      }
    );
    
    const data = await response.json();
    
    if (confirm) {
      showEmergencyScreen(data);
    } else {
      showMessage(data.message);
      setEscalationDetails(null);
      setCurrentNudge(null);
      setNudgePolling(true); // Resume polling
    }
  } catch (error) {
    console.error('Escalation error:', error);
  }
};
```

## Error Handling

- **Network errors**: Retry polling after delay
- **401/403**: Redirect to login
- **500**: Show error message, allow retry
- **Timeout**: Handle gracefully, continue polling

## Testing Checklist

- [ ] Polling works every 10-20 seconds
- [ ] Nudge modal appears when received
- [ ] "Yes, I'm Safe" works correctly
- [ ] "No, I'm in Trouble" shows confirmation modal
- [ ] "Yes, Send Alerts" sends escalation
- [ ] "Cancel" cancels escalation
- [ ] Polling resumes after response
- [ ] Error handling works
- [ ] Works offline (queue requests, send when online)

## API Base Configuration

Create a config file:
```javascript
// config/api.js
export const API_CONFIG = {
  BASE_URL: __DEV__ 
    ? 'http://localhost:8005' 
    : 'https://gosolo-4.onrender.com',
  POLL_INTERVAL: 15000, // 15 seconds
  TIMEOUT: 10000 // 10 seconds
};
```

That's everything you need to integrate the safety check polling and emergency escalation features. The backend handles all the logic, you just need to poll, display modals, and make POST requests based on user actions.

