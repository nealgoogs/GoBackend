package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Match struct {
	Team1 string `json:"team1"`
	Team2 string `json:"team2"`
	Time  string `json:"time"`
}

type MatchesByTime map[string][]Match

// scrapeEPLSchedule scrapes the EPL schedule and groups matches by time.
func scrapeEPLSchedule() (MatchesByTime, error) {
	c := colly.NewCollector()

	// Set a User-Agent header to simulate a real browser request
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

	matchesByTime := make(MatchesByTime)

	// Scrape match details from the page
	c.OnHTML("tbody tr.TableBase-bodyTr", func(e *colly.HTMLElement) {
		// Get the team names
		homeTeam := e.ChildText("td:nth-child(1) .TeamName") // Home team name
		awayTeam := e.ChildText("td:nth-child(2) .TeamName") // Away team name

		// Get the match time
		matchTime := e.ChildText("td:nth-child(3) .CellGame") // Match time

		// Clean up the match time
		matchTime = strings.TrimSpace(matchTime)

		// Use regex to extract only the time (e.g., "9:00 am" or "11:30 am")
		re := regexp.MustCompile(`^\d{1,2}:\d{2} [ap]m`)
		if match := re.FindString(matchTime); match != "" {
			matchTime = match
		} else {
			return // Skip if no valid time found
		}

		// Ensure both teams and match time are present
		if homeTeam != "" && awayTeam != "" && matchTime != "" {
			match := Match{
				Team1: strings.TrimSpace(homeTeam),
				Team2: strings.TrimSpace(awayTeam),
				Time:  matchTime,
			}

			// Group matches by match time
			if _, exists := matchesByTime[match.Time]; !exists {
				matchesByTime[match.Time] = []Match{}
			}
			matchesByTime[match.Time] = append(matchesByTime[match.Time], match)
		}
	})

	// Visit the target URL
	err := c.Visit("https://www.cbssports.com/soccer/premier-league/schedule/20241110/")
	if err != nil {
		return nil, fmt.Errorf("error visiting the URL: %v", err)
	}

	return matchesByTime, nil
}

// saveMatchesToFile saves the grouped match data to a JSON file
func saveMatchesToFile(data MatchesByTime) error {
	file, err := os.Create("epl_match_data.json")
	if err != nil {
		return fmt.Errorf("error creating the file: %v", err)
	}
	defer file.Close()

	// Marshal the data into JSON with indentation
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling the data to JSON: %v", err)
	}

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("error writing to the file: %v", err)
	}

	return nil
}
