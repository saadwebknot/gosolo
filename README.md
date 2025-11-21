## Go-SOLO: Solo Traveller Safety MVP
FeelSafe prototype for hotels—now branded **Go-SOLO** with the **SOLO-BUDDY** pro guide feature.

### Core API Surface
- `POST /gosolo/travellers` — onboard a traveller (name, phone, hotel info, location permission)
- `POST /gosolo/travellers/{id}/location` — push the latest coordinates + network status (stores last known location if offline)
- `POST /gosolo/travellers/{id}/sos` — trigger SOS (enables continuous nudges until stopped)
- `GET /gosolo/travellers/{id}/nudges` — fetch conditional nudges (area, weather, SOS follow-up, network)
- `GET /gosolo/travellers/{id}/hotel-brief` — returns localized hotel prompt for in-app audio display
- `GET /gosolo/travellers/{id}/solo-buddy` — list available local guides (SOLO-BUDDY pro feature)

### Database Entities
`travellers`, `sos_events`, `hazard_zones`, `weather_alerts`, `solo_guides` added to `db/init.sql` alongside sample hazard/weather/guide data.

## Required Packages
- Go 1.20+ (module aware)
- Docker & Docker Compose
- MySQL 5.7 (or let docker-compose manage it)

## Quick Run Project
First clone the repo then go to go-mysql-crud folder. After that build your image and run by docker. Make sure you have docker in your machine. 

```
git clone https://github.com/s1s1ty/go-mysql-crud.git

cd go-mysql-crud

chmod +x run.sh
./run.sh

docker compose up --build
```

### Render.com Deployment

**Build Command:** `./run.sh`  
**Start Command:** `./dist/go-mysql-crud`

**Required Environment Variables** (set in Render dashboard):
- `DB_HOST` - Your MySQL hostname (e.g., from PlanetScale, AWS RDS, or external service)
- `DB_PORT` - MySQL port (usually `3306`)
- `DB_USER` - MySQL username (defaults to `root` if not set)
- `DB_PASS` - MySQL password
- `DB_NAME` - Database name (e.g., `go-mysql-crud`)

**Note:** Render doesn't provide managed MySQL. You'll need:
- An external MySQL service (PlanetScale, AWS RDS, etc.), OR
- A separate Render service running MySQL via Docker

**Quick PlanetScale Setup (Recommended - 5 minutes):**
1. Sign up at https://planetscale.com (free tier available)
2. Create a new database (e.g., `gosolo`)
3. Go to "Connect" → copy the connection string
4. It looks like: `mysql://USER:PASS@HOST:PORT/DATABASE?ssl-mode=REQUIRED`
5. Extract values and set in Render:
   - `DB_HOST` = the hostname (e.g., `aws.connect.psdb.cloud`)
   - `DB_PORT` = `3306` (usually)
   - `DB_USER` = username from connection string
   - `DB_PASS` = password from connection string
   - `DB_NAME` = database name
6. Run your `db/init.sql` schema via PlanetScale's SQL editor or CLI

**Local Docker Compose defaults** (for reference):
- `DB_HOST=mysql` (service name)
- `DB_PORT=3306`
- `DB_USER=root`
- `DB_PASS=12345`
- `DB_NAME=go-mysql-crud`

### Database Maintenance (create/drop tables manually)
Run SQL through the MySQL container:
```bash
docker compose exec mysql mysql -uroot -p12345 go-mysql-crud -e "CREATE TABLE demo (...);"
docker compose exec mysql mysql -uroot -p12345 go-mysql-crud -e "DROP TABLE demo;"
```

### Handy cURL Tests
1. Register traveller
   ```bash
   curl -s -X POST http://localhost:8005/gosolo/travellers \
     -H "Content-Type: application/json" \
     -d '{"name":"Ava Solo","phone":"+1-202-555-0199","hotel_name":"Aurora Suites","hotel_address":"221B MG Road","hotel_language_prompt":"आप ऑरोरा सूट्स में ठहर रहे हैं।","location_permission":true}'
   ```
2. Update location (store last known coordinates, even if the network drops)
   ```bash
   curl -s -X POST http://localhost:8005/gosolo/travellers/1/location \
     -H "Content-Type: application/json" \
     -d '{"lat":28.6139,"lng":77.2090,"source":"gps","network_lost":false}'
   ```
3. Fetch dynamic nudges
   ```bash
   curl -s "http://localhost:8005/gosolo/travellers/1/nudges" | jq
   ```
4. Trigger SOS
   ```bash
   curl -s -X POST http://localhost:8005/gosolo/travellers/1/sos \
     -H "Content-Type: application/json" \
     -d '{"channel":"app","notes":"Walking back to hotel"}'
   ```
5. List SOLO-BUDDY guides
   ```bash
   curl -s http://localhost:8005/gosolo/travellers/1/solo-buddy | jq
   ```

