package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// load env file

	if os.Getenv("GO_ENV") != "heroku" {
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
	}

	// handle routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/tsowohang", handleTsoWoHang)
	mux.HandleFunc("/learningcenter", handleLearningCenter)
	mux.HandleFunc("/akungkok", handleAKungKok)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal("Error occurred while starting the server:", err)
	}
}
