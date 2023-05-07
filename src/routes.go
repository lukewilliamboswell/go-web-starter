package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func RegisterHandlers(router *chi.Mux) {
	router.Get("/", handleGetRoot)
	router.Get("/version", handleGetVersion)
	router.Get("/headers", handleGetHeaders)
}

func handleGetRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!\n"))
}

func handleGetVersion(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(version + "\n"))
}

func handleGetHeaders(w http.ResponseWriter, r *http.Request) {
	// print out all the request headers in response body
	for name, headers := range r.Header {
		for _, h := range headers {
			w.Write([]byte(name + ": " + h))
		}
	}

	w.Write([]byte("\n"))
}
