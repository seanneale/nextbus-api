package main

import (
	"html/template"
	"net/http"
)

type Table struct {
	Name    string
	BusTime []busTime
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Error loading home page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func handleTsoWoHang(w http.ResponseWriter, r *http.Request) {
	tables := []Table{
		{Name: "Tso Wo Hang", BusTime: kmb_bus_times([]string{"B959226950B0DEA7"})},
		{Name: "Mak Pin (Towards Ma On Shan)", BusTime: kmb_bus_times([]string{"B2F4485FA517FEED"})},
		{Name: "Mak Pin (Towards Sai Kung)", BusTime: kmb_bus_times([]string{"C7548A3C37ADC1AA"})},
	}

	renderFixedTable(tables, w, r)
}

func handleLearningCenter(w http.ResponseWriter, r *http.Request) {
	tables := []Table{
		{Name: "Sha Kok Mei", BusTime: kmb_bus_times([]string{"F85F7F6FEB0812B5"})},
		{Name: "Hong Kong Academy (8 Min Walk)", BusTime: kmb_bus_times([]string{"1AE7D87716BC52B1"})},
	}

	renderFixedTable(tables, w, r)
}

func handleAKungKok(w http.ResponseWriter, r *http.Request) {
	tables := []Table{
		{Name: "A Kung Kok", BusTime: kmb_bus_times([]string{"4B0E21EDF07F9C83", "B8BBCCA288E1F862", "6F03E19C5E800893"})},
	}
	renderFixedTable(tables, w, r)
}

func handleGoPark(w http.ResponseWriter, r *http.Request) {
	tables := []Table{
		{Name: "Tin Liu (Go Park: Toward Ma On Shan)", BusTime: kmb_bus_times([]string{"79B88295548184C1"})},
		{Name: "Ma Kwu Lam (Go Park: Toward Sai Kung)", BusTime: kmb_bus_times([]string{"C67E139743E66D71"})},
	}
	renderFixedTable(tables, w, r)
}

func renderFixedTable(tables []Table, w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/fixed_table.html")
	if err != nil {
		http.Error(w, "Error loading home page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, tables)
	if err != nil {
		return
	}
}
