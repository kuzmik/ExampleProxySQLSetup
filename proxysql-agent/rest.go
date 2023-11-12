package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GetPodsRequest struct {
	Type string `json:"type"`
}

// curl -X POST http://localhost:8080/status
func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status": "ok"}`)
}

// curl -X POST http://localhost:8080/resync
func resyncHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	// Call your resync function here
	// Example: resync()
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Resync operation initiated")
}

// curl -d '{"type": "core"}' localhost:8080/get_pods
func getPodsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	var req GetPodsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	if req.Type != "core" && req.Type != "satellite" {
		http.Error(w, "Invalid 'type' value. Use 'core' or 'satellite'", http.StatusBadRequest)
		return
	}

	// Handle the request based on the 'type' value
	if req.Type == "core" {
		// Handle core type
	} else if req.Type == "satellite" {
		// Handle satellite type
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Get pods operation completed")
}

func StartAPI() {
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/resync", resyncHandler)
	http.HandleFunc("/get_pods", getPodsHandler)

	// FIXME: make configurable
	port := ":8080"

	// Start the HTTP server
	fmt.Printf("Server is running on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error starting the server: %v\n", err)
	}
}
