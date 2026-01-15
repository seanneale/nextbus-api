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

func handleWuKaiSha(w http.ResponseWriter, r *http.Request) {
	tables := []Table{
		{Name: "Wu Kai Sha Station", BusTime: kmb_bus_times([]string{"9AD30F8EDBC3F139", "8134FC79F33F203D", "0CB1F7979192FBB2", "541FF05FA053F2B5", "7D54FE486D057070", "91CD1740B6AC752B", "E22B819E638307CC", "BA72214DFE48AA86"})},
	}
	renderFixedTable(tables, w, r)
}

func handleKeiLingHaLoWai(w http.ResponseWriter, r *http.Request) {
	tables := []Table{
		{Name: "Kei Ling Ha Lo Wai (Towards Sai Kung)", BusTime: kmb_bus_times([]string{"46E8837F566582D2"})},
		{Name: "Kei Ling Ha Lo Wai (Towards Ma On Shan)", BusTime: kmb_bus_times([]string{"6BA43C06A9AD502A"})},
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
