package main

import (
	"Backend/handlers"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

// Input structure for the incoming request
type Input struct {
	Text string `json:"text"`
}

// CORS middleware
func enableCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// Main function to set up the server
func main() {
	// Run the scraper logic manually (without HTTP context)
	err := handlers.ScrapeAndSaveDataScraping()
	if err != nil {
		log.Println("Error in manual scrape:", err)
	} else {
		log.Println("Scraping test successful!")
	}

	// Set up routes
	http.HandleFunc("/generate-image", generateImageHandler)
	http.HandleFunc("/api/matches", handlers.ServeData)
	http.HandleFunc("/scraper", handlers.ScrapeAndSaveData)
	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

// Response structure
type Response struct {
	ImageURL string `json:"image_url"`
}

// Generate an image based on text input
func generateImageHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r) // Ensure CORS headers are set

	if r.Method == http.MethodOptions {
		// Respond to preflight requests
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input Input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("Received text: %s", input.Text)

	// URL encode the prompt
	encodedPrompt := url.QueryEscape(input.Text)
	apiURL := "https://image.pollinations.ai/prompt/" + encodedPrompt // API endpoint

	// Send GET request to Pollinations API
	resp, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, "Error calling Pollinations API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to generate image", http.StatusInternalServerError)
		return
	}

	// Return the image URL
	response := Response{ImageURL: apiURL} // Directly return the request URL for the generated image
	json.NewEncoder(w).Encode(response)
}
