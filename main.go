package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	"github.com/s1s1ty/go-mysql-crud/driver"
	ph "github.com/s1s1ty/go-mysql-crud/handler/http"
)

func main() {
	// Load .env file if it exists (ignore errors for production where env vars are set directly)
	_ = godotenv.Load()

	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPass == "" || dbName == "" {
		log.Fatal("Missing required database environment variables. Check DB_HOST, DB_PORT, DB_USER, DB_PASS, DB_NAME")
	}

	connection, err := driver.ConnectSQL(dbHost, dbPort, dbUser, dbPass, dbName)
	if err != nil {
		log.Fatalf("Database connection failed: %v\nCheck your DB credentials (DB_HOST=%s, DB_PORT=%s, DB_USER=%s, DB_NAME=%s)", err, dbHost, dbPort, dbUser, dbName)
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	pHandler := ph.NewPostHandler(connection)
	tHandler := ph.NewTestHandler(connection)
	goSoloHandler := ph.NewGoSoloHandler(connection)
	r.Route("/", func(rt chi.Router) {
		rt.Mount("/posts", postRouter(pHandler))
		rt.Get("/tests", tHandler.Fetch)
		rt.Mount("/gosolo", goSoloRouter(goSoloHandler))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8005"
	}

	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("Server listen at %s\n", addr)
	http.ListenAndServe(addr, r)
}

// A completely separate router for posts routes
func postRouter(pHandler *ph.Post) http.Handler {
	r := chi.NewRouter()
	r.Get("/", pHandler.Fetch)
	r.Get("/{id:[0-9]+}", pHandler.GetByID)
	r.Post("/", pHandler.Create)
	r.Put("/{id:[0-9]+}", pHandler.Update)
	r.Delete("/{id:[0-9]+}", pHandler.Delete)

	return r
}

func goSoloRouter(handler *ph.GoSolo) http.Handler {
	r := chi.NewRouter()
	r.Post("/travellers", handler.RegisterTraveller)
	r.Post("/travellers/{travellerID:[0-9]+}/location", handler.UpdateLocation)
	r.Post("/travellers/{travellerID:[0-9]+}/sos", handler.TriggerSOS)
	r.Post("/travellers/{travellerID:[0-9]+}/sos-button", handler.TriggerSOSButton)
	r.Get("/travellers/{travellerID:[0-9]+}/nudges", handler.FetchNudges)
	r.Get("/travellers/{travellerID:[0-9]+}/hotel-brief", handler.HotelBrief)
	r.Get("/travellers/{travellerID:[0-9]+}/solo-buddy", handler.FindSoloBuddy)
	// Nudge preferences management
	r.Put("/travellers/{travellerID:[0-9]+}/nudge-settings", handler.UpdateNudgeSettings)
	r.Post("/travellers/{travellerID:[0-9]+}/pause-nudges", handler.PauseNudges)
	r.Post("/travellers/{travellerID:[0-9]+}/resume-nudges", handler.ResumeNudges)
	// Safety check polling and responses
	r.Get("/travellers/{travellerID:[0-9]+}/poll-safety-nudges", handler.PollSafetyNudges)
	r.Post("/travellers/{travellerID:[0-9]+}/respond-safety-check", handler.RespondToSafetyCheck)
	r.Post("/travellers/{travellerID:[0-9]+}/confirm-escalation", handler.ConfirmEmergencyEscalation)
	return r
}
