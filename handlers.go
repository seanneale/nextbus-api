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

	err = tmpl.Execute(w, response{TsoWoHang: kmb_bus_times([]string{"B959226950B0DEA7"}), MakPin: kmb_bus_times([]string{"B2F4485FA517FEED"})})
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
		http.Error(w, "Error loading home page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, response{ShaKokMei: kmb_bus_times([]string{"F85F7F6FEB0812B5"}), HongKongAcademy: kmb_bus_times([]string{"1AE7D87716BC52B1"})})
	if err != nil {
		return
	}
}

func handleAKungKok(w http.ResponseWriter, r *http.Request) {
	type response struct {
		AKungKok []busTime
	}

	tmpl, err := template.ParseFiles("templates/akungkok.html")
	if err != nil {
		http.Error(w, "Error loading home page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, response{AKungKok: kmb_bus_times([]string{"4B0E21EDF07F9C83", "B8BBCCA288E1F862", "6F03E19C5E800893"})})
	if err != nil {
		return
	}
}
