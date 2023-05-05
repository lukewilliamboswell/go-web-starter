package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

func TestAppRouter(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)
	rr := executeRequest(req, hanldleGetRoot)

	// check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200")
	}

	// check the response body
	if rr.Body.String() != "Hello, World!" {
		t.Errorf("unexpected body")
	}
}

func executeRequest(req *http.Request, handler http.HandlerFunc) *httptest.ResponseRecorder {
	// Create new recorder
	rr := httptest.NewRecorder()

	// Create a router and setup handler
	r := chi.NewRouter()
	r.HandleFunc("/", handler)

	// Execute request
	r.ServeHTTP(rr, req)

	// Return recorder
	return rr
}
