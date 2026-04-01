// used as reference: https://www.youtube.com/watch?v=k4ZULXJwDGo
package main

import (
	"database/sql"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	_ "github.com/lib/pq" // what does _ do?
)

var DB *sql.DB

// potential naming clash with stopInfo in kmb.go
type StopInfo struct {
	Id, KmbStopId, GmbStopId string
	Latitude, Longitude      float64
}

type RouteInfo struct {
	Id, RouteNo, Bound, ServiceType, Company, GmbRouteId string
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

func PopulateKmbRoutesTable() {
	routes := kmb_route_list()
	// TODO: Add back chinese names
	sqlString := "INSERT INTO nextbus.routes (route_no, company, bound, service_type, orig_en, dest_en) VALUES"
	var sqlStrings []string
	replacer := strings.NewReplacer("(", "[", ")", "]", "'", "")
	for _, route := range routes {
		sqlStrings = append(sqlStrings, fmt.Sprintf("('%s', '%s', '%s', %d, '%s', '%s')", route.RouteNo, route.Company, route.Bound, route.ServiceType, replacer.Replace(route.OrigEn), replacer.Replace(route.DestEn)))
	}
	sqlString += strings.Join(sqlStrings, ",")
	writeToDb(sqlString)
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
	sqlString := "INSERT INTO nextbus.stops (name, kmb_stop_id, kmb_name_en, kmb_name_sc, kmb_name_tc, latitude, longitude) VALUES"
	var sqlStrings []string

	replacer := strings.NewReplacer("(", "[", ")", "]", "'", "")
	for _, stop := range groupedStopsByLatLong {
		// Update naming convertion
		name := replacer.Replace(stop.KmbNameEn)
		sqlStrings = append(sqlStrings, fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', %f, %f)", name, stop.KmbStopId, replacer.Replace(stop.KmbNameEn), stop.KmbNameSc, stop.KmbNameTc, stop.Latitude, stop.Longitude))
	}
	sqlString += strings.Join(sqlStrings, ",")
	writeToDb(sqlString)
}

func PopulateRouteStopsTable() {
	routeStops := kmbRouteStopList()

	// retrieve all routes
	routeRows, err := DB.Query("SELECT id, route_no, bound, service_type, company FROM nextbus.routes WHERE company='KMB';")
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
	stopRows, err := DB.Query("SELECT id, kmb_stop_id FROM nextbus.stops WHERE kmb_stop_id IS NOT NULL;")
	if err != nil {
		log.Printf("error reading data from Postgresql DB: %v", err)
	}
	var allStops []StopInfo
	for stopRows.Next() {
		var stopInfo StopInfo
		stopRows.Scan(&stopInfo.Id, &stopInfo.KmbStopId)
		allStops = append(allStops, stopInfo)
	}

	sqlString := "INSERT INTO nextbus.routestops (route_id, stop_id) VALUES"
	var sqlStrings []string

	for _, routeStop := range routeStops {
		var routeStopInfo RouteStopInfo

		for _, stop := range allStops {
			stopIds := strings.Split(stop.KmbStopId, ",")

			if slices.Contains(stopIds, routeStop.KmbStopId) {
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
			sqlStrings = append(sqlStrings, fmt.Sprintf("('%s', '%s')", routeStopInfo.RouteId, routeStopInfo.StopId))
		}
	}
	if len(sqlStrings) > 0 {
		sqlString += strings.Join(sqlStrings, ",")
		fmt.Println(sqlString)
		writeToDb(sqlString)
	}
}

func PopulateGmbRoutesTable() {
	routes := gmbRouteList()
	// TODO: Add back chinese names
	sqlString := "INSERT INTO nextbus.routes (route_no, company, bound, orig_en, dest_en, region, description_en, gmb_route_id) VALUES"
	var sqlStrings []string
	replacer := strings.NewReplacer("(", "[", ")", "]", "'", "")
	for _, route := range routes {
		sqlStrings = append(sqlStrings, fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%s', '%s', %d)", route.RouteNo, route.Company, route.Bound, replacer.Replace(route.OrigEn), replacer.Replace(route.DestEn), route.Region, replacer.Replace(route.DescriptionEn), route.GmbRouteId))
	}
	sqlString += strings.Join(sqlStrings, ",")
	writeToDb(sqlString)
}

func PopulateGmbStopsTable() {
	// The GMB API does not return all stops from one endpoint like the KMB API
	// We need to request the stops for each route, then group the ones with the same stop id to request data about individually.

	// Retrieve RouteID and Sequence Type from DB
	fmt.Println(currentTime(), "retrieveGmbRouteDataFromDb")
	allRoutes := retrieveGmbRouteDataFromDb()

	// Retrive stop data for each route
	fmt.Println(currentTime(), "gmbRouteStopList")
	newRouteStops, newStops := gmbRouteStopList(allRoutes)

	// Retrieve lat/long for stops from Stops API
	fmt.Println(currentTime(), "gmbStopLatLongList")
	newStops = gmbStopLatLongList(newStops)

	// Retrieve existing stops from DB
	// for GMB there's only unmatched stops, so we can skip this step for now
	// fmt.Println(currentTime(), "retrieveExistingStopsFromDb")
	// existingStops := retrieveExistingStopsFromDb()

	// find GMB stops with the same lat/long as the stops already in the DB
	// for GMB there's only unmatched stops, so we can skip this step for now
	// fmt.Println(currentTime(), "findMatchingStops")
	// unmatchedStops := findMatchingStops(existingStops, newStops)

	// insert stops into DB
	fmt.Println(currentTime(), "insertNewStopsIntoDB")
	insertNewStopsIntoDB(newStops)

	// extract newStops from DB
	fmt.Println(currentTime(), "extractNewStopsFromDB")
	newStopsDBInfo := extractNewStopsFromDB(newStops)

	fmt.Println(currentTime(), "insertNewRouteStopsIntoDb")
	insertNewRouteStopsIntoDb(newRouteStops, newStopsDBInfo, allRoutes)

	fmt.Println(currentTime(), "FINISHED")
}

func currentTime() string {
	timeLayout := "15:04"
	return time.Now().Format(timeLayout)
}

func retrieveGmbRouteDataFromDb() []RouteInfo {
	routeRows, err := DB.Query("SELECT id, route_no, bound, gmb_route_id FROM nextbus.routes WHERE company='GMB';")
	if err != nil {
		log.Printf("error reading data from Postgresql DB: %v", err)
	}
	var allRoutes []RouteInfo
	for routeRows.Next() {
		var routeInfo RouteInfo

		routeRows.Scan(&routeInfo.Id, &routeInfo.RouteNo, &routeInfo.Bound, &routeInfo.GmbRouteId)
		allRoutes = append(allRoutes, routeInfo)
	}

	return allRoutes
}

func retrieveExistingStopsFromDb() []StopInfo {
	stopRows, err := DB.Query("SELECT id, latitude, longitude FROM nextbus.stops;")
	if err != nil {
		log.Printf("error reading data from Postgresql DB: %v", err)
	}
	var existingStops []StopInfo
	for stopRows.Next() {
		var stopInfo StopInfo

		stopRows.Scan(&stopInfo.Id, &stopInfo.Latitude, &stopInfo.Longitude)
		existingStops = append(existingStops, stopInfo)
	}

	return existingStops
}

func findMatchingStops(existingStops []StopInfo, newStops []stopInfo) []stopInfo {
	var unmatchedStops []stopInfo
	for _, existStop := range existingStops {
		for _, newStop := range newStops {
			if existStop.Latitude == newStop.Latitude && existStop.Longitude == newStop.Longitude {
				// update the existing stop to include the GmbStopId
				// no matches for GMB and KMB, to be updated when Citybus are added
				break
			} else {
				// append the stop info to the sql string to be created in the DB later
				unmatchedStops = append(unmatchedStops, newStop)
			}
		}
	}

	return unmatchedStops
}

func insertNewStopsIntoDB(newStops []stopInfo) {
	groupedStopsByLatLong := make(map[string]stopInfo)

	for _, stop := range newStops {
		key := fmt.Sprintf("%f,%f", stop.Latitude, stop.Longitude)
		if _, exists := groupedStopsByLatLong[key]; exists {
			existingStop := groupedStopsByLatLong[key]
			existingStop.GmbStopId += fmt.Sprintf(",%s", stop.GmbStopId)
			groupedStopsByLatLong[key] = existingStop
		} else {
			groupedStopsByLatLong[key] = stop
		}
	}

	// Large amount of data being inserted so using a single query to avoid multiple round trips to the database
	sqlString := "INSERT INTO nextbus.stops (name, gmb_stop_id, gmb_name_en, gmb_name_sc, gmb_name_tc, latitude, longitude) VALUES"
	var sqlStrings []string

	replacer := strings.NewReplacer("(", "[", ")", "]", "'", "")
	for _, stop := range groupedStopsByLatLong {
		// Update naming convertion
		name := replacer.Replace(stop.GmbNameEn)
		sqlStrings = append(sqlStrings, fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', %f, %f)", name, stop.GmbStopId, replacer.Replace(stop.GmbNameEn), replacer.Replace(stop.GmbNameSc), replacer.Replace(stop.GmbNameTc), stop.Latitude, stop.Longitude))
	}
	sqlString += strings.Join(sqlStrings, ",")
	writeToDb(sqlString)
}

func extractNewStopsFromDB(newStops []stopInfo) []StopInfo {
	// make list of new GMB Stop Ids
	var gmbStopIDs []string
	for _, stop := range newStops {
		gmbStopIDs = append(gmbStopIDs, stop.GmbStopId)
	}
	sqlString := fmt.Sprintf("SELECT id, gmb_stop_id FROM nextbus.stops WHERE gmb_stop_id IN ('%s')", strings.Join(gmbStopIDs, "', '"))
	stopRows, err := DB.Query(sqlString)
	if err != nil {
		log.Printf("error reading data from Postgresql DB: %v", err)
	}

	var existingStops []StopInfo
	for stopRows.Next() {
		var stopInfo StopInfo

		stopRows.Scan(&stopInfo.Id, &stopInfo.GmbStopId)
		existingStops = append(existingStops, stopInfo)
	}

	return existingStops
}

func insertNewRouteStopsIntoDb(routeStops []routeStopInfo, allStops []StopInfo, allRoutes []RouteInfo) {
	// Notes for tomorrow:
	// 1. Pass allRoutes into method
	// 2. Retrieve GmbStopId and Id from DB
	sqlString := "INSERT INTO nextbus.routestops (route_id, stop_id) VALUES"
	var sqlStrings []string

	for _, routeStop := range routeStops {
		var routeStopInfo RouteStopInfo

		for _, stop := range allStops {
			stopIds := strings.Split(stop.GmbStopId, ",")

			if slices.Contains(stopIds, routeStop.GmbStopId) {
				routeStopInfo.StopId = stop.Id
				// TODO: Move to method and return a value
			}
		}

		for _, route := range allRoutes {
			if route.GmbRouteId == routeStop.GmbRouteId {
				routeStopInfo.RouteId = route.Id
				// TODO: Move to method and return a value
			}
		}
		if routeStopInfo.RouteId != "" && routeStopInfo.StopId != "" {
			sqlStrings = append(sqlStrings, fmt.Sprintf("('%s', '%s')", routeStopInfo.RouteId, routeStopInfo.StopId))
		}
	}
	if len(sqlStrings) > 0 {
		sqlString += strings.Join(sqlStrings, ",")
		writeToDb(sqlString)
	}

}

func writeToDb(sqlString string) {
	err := DB.QueryRow(sqlString).Err()
	if err != nil {
		log.Printf("error creating data from Postgresql DB: %v", err)
	}
}
