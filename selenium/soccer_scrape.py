from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.service import Service
from webdriver_manager.chrome import ChromeDriverManager
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC

def extract_champions_league_matches():
    driver = webdriver.Chrome(service=Service(ChromeDriverManager().install()))
    driver.get("https://www.espn.com/soccer/schedule/_/date/20241106")

    try:
        # Wait for the Champions League table to be present
        champions_league_table = WebDriverWait(driver, 10).until(
            EC.presence_of_element_located((By.XPATH, "//div[text()='UEFA Champions League']//following-sibling::table"))
        )

        # Find all rows in the Champions League table
        match_rows = champions_league_table.find_elements(By.CSS_SELECTOR, "tbody.Table__TBODY tr.Table__TR")

        matches = []
        for match_row in match_rows:
            teams = match_row.find_elements(By.CSS_SELECTOR, ".matchTeams .Table__Team")
            home_team = teams[0].text
            away_team = teams[1].text
            matches.append(f"{home_team} vs {away_team}")

        return matches

    except Exception as e:
        print(f"An error occurred: {e}")
    finally:
        driver.quit()

if __name__ == "__main__":
    champions_league_matches = extract_champions_league_matches()
    for match in champions_league_matches:
        print(match)