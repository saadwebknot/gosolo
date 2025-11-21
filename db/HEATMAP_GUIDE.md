# Heat Map Data Guide - JP Nagar, Bengaluru

## Your Location (Center Point)
- **Latitude:** 12.914599
- **Longitude:** 77.595036
- **Address:** 8th Main Road, JP Nagar, Bengaluru - 560078

## Hazard Zones Distribution (within 5km radius)

### High Crime Areas (Red Zones - safe_rating 1.0-2.5)
1. **JP Nagar Market Backlanes** (12.910000, 77.590000)
   - 0.5km from center
   - Severity: HIGH
   - Safe Rating: 2.1

2. **Bannerghatta Road Underpass** (12.900000, 77.600000)
   - 1.6km from center
   - Severity: HIGH
   - Safe Rating: 1.8

3. **BTM Layout Industrial Area** (12.920000, 77.610000)
   - 1.8km from center
   - Severity: HIGH
   - Safe Rating: 2.3

### Medium Risk Areas (Orange Zones - safe_rating 2.6-3.5)
4. **Jayanagar 4th Block Park Perimeter** (12.930000, 77.580000)
   - 1.7km from center
   - Severity: MEDIUM
   - Safe Rating: 3.2

5. **Banashankari Bus Stand Area** (12.905000, 77.575000)
   - 2.2km from center
   - Severity: MEDIUM
   - Safe Rating: 3.0

6. **HSR Layout Outer Ring Road** (12.940000, 77.640000)
   - 5.0km from center
   - Severity: MEDIUM
   - Safe Rating: 3.4

### Lower Risk Areas (Yellow Zones - safe_rating 3.6-4.0)
7. **JP Nagar Phase 7 Construction Zone** (12.915000, 77.605000)
   - 1.0km from center
   - Severity: MEDIUM
   - Safe Rating: 3.8

8. **Bilekahalli Main Road** (12.925000, 77.620000)
   - 2.8km from center
   - Severity: LOW
   - Safe Rating: 3.9

### Safe Areas (Green Zones - safe_rating 4.1-5.0)
9. **JP Nagar Metro Station Area** (12.914599, 77.595036)
   - 0km from center (your location)
   - Severity: LOW
   - Safe Rating: 4.8 (won't trigger nudges)

10. **Forum Mall Vicinity** (12.920000, 77.590000)
    - 0.6km from center
    - Severity: LOW
    - Safe Rating: 4.9 (won't trigger nudges)

## Heat Map Color Coding

- **🔴 Red (High Crime):** safe_rating 1.0 - 2.5
- **🟠 Orange (Medium Risk):** safe_rating 2.6 - 3.5
- **🟡 Yellow (Caution):** safe_rating 3.6 - 4.0
- **🟢 Green (Safe):** safe_rating 4.1 - 5.0

## Testing Scenarios

### Scenario 1: At Your Current Location
- Location: 12.914599, 77.595036
- Expected: No nudges (you're in a safe zone with rating 4.8)

### Scenario 2: Near Market Backlanes
- Location: 12.910000, 77.590000
- Expected: HIGH severity nudge about market backlanes

### Scenario 3: Near Bannerghatta Underpass
- Location: 12.900000, 77.600000
- Expected: HIGH severity nudge about dark underpass

### Scenario 4: Near Construction Zone
- Location: 12.915000, 77.605000
- Expected: MEDIUM severity nudge about construction

## How to Test

1. Update your traveller's location:
```bash
curl -X POST http://localhost:8005/gosolo/travellers/1/location \
  -H "Content-Type: application/json" \
  -d '{"lat":12.910000,"lng":77.590000,"source":"gps","network_lost":false}'
```

2. Fetch nudges:
```bash
curl http://localhost:8005/gosolo/travellers/1/nudges
```

3. Try different coordinates to see different nudges!

