package main

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

// Test the health endpoint
func TestHealth(t *testing.T) {

	dummyCheckDBHealth := func() bool { return true }

	req, _ := http.NewRequest("GET", "/", nil)
	rr := executeRequest(req, handleGetHealth(dummyCheckDBHealth))

	// check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	// unmarshal response body into struct
	var data struct {
		Status    string `json:"db"`
		Version   string `json:"version"`
		Timestamp string `json:"timestamp"`
	}
	err := json.Unmarshal(rr.Body.Bytes(), &data)
	if err != nil {
		t.Fatal(err)
	}

	// check database status
	if data.Status != "Ok" {
		t.Errorf("expected Ok, got %s", data.Status)
	}

	// check version
	if data.Version != version {
		t.Errorf("expected %s, got %s", version, data.Version)
	}

	// check timestamp
	_, err = time.Parse(time.RFC3339, data.Timestamp)
	if err != nil {
		t.Fatal(err)
	}
}
