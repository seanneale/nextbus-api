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

	// Build the DB - Uncomment if required
	PopulateRoutesTable()

	// fmt.Println("testing stops")
	// // testing DB connection - CREATE
	// err = DB.QueryRow("INSERT INTO nextbus.stops (name) VALUES ($1)", "becky").Err()
	// if err != nil {
	// 	log.Printf("error creating data from Postgresql DB: %v", err)
	// }

	// //  testing DB connection - GET ALL
	// rows, err := DB.Query("SELECT id, name FROM nextbus.stops;")
	// if err != nil {
	// 	log.Printf("error reading data from Postgresql DB: %v", err)
	// }
	// for rows.Next() {
	// 	var stopTest StopTest
	// 	rows.Scan(&stopTest.Id, &stopTest.Name)
	// 	log.Printf("%v", stopTest)
	// }

	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal("Error occurred while starting the server:", err)
	}
}
