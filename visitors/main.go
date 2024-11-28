package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

type Visitor struct {
	IPv4      string  `json:"ipv4"`
	IPv6      string  `json:"ipv6"`
	Timestamp string  `json:"timestamp"`
	Country   string  `json:"country"`
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Connect to SQLite database
func connectDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./visitors.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getIPAddresses(r *http.Request) (ipv4, ipv6 string) {
	// Extract IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "Unknown", "Unknown"
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return "Unknown", "Unknown"
	}

	// Check if it's IPv4 or IPv6
	if parsedIP.To4() != nil {
		return parsedIP.String(), "" // IPv4
	}
	return "", parsedIP.String() // IPv6
}

// Log visitor data into the database
func logVisitor(db *sql.DB, visitor Visitor) error {
	_, err := db.Exec(`
		INSERT INTO visitors (ipv4, ipv6, timestamp, country, city, latitude, longitude)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, visitor.IPv4, visitor.IPv6, visitor.Timestamp, visitor.Country, visitor.City, visitor.Latitude, visitor.Longitude)
	return err
}

// Handler for logging visitors
func visitorHandler(w http.ResponseWriter, r *http.Request) {
	// Extract IPv4 and IPv6 addresses

	log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)

	// Ignore requests for favicon
	if r.URL.Path == "/favicon.ico" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	ipv4, ipv6 := getIPAddresses(r)

	// Fetch current date
	today := time.Now().Format("02-01-2006") // DD-MM-YYYY format

	// Simulate GeoIP data (replace with real API calls if needed)
	geoData := Visitor{
		IPv4:      ipv4,
		IPv6:      ipv6,
		Timestamp: today,
		Country:   "Unknown",
		City:      "Unknown",
		Latitude:  0.0,
		Longitude: 0.0,
	}

	// Connect to the database
	db, err := connectDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		log.Printf("Database error: %v", err)
		return
	}
	defer db.Close()

	// Log visitor data
	err = logVisitor(db, geoData)
	if err != nil {
		http.Error(w, "Failed to log visitor", http.StatusInternalServerError)
		log.Printf("Logging error: %v", err)
		return
	}

	// Respond to the visitor
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(geoData)
}

func getGeoData(ip string) (Visitor, error) {
	apiKey := "f4407f7833ff803650e7fe9be85bed5d"
	resp, err := http.Get(fmt.Sprintf("http://api.ipstack.com/%s?access_key=%s", ip, apiKey))
	if err != nil {
		return Visitor{}, err
	}
	defer resp.Body.Close()

	var geoData struct {
		CountryName string  `json:"country_name"`
		City        string  `json:"city"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&geoData); err != nil {
		return Visitor{}, err
	}

	return Visitor{
		Country:   geoData.CountryName,
		City:      geoData.City,
		Latitude:  geoData.Latitude,
		Longitude: geoData.Longitude,
	}, nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080 if PORT is not set
	}
	http.HandleFunc("/", visitorHandler)
	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
