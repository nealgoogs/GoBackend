package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
)

// Match represents a match structure
type Match struct {
	Team1 string `json:"team1"`
	Team2 string `json:"team2"`
}

// MatchResponse represents the structure of the JSON response
type MatchResponse map[string][]Match // A map where the key is a string (time) and the value is a slice of Match

// Exported handler function to get match data
func GetMatches(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")             // Allow all origins
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS") // Allow GET and OPTIONS methods

	// Handle preflight request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	cmd := exec.Command("python", "selenium/soccer_scrape.py") // Or "python3" if necessary

	// Run the Python script and capture output
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error running scraper: %v", err) // Log the error
		http.Error(w, "Error running scraper", http.StatusInternalServerError)
		return
	}

	// Parse JSON from Python script output
	var matches MatchResponse
	err = json.Unmarshal(output, &matches)
	if err != nil {
		log.Printf("Error parsing JSON: %v", err) // Log the error
		http.Error(w, "Error parsing JSON", http.StatusInternalServerError)
		return
	}

	// Set JSON response headers
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matches) // Directly encode the matches map
}
