package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type StopTest struct {
	Id, Name string
}

type RouteTest struct {
	Id, RouteNo string
}

func main() {
	var err error
	// load env file
	if os.Getenv("GO_ENV") != "heroku" {
		err = godotenv.Load()
		if err != nil {
			panic(err)
		}
	}

	err = OpenDatabase()
	if err != nil {
		log.Printf("error connecting to Postgresql DB: %v", err)
	}
	defer CloseDatabase()

	// handle routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/tsowohang", handleTsoWoHang)
	mux.HandleFunc("/learningcenter", handleLearningCenter)
	mux.HandleFunc("/akungkok", handleAKungKok)
	mux.HandleFunc("/gopark", handleGoPark)
	mux.HandleFunc("/wukaisha", handleWuKaiSha)
	mux.HandleFunc("/keilinghalowai", handleKeiLingHaLoWai)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	// Build the following DB tables - Uncomment if required
	// - Routes
	// - Stops
	// - RouteStops (Join the tables together)
	// PopulateRoutesTable()
	// PopulateStopsTable()
	// PopulateRouteStopsTable()

	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal("Error occurred while starting the server:", err)
	}
}
