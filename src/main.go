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

	"github.com/lukewilliamboswell/go-web-starter/src/user"
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
	dbServer string
	dbPort   int
	dbUser   string
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
	dbServer = os.Getenv("DB_SERVER")
	portStr := os.Getenv("DB_PORT")
	dbUser = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	database = os.Getenv("DB_NAME")
	secret = os.Getenv("APP_SECRET")

	// Check secret environment variable has been set
	if secret == "" {
		log.Fatalf("Error secret not set\n")
	}

	// Convert port string to integer
	dbPort, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Error converting port string to integer: %s\n", err.Error())
	}

	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		dbServer, dbUser, password, dbPort, database)

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
	log.Printf("Connected to DB %s:%d as %s\n", dbServer, dbPort, dbUser)

	// Create repositories
	userRepository := user.NewRepository(db)

	// Create controllers
	userController := user.NewController(userRepository)

	// Create router
	router := chi.NewRouter()

	// Add middleware
	router.Use(middleware.Logger)
	router.Use(userController.UserAuthenticationMiddleware())

	// Add database handle to request context
	router.Use(middleware.WithValue(DB_KEY, db))

	// Register handlers
	router.Get("/", handleGetRoot)
	router.Get("/headers", handleGetHeaders) // TODO remove
	router.Get("/health", handleGetHealth(checkDBHealth(db)))
	userController.RegisterRoutes(router)

	// Start server
	log.Printf("Starting server on port 8080\n")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Fatal error starting server: %s\n", err.Error())
	}
}
