package servers

import (
	"Aitu-Bet/config"
	"Aitu-Bet/internal/models"
	"Aitu-Bet/logging"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (s *Server) fetchFootballMatches() ([]models.Fixture, error) {
	var fixtures []models.Fixture
	maxRetries := 3

	date := "2025-02-06"
	url := fmt.Sprintf("https://v3.football.api-sports.io/fixtures?date=%s", date)

	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %w", err)
		}

		req.Header.Set("x-apisports-key", config.ApiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Attempt %d: error fetching data: %v", i+1, err)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		// Define a struct for the response to map the incoming data
		var response models.FootballResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("error decoding response: %w", err)
		}

		// Directly use the Response from the decoded data
		fixtures = response.Response

		// If successful, break out of the retry loop
		break
	}

	return fixtures, nil
}

func (s *Server) saveFootballMatchesToDB(fixtures []models.Fixture) error {
	for _, fixture := range fixtures {
		startTime := fixture.FixtureDetails.Date

		matchName := fixture.Teams.Home.Name + " vs " + fixture.Teams.Away.Name

		var eventExists bool
		err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM events WHERE name=$1 AND start_time=$2)",
			matchName, startTime).Scan(&eventExists)
		if err != nil {
			log.Println("Error checking if event exists:", err)
			return err
		}

		if !eventExists {
			_, err = s.db.Exec(`
				INSERT INTO events (name, description, start_time, category, referee, venue_name, venue_city)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
			`, matchName,
				"Football match description", startTime, "Sports", fixture.FixtureDetails.Referee,
				fixture.FixtureDetails.Venue.Name, fixture.FixtureDetails.Venue.City)
			if err != nil {
				log.Println("Error inserting match data:", err)
				return err
			}
			log.Printf("Inserted match: %s vs %s", fixture.Teams.Home.Name, fixture.Teams.Away.Name)
		} else {
			log.Printf("Match %s already exists", matchName)
		}
	}

	return nil
}

func (s *Server) FetchFootballMatchesHandler(w http.ResponseWriter, r *http.Request) {
	matches, err := s.fetchFootballMatches()
	if err != nil {
		log.Printf("Error fetching football matches: %v", err)
		http.Error(w, "Failed to retrieve football matches", http.StatusInternalServerError)
		return
	}
	err = s.saveFootballMatchesToDB(matches)
	if err != nil {
		logging.Error("Failed to save football match data to DB", err)
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(matches); err != nil {
		log.Printf("Error encoding football matches to JSON: %v", err)
		http.Error(w, "Failed to encode football matches to JSON", http.StatusInternalServerError)
	}
}
