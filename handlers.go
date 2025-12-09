package main

import (
	"fmt"
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

	// var tmplFile = "templates/tsowohang.html"
	tmpl, err := template.ParseFiles("templates/tsowohang.html")
	if err != nil {
		http.Error(w, "Error loading home page 1", http.StatusInternalServerError)
		return
	}

	tswResponse := response{TsoWoHang: kmb_bus_times("B959226950B0DEA7"), MakPin: kmb_bus_times("B2F4485FA517FEED")}

	// fmt.Println(time.Now())
	fmt.Println("Rendering....")
	fmt.Println(tswResponse.MakPin)

	for _, bus := range tswResponse.MakPin {
		fmt.Println(bus.RouteNo)
	}

	err = tmpl.Execute(w, tswResponse)
	if err != nil {
		// http.Error(w, "Error loading home page 2", http.StatusInternalServerError)
		return
	}
}

func handleLearningCenter(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ShaKokMei, HongKongAcademy []busTime
	}

	// var tmplFile = "templates/tsowohang.html"
	tmpl, err := template.ParseFiles("templates/learningcenter.html")
	if err != nil {
		http.Error(w, "Error loading home page 1", http.StatusInternalServerError)
		return
	}

	lcResponse := response{ShaKokMei: kmb_bus_times("F85F7F6FEB0812B5"), HongKongAcademy: kmb_bus_times("1AE7D87716BC52B1")}

	err = tmpl.Execute(w, lcResponse)
	if err != nil {
		// http.Error(w, "Error loading home page 2", http.StatusInternalServerError)
		return
	}
}
