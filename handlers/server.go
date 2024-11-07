package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// ScrapeAndSaveDataScraping performs scraping logic without HTTP context
func ScrapeAndSaveDataScraping() error {
	// Perform the scraping
	data, err := scrapeEPLSchedule()
	if err != nil {
		log.Println("Error scraping:", err)
		return err
	}

	// Save the scraped data
	err = saveMatchesToFile(data)
	if err != nil {
		log.Println("Error saving data:", err)
		return err
	}

	log.Println("Data scraped and saved successfully")
	return nil
}

// ScrapeAndSaveData is the HTTP handler for scraping and saving data
func ScrapeAndSaveData(w http.ResponseWriter, r *http.Request) {
	// Trigger scraping and saving inside the HTTP handler
	err := ScrapeAndSaveDataScraping()
	if err != nil {
		http.Error(w, "Scraping failed", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"status": "Data scraped and saved successfully"}
	json.NewEncoder(w).Encode(response)
}

// ServeData handles requests to serve the JSON data from the file
func ServeData(w http.ResponseWriter, r *http.Request) {
	// Log the request for debugging
	log.Println("Received request:", r.Method, r.URL)

	// Handle OPTIONS method for CORS preflight
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Handle GET method
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Content-Type", "application/json")

	// Open the match data file
	file, err := os.Open("epl_match_data.json") // Change this line to use epl_match_data.json
	if err != nil {
		log.Println("Error opening file:", err)
		http.Error(w, "Could not read data", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Log if file is successfully opened
	log.Println("File opened successfully")

	// Read the file content into a variable
	fileContent, err := os.ReadFile("epl_match_data.json") // Same here, change the file name
	if err != nil {
		log.Println("Error reading file:", err)
		http.Error(w, "Failed to read data", http.StatusInternalServerError)
		return
	}

	// Log the response before sending
	log.Println("Serving data to client")

	// Serve the JSON data
	w.Write(fileContent)
}
