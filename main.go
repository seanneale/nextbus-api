package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	// load env file
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// handle routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/tsowohang", handleTsoWoHang)
	mux.HandleFunc("/learningcenter", handleLearningCenter)

	// Start the server
	err = http.ListenAndServe(":4001", mux)
	if err != nil {
		log.Fatal("Error occurred while starting the server:", err)
	}
}
