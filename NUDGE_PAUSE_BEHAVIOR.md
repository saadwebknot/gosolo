# Nudge Pause Behavior - Updated Logic

## Overview
When nudges are paused, **safety check nudges** ("Are you alright?") are stopped, but **conditional nudges** (weather, crime rate, area hazards, SOS) continue to work normally.

## Two Types of Nudges

### 1. Safety Check Nudges (Affected by Pause)
- **Endpoint**: `GET /gosolo/travellers/{id}/poll-safety-nudges`
- **Type**: Time-based frequency nudges
- **Message**: "Are you alright?"
- **Behavior when paused**: Returns empty array `[]` - no nudges sent

### 2. Conditional Nudges (NOT Affected by Pause)
- **Endpoint**: `GET /gosolo/travellers/{id}/nudges`
- **Types**:
  - Area/Crime rate warnings (based on location)
  - Weather alerts (based on location)
  - SOS follow-up nudges (if SOS is active)
  - Network loss alerts
- **Behavior when paused**: Continue to work normally - still returned

## Pause Endpoints

### Pause Nudges
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

### Resume Nudges
```bash
curl -X POST http://localhost:8005/gosolo/travellers/1/resume-nudges \
  -H "Content-Type: application/json"
```

## Testing Behavior

### Test 1: When Nudges are Active (Not Paused)

**Safety Check Polling:**
```bash
curl http://localhost:8005/gosolo/travellers/1/poll-safety-nudges
```
**Result**: Returns safety check nudge when frequency threshold is hit

**Conditional Nudges:**
```bash
curl http://localhost:8005/gosolo/travellers/1/nudges
```
**Result**: Returns area/weather/SOS nudges based on location

---

### Test 2: When Nudges are Paused

**Step 1: Pause nudges**
```bash
curl -X POST http://localhost:8005/gosolo/travellers/1/pause-nudges \
  -H "Content-Type: application/json" \
  -d '{"duration": "1hour"}'
```

**Step 2: Check Safety Check Polling**
```bash
curl http://localhost:8005/gosolo/travellers/1/poll-safety-nudges
```
**Result**: Returns empty array `[]` - no safety check nudges

**Step 3: Check Conditional Nudges**
```bash
curl http://localhost:8005/gosolo/travellers/1/nudges
```
**Result**: Still returns area/weather/SOS nudges - these continue working!

**Step 4: Resume nudges**
```bash
curl -X POST http://localhost:8005/gosolo/travellers/1/resume-nudges \
  -H "Content-Type: application/json"
```

**Step 5: Verify Safety Check Resumed**
```bash
curl http://localhost:8005/gosolo/travellers/1/poll-safety-nudges
```
**Result**: Safety check nudges resume (returns nudge when frequency hit)

## Rationale

**Why conditional nudges continue when paused:**
- **Safety critical** - Weather alerts and crime warnings are safety-related, not annoyances
- **Location-based** - These are contextual warnings about immediate surroundings
- **Emergency information** - Users need to know about nearby hazards even if they paused regular check-ins

**Why safety check nudges stop when paused:**
- **User preference** - User explicitly paused check-in nudges
- **Frequency-based** - These are scheduled reminders, not safety alerts
- **Reduces notification fatigue** - Prevents constant "Are you alright?" messages in safe zones

## Summary Table

| Nudge Type | Endpoint | Affected by Pause? |
|-----------|----------|-------------------|
| Safety Check ("Are you alright?") | `/poll-safety-nudges` | ✅ YES - Returns `[]` when paused |
| Area/Crime warnings | `/nudges` | ❌ NO - Continue normally |
| Weather alerts | `/nudges` | ❌ NO - Continue normally |
| SOS follow-ups | `/nudges` | ❌ NO - Continue normally |
| Network loss alerts | `/nudges` | ❌ NO - Continue normally |

## React Native Implementation Note

When implementing the pause feature in React Native:
- **Stop polling** `/poll-safety-nudges` when paused (or poll but ignore empty responses)
- **Continue fetching** `/nudges` for conditional alerts even when paused
- Show indicator that "Check-in nudges paused" but "Safety alerts active"

