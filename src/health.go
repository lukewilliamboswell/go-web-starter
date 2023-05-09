package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Returns a JSON response with the database status, version and timestamp
// {"db":"Ok","version":"latest","timestamp":"2023-05-07T06:28:46Z"}
func handleGetHealth(checkDBHealth func() bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// check database connection
		var db_status string
		if checkDBHealth() {
			db_status = "Ok"
		} else {
			db_status = "No Response"
		}

		// create JSON response
		data := struct {
			Status    string `json:"db"`
			Version   string `json:"version"`
			Timestamp string `json:"timestamp"`
		}{
			Status:    db_status,
			Version:   version,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}

		// convert data to JSON
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}

		// Set content type header
		w.Header().Set("Content-Type", "application/json")

		// write JSON response to client
		w.Write(jsonData)
	}
}

// Returns a function that checks the database connection
func checkDBHealth(db *sql.DB) func() bool {
	return func() bool {
		err := db.Ping()
		if err != nil {
			return false
		} else {
			return true
		}
	}
}
