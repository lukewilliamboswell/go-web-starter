package main

import (
	"net/http"
)

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
