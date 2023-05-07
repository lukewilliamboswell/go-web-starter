package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type ContextKey string

// Constants
const DB_KEY ContextKey = "db"

// Version is changed at build time
var version string = ""

// These variables read from environment variables
var (
	secret   string
	db       *sql.DB
	server   string
	port     int
	user     string
	password string
	database string
)

func main() {
	var err error

	// Check version variable has been set at build time
	if version == "" {
		log.Fatalf("Error version not set at build time\n")
	}

	// Get server details from environment variables
	server = os.Getenv("DB_SERVER")
	portStr := os.Getenv("DB_PORT")
	user = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	database = os.Getenv("DB_NAME")
	secret = os.Getenv("APP_SECRET")

	// Check secret environment variable has been set
	if secret == "" {
		log.Fatalf("Error secret not set\n")
	}

	// Convert port string to integer
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Error converting port string to integer: %s\n", err.Error())
	}

	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatalf("Error creating connection pool: %s\n", err.Error())
	}

	// Ping server to verify connection
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Connected to DB %s:%d as %s\n", server, port, user)

	// Create router
	router := chi.NewRouter()

	// Add middleware
	router.Use(middleware.Logger)

	// Add database handle to request context
	router.Use(middleware.WithValue(DB_KEY, db))

	// Register handlers
	RegisterHandlers(router)

	// Start server
	log.Printf("Starting server on port 8080\n")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Fatal error starting server: %s\n", err.Error())
	}
}
