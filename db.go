// used as reference: https://www.youtube.com/watch?v=k4ZULXJwDGo
package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq" // what does _ do?
)

var DB *sql.DB

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
