package servers

import (
	"Aitu-Bet/config"
	"Aitu-Bet/internal/models"
	"Aitu-Bet/internal/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (s *Server) fetchFootballMatches() ([]models.Fixture, error) {
	var fixtures []models.Fixture
	maxRetries := 3

	currentDate := time.Now().Format("2006-01-02")
	url := fmt.Sprintf("https://v3.football.api-sports.io/fixtures?date=%s", currentDate)

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

		var response models.FootballResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("error decoding response: %w", err)
		}

		fixtures = response.Response
		break
	}

	return fixtures, nil
}

func (s *Server) FetchFootballMatchesHandler(w http.ResponseWriter, r *http.Request) {
	leagueID := r.URL.Query().Get("league_id")

	matches, err := s.fetchFootballMatches()
	if err != nil {
		log.Printf("Error fetching football matches: %v", err)
		http.Error(w, "Failed to retrieve football matches", http.StatusInternalServerError)
		return
	}

	var filteredMatches []models.Fixture
	if leagueID != "" {
		filteredMatches = filterMatchesByLeague(matches, leagueID)
	} else {
		filteredMatches = matches
	}

	err = s.saveFootballMatchesToDB(filteredMatches)
	if err != nil {
		log.Printf("Failed to save filtered football matches to DB: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(filteredMatches); err != nil {
		log.Printf("Error encoding football matches to JSON: %v", err)
		http.Error(w, "Failed to encode football matches to JSON", http.StatusInternalServerError)
	}
}

func filterMatchesByLeague(fixtures []models.Fixture, leagueID string) []models.Fixture {
	var filteredFixtures []models.Fixture
	for _, fixture := range fixtures {
		if fmt.Sprintf("%d", fixture.League.ID) == leagueID {
			filteredFixtures = append(filteredFixtures, fixture)
		}
	}
	return filteredFixtures
}

func (s *Server) saveFootballMatchesToDB(fixtures []models.Fixture) error {
	for _, fixture := range fixtures {
		startTime := fixture.FixtureDetails.Date
		matchName := fixture.Teams.Home.Name + " vs " + fixture.Teams.Away.Name

		// Check if the event already exists
		var eventExists bool
		err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM events WHERE name=$1 AND start_time=$2)",
			matchName, startTime).Scan(&eventExists)
		if err != nil {
			log.Println("Error checking if event exists:", err)
			return err
		}

		// Create odds for the fixture
		odds, err := services.CreateOddsService(fixture, s.db)
		if err != nil {
			log.Println("Error creating odds:", err)
			return err
		}

		if !eventExists {
			// Insert new event if it doesn't exist
			_, err = s.db.Exec(`
				INSERT INTO events (
					name, description, start_time, category, referee, venue_name, venue_city, home_win_odds, away_win_odds, draw_odds, match_status
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			`, matchName,
				"Football match description", startTime, "Sports", fixture.FixtureDetails.Referee,
				fixture.FixtureDetails.Venue.Name, fixture.FixtureDetails.Venue.City, odds.HomeWin, odds.AwayWin, odds.Draw, fixture.FixtureDetails.Status.Long)
			if err != nil {
				log.Println("Error inserting match data:", err)
				return err
			}
			log.Printf("Inserted match: %s vs %s", fixture.Teams.Home.Name, fixture.Teams.Away.Name)
		} else {
			_, err = s.db.Exec(`
				UPDATE events 
				SET description = $1, category = $2, referee = $3, venue_name = $4, venue_city = $5, 
				    home_win_odds = $6, away_win_odds = $7, draw_odds = $8, match_status = $9
				WHERE name = $10 AND start_time = $11
			`, "Football match description", "Sports", fixture.FixtureDetails.Referee,
				fixture.FixtureDetails.Venue.Name, fixture.FixtureDetails.Venue.City,
				odds.HomeWin, odds.AwayWin, odds.Draw, fixture.FixtureDetails.Status.Long,
				matchName, startTime)
			if err != nil {
				log.Println("Error updating match data:", err)
				return err
			}
			log.Printf("Updated match: %s", matchName)
		}
	}

	return nil
}
