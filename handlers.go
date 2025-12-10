package main

import (
	"html/template"
	"net/http"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Error loading home page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func handleTsoWoHang(w http.ResponseWriter, r *http.Request) {
	type response struct {
		TsoWoHang, MakPin []busTime
	}

	tmpl, err := template.ParseFiles("templates/tsowohang.html")
	if err != nil {
		http.Error(w, "Error loading home page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, response{TsoWoHang: kmb_bus_times("B959226950B0DEA7"), MakPin: kmb_bus_times("B2F4485FA517FEED")})
	if err != nil {
		return
	}
}

func handleLearningCenter(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ShaKokMei, HongKongAcademy []busTime
	}

	tmpl, err := template.ParseFiles("templates/learningcenter.html")
	if err != nil {
		http.Error(w, "Error loading home page 1", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, response{ShaKokMei: kmb_bus_times("F85F7F6FEB0812B5"), HongKongAcademy: kmb_bus_times("1AE7D87716BC52B1")})
	if err != nil {
		return
	}
}
