package main

import (
	"compress/gzip"
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/lukewilliamboswell/go-web-starter/src/user"
)

//go:embed public/*
var staticFiles embed.FS

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
	dbPort, err = strconv.Atoi(portStr)
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

	// If this is a development build, add middleware to set the user header
	if version == "dev" {
		router.Use(devMiddlewareSetUserHeader)
	}

	// Add middleware to authenticate user, except for the following routes
	exceptions := []string{
		"/health",
		"/registered.html",
		"/styles.css",
		"/favicon.ico",
	}
	router.Use(userController.UserAuthenticationMiddleware(version, exceptions))

	// Add database handle to request context
	router.Use(middleware.WithValue(DB_KEY, db))

	// Register handlers
	router.Get("/headers", handleGetHeaders) // TODO remove
	router.Get("/health", handleGetHealth(checkDBHealth(db)))

	// Create a API router
	api := chi.NewRouter()

	// Register API routes
	userController.RegisterRoutes(api)

	// Mount API router
	router.Mount("/api", api)

	// Serve static files
	router.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/public" + r.URL.Path

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()

			gzw := gzipResponseWriter{ResponseWriter: w, Writer: gz}
			http.FileServer(http.FS(staticFiles)).ServeHTTP(gzw, r)
		} else {
			http.FileServer(http.FS(staticFiles)).ServeHTTP(w, r)
		}
	})

	// Start server
	log.Printf("Starting server on port 8080\n")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Fatal error starting server: %s\n", err.Error())
	}
}

type gzipResponseWriter struct {
	http.ResponseWriter
	*gzip.Writer
}

func (w gzipResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w gzipResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}

// Set the X-Ms-Client-Principal-Id for development builds
func devMiddlewareSetUserHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r.Header.Set("X-Ms-Client-Principal-Id", "bdd00b62-f7b6-40b4-b78f-d4b30f9dfe15")

		next.ServeHTTP(w, r)
	})
}
