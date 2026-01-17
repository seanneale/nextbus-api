// used as reference: https://www.youtube.com/watch?v=k4ZULXJwDGo
package main

import (
	"database/sql"

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
