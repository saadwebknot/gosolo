# go-mysql-crud
Sample crud operation using Golang and MySql

## API ENDPOINTS

### All Posts
- Path : `/posts`
- Method: `GET`
- Response: `200`

### Create Post
- Path : `/posts`
- Method: `POST`
- Fields: `title, content`
- Response: `201`

### Details a Post
- Path : `/posts/{id}`
- Method: `GET`
- Response: `200`

### Update Post
- Path : `/posts/{id}`
- Method: `PUT`
- Fields: `title, content`
- Response: `200`

### Delete Post
- Path : `/posts/{id}`
- Method: `DELETE`
- Response: `204`

### List Test Records
- Path : `/tests`
- Method: `GET`
- Response: `200`
- Description: Returns all rows from the `test` table (name + age).

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

