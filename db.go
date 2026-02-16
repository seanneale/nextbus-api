// used as reference: https://www.youtube.com/watch?v=k4ZULXJwDGo
package main

import (
	"database/sql"
	"fmt"
	"log"
	"slices"
	"strings"

	_ "github.com/lib/pq" // what does _ do?
)

var DB *sql.DB

type StopInfo struct {
	Id, KmbStopId string
}

type RouteInfo struct {
	Id, RouteNo, Bound, ServiceType, Company string
}

type RouteStopInfo struct {
	RouteId, StopId string
}

func OpenDatabase() error {
	var err error
	DB, err = sql.Open("postgres", "user=seanneale dbname=nextbus sslmode=disable")
	if err != nil {
		return err
	}
	return nil
}

func CloseDatabase() error {
	return DB.Close()
}

func PopulateRoutesTable() {
	routes := kmb_route_list()
	// TODO: Add back chinese names
	sql_string := "INSERT INTO nextbus.routes (route_no, company, bound, service_type, orig_en, dest_en) VALUES"
	var sql_strings []string
	replacer := strings.NewReplacer("(", "[", ")", "]", "'", "")
	for _, route := range routes {
		sql_strings = append(sql_strings, fmt.Sprintf("('%s', '%s', '%s', %d, '%s', '%s')", route.RouteNo, route.Company, route.Bound, route.ServiceType, replacer.Replace(route.OrigEn), replacer.Replace(route.DestEn)))
	}
	sql_string += strings.Join(sql_strings, ",")
	err := DB.QueryRow(sql_string).Err()
	if err != nil {
		log.Printf("error creating data from Postgresql DB: %v", err)
	}
}

func PopulateStopsTable() {
	stops := kmbStopList()

	// Grouping together stops with the same lat and long so they can be treated as one
	groupedStopsByLatLong := make(map[string]stopInfo)

	for _, stop := range stops {
		key := fmt.Sprintf("%f,%f", stop.Latitude, stop.Longitude)
		if _, exists := groupedStopsByLatLong[key]; exists {
			existingStop := groupedStopsByLatLong[key]
			existingStop.KmbStopId += fmt.Sprintf(",%s", stop.KmbStopId)
			groupedStopsByLatLong[key] = existingStop
		} else {
			groupedStopsByLatLong[key] = stop
		}
	}

	// Large amount of data being inserted so using a single query to avoid multiple round trips to the database
	sql_string := "INSERT INTO nextbus.stops (name, kmb_stop_id, kmb_name_en, kmb_name_sc, kmb_name_tc, latitude, longitude) VALUES"
	var sql_strings []string

	replacer := strings.NewReplacer("(", "[", ")", "]", "'", "")
	for _, stop := range groupedStopsByLatLong {
		// Update naming convertion
		name := replacer.Replace(stop.KmbNameEn)
		sql_strings = append(sql_strings, fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', %f, %f)", name, stop.KmbStopId, replacer.Replace(stop.KmbNameEn), stop.KmbNameSc, stop.KmbNameTc, stop.Latitude, stop.Longitude))
	}
	sql_string += strings.Join(sql_strings, ",")
	err := DB.QueryRow(sql_string).Err()
	if err != nil {
		log.Printf("error creating data from Postgresql DB: %v", err)
	}
}

func PopulateRouteStopsTable() {
	routeStops := kmbRouteStopList()

	// retrieve all routes
	routeRows, err := DB.Query("SELECT id, route_no, bound, service_type, company FROM nextbus.routes;")
	if err != nil {
		log.Printf("error reading data from Postgresql DB: %v", err)
	}
	var allRoutes []RouteInfo
	for routeRows.Next() {
		var routeInfo RouteInfo
		routeRows.Scan(&routeInfo.Id, &routeInfo.RouteNo, &routeInfo.Bound, &routeInfo.ServiceType, &routeInfo.Company)
		allRoutes = append(allRoutes, routeInfo)
	}

	// retrieve all stops
	stopRows, err := DB.Query("SELECT id, kmb_stop_id FROM nextbus.stops;")
	if err != nil {
		log.Printf("error reading data from Postgresql DB: %v", err)
	}
	var allStops []StopInfo
	for stopRows.Next() {
		var stopInfo StopInfo
		stopRows.Scan(&stopInfo.Id, &stopInfo.KmbStopId)
		allStops = append(allStops, stopInfo)
	}

	sql_string := "INSERT INTO nextbus.routestops (route_id, stop_id) VALUES"
	var sql_strings []string

	var routeStopInfos []RouteStopInfo

	for _, routeStop := range routeStops {
		var routeStopInfo RouteStopInfo

		for _, stop := range allStops {
			stopIds := strings.Split(stop.KmbStopId, ",")

			if slices.Contains(stopIds, stop.KmbStopId) {
				routeStopInfo.StopId = stop.Id
				// TODO: Move to method and return a value
			}
		}

		for _, route := range allRoutes {
			if route.RouteNo == routeStop.RouteNo && route.Bound == routeStop.Bound && route.ServiceType == routeStop.ServiceType && route.Company == "KMB" {
				routeStopInfo.RouteId = route.Id
				// TODO: Move to method and return a value
			}
		}
		if routeStopInfo.RouteId != "" && routeStopInfo.StopId != "" {
			sql_strings = append(sql_strings, fmt.Sprintf("('%s', '%s')", routeStopInfo.RouteId, routeStopInfo.StopId))
		}
	}
	sql_string += strings.Join(sql_strings, ",")
	err = DB.QueryRow(sql_string).Err()
	if err != nil {
		log.Printf("error creating data from Postgresql DB: %v", err)
	}
	fmt.Println(routeStopInfos)
}
