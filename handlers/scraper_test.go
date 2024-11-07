package handlers

import (
	"log"
	"testing"
)

// TestScrapeAndSave tests the scrapeMatchData and saveDataToFile functions
func TestScrapeAndSave(t *testing.T) {
	// Call the scraper function
	data, err := scrapeEPLSchedule()
	if err != nil {
		t.Fatalf("scrapeMatchData failed: %v", err)
	}

	// Check if any matches were scraped
	if len(data) == 0 {
		t.Error("No match data found; verify the CSS selectors or the page structure.")
	} else {
		log.Printf("Scraped %d match times\n", len(data))
		// Optionally print the match data for verification
		for matchTime, matches := range data {
			for _, match := range matches {
				log.Printf("Match: %s vs %s at %s", match.Team1, match.Team2, matchTime)
			}
		}
	}

	// Attempt to save the data to a file
	err = saveMatchesToFile(data)
	if err != nil {
		t.Fatalf("saveDataToFile failed: %v", err)
	}
}
