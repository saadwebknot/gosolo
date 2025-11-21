package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/s1s1ty/go-mysql-crud/driver"
	ph "github.com/s1s1ty/go-mysql-crud/handler/http"
)

func main() {
	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}

	connection, err := driver.ConnectSQL(dbHost, dbPort, dbUser, dbPass, dbName)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
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

	fmt.Println("Server listen at :8005")
	http.ListenAndServe(":8005", r)
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
	r.Get("/travellers/{travellerID:[0-9]+}/nudges", handler.FetchNudges)
	r.Get("/travellers/{travellerID:[0-9]+}/hotel-brief", handler.HotelBrief)
	r.Get("/travellers/{travellerID:[0-9]+}/solo-buddy", handler.FindSoloBuddy)
	return r
}
