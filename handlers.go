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
		tsoWoHang, makPin []busTime
	}

	// var tmplFile = "templates/tsowohang.html"
	tmpl, err := template.ParseFiles("templates/tsowohang.html")
	if err != nil {
		http.Error(w, "Error loading home page 1", http.StatusInternalServerError)
		return
	}

	tswResponse := response{tsoWoHang: kmb_bus_times("B959226950B0DEA7"), makPin: kmb_bus_times("C7548A3C37ADC1AA")}

	// fmt.Println(time.Now())
	fmt.Println("Rendering....")
	fmt.Println(tswResponse.makPin)

	err = tmpl.Execute(w, kmb_bus_times("C7548A3C37ADC1AA"))
	if err != nil {
		// http.Error(w, "Error loading home page 2", http.StatusInternalServerError)
		return
	}
}
