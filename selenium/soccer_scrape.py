from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
import json

# Set up Chrome options to ignore SSL errors
chrome_options = Options()
chrome_options.add_argument('--ignore-certificate-errors')
chrome_options.add_argument("--headless")  # Run in headless mode
chrome_options.add_argument("--no-sandbox")  # Bypass OS security model
chrome_options.add_argument("--disable-dev-shm-usage")

# Initialize WebDriver
driver = webdriver.Chrome(options=chrome_options)

def scrape_epl_schedule():
    # Dictionary to hold matches grouped by time
    matches_by_time = {}

    try:
        # Navigate to the soccer schedule page
        driver.get("https://www.espn.com/soccer/schedule/_/date/20241109")

        # Wait for the English Premier League section to load
        WebDriverWait(driver, 10).until(
            EC.presence_of_element_located((By.CLASS_NAME, "ScheduleTables--premier-league"))
        )

        # Locate the English Premier League schedule table
        premier_league_table = driver.find_element(By.CLASS_NAME, "ScheduleTables--premier-league")

        # Extract match rows within the table body
        match_rows = premier_league_table.find_elements(By.CSS_SELECTOR, ".Table__TBODY .Table__TR")

        # Loop through each match row and extract details
        for row in match_rows:
            teams = row.find_elements(By.CSS_SELECTOR, ".Table__Team")
            time_element = row.find_element(By.CSS_SELECTOR, ".date__col")

            match_time = time_element.text if time_element else "Unknown Time"
            match = {
                "team1": teams[0].text if len(teams) > 0 else None,
                "team2": teams[1].text if len(teams) > 1 else None
            }

            # Group matches by time
            if match_time not in matches_by_time:
                matches_by_time[match_time] = []
            matches_by_time[match_time].append(match)

    finally:
        driver.quit()  # Ensure the browser closes even if there's an error

    # Convert the grouped match data to JSON format
    matches_by_time_json = {
        time: [{"team1": match['team1'], "team2": match['team2']} for match in matches]
        for time, matches in matches_by_time.items()
    }

    # Convert to JSON and print
    json_output = json.dumps(matches_by_time_json, indent=4)
    print(json_output)

# Run the function to scrape and group EPL schedule by time
scrape_epl_schedule()
