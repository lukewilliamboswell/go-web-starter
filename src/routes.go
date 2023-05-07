package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

func RegisterHandlers(router *chi.Mux) {
	router.Get("/", handleGetRoot)
	router.Get("/headers", handleGetHeaders)
	router.Get("/health", handleGetHealth)
}

func handleGetRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!\n"))
}

func handleGetHeaders(w http.ResponseWriter, r *http.Request) {
	// print out all the request headers in response body
	for name, headers := range r.Header {
		for _, h := range headers {
			w.Write([]byte(name + ": " + h + "\n"))
		}
	}

	w.Write([]byte("\n"))
}

// Returns a JSON response with the database status, version and timestamp
// {"db":"Ok","version":"latest","timestamp":"2023-05-07T06:28:46Z"}
func handleGetHealth(w http.ResponseWriter, r *http.Request) {

	// get database handle from context
	ctx := r.Context()
	db, ok := ctx.Value(DB_KEY).(*sql.DB)
	if !ok {
		log.Print("Error could not get database handle from context")
		http.Error(w, "could not get database handle from context", http.StatusInternalServerError)
		return
	}

	// check database connection
	var db_status string
	if db.PingContext(ctx) != nil {
		db_status = "No Response"
	} else {
		db_status = "Ok"
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
