package servers

import (
	"Aitu-Bet/config"
	"Aitu-Bet/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (s *Server) FetchAndSaveMatchesHandler(w http.ResponseWriter, r *http.Request) {
	var allMatches []models.Fixture

	date := r.URL.Query().Get("date")
	if date == "" {
		http.Error(w, "Date query parameter is required", http.StatusBadRequest)
		return
	}

	log.Printf("Fetching matches for date: %s", date)
	matches, err := s.fetchFootballMatchesForDay(date)
	if err != nil {
		log.Printf("Error fetching football matches for date %s: %v", date, err)
		http.Error(w, "Failed to fetch matches", http.StatusInternalServerError)
		return
	}

	err = s.saveMatchesToDB(matches)
	if err != nil {
		log.Printf("Failed to save football match data to DB for date %s: %v", date, err)
		http.Error(w, "Failed to save matches", http.StatusInternalServerError)
		return
	}

	allMatches = append(allMatches, matches...)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(allMatches); err != nil {
		log.Printf("Error encoding football matches to JSON: %v", err)
		http.Error(w, "Failed to encode football matches to JSON", http.StatusInternalServerError)
	}
}

func (s *Server) fetchFootballMatchesForDay(date string) ([]models.Fixture, error) {
	var fixtures []models.Fixture
	url := fmt.Sprintf("https://v3.football.api-sports.io/fixtures?date=%s", date)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("x-apisports-key", config.ApiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %w", err)
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
	return fixtures, nil
}

func (s *Server) saveMatchesToDB(fixtures []models.Fixture) error {
	for _, fixture := range fixtures {
		matchDate := fixture.FixtureDetails.Date

		var matchExists bool
		err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM matches WHERE home_team_id=$1 AND away_team_id=$2 AND match_date=$3)",
			fixture.Teams.Home.ID, fixture.Teams.Away.ID, matchDate).Scan(&matchExists)
		if err != nil {
			log.Println("Error checking if match exists:", err)
			return err
		}

		if !matchExists {
			_, err = s.db.Exec(`
				INSERT INTO matches (home_team_id, away_team_id, home_goals, away_goals, match_date, league_id, referee, venue_name, venue_city)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			`, fixture.Teams.Home.ID, fixture.Teams.Away.ID,
				fixture.Goals.Home, fixture.Goals.Away, matchDate,
				fixture.League.ID, fixture.FixtureDetails.Referee,
				fixture.FixtureDetails.Venue.Name, fixture.FixtureDetails.Venue.City)
			if err != nil {
				log.Println("Error inserting match data:", err)
				return err
			}
			log.Printf("Inserted match: %s vs %s on %s", fixture.Teams.Home.Name, fixture.Teams.Away.Name, matchDate)
		} else {
			log.Printf("Match already exists: %s vs %s on %s", fixture.Teams.Home.Name, fixture.Teams.Away.Name, matchDate)
		}
	}

	return nil
}

func (s *Server) FetchLeagueMatchesHandler(w http.ResponseWriter, r *http.Request) {
	var allMatches []models.Fixture

	// Parse the league from query parameter
	league := r.URL.Query().Get("league")
	if league == "" {
		http.Error(w, "League query parameter is required", http.StatusBadRequest)
		return
	}

	log.Printf("Fetching matches for league: %s", league)
	matches, err := s.fetchFootballMatchesForLeague(league)
	if err != nil {
		log.Printf("Error fetching football matches for league %s: %v", league, err)
		http.Error(w, "Failed to fetch matches", http.StatusInternalServerError)
		return
	}

	err = s.saveLeagueMatchesToDB(matches)
	if err != nil {
		log.Printf("Failed to save football match data to DB for league %s: %v", league, err)
		http.Error(w, "Failed to save matches", http.StatusInternalServerError)
		return
	}

	allMatches = append(allMatches, matches...)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(allMatches); err != nil {
		log.Printf("Error encoding football matches to JSON: %v", err)
		http.Error(w, "Failed to encode football matches to JSON", http.StatusInternalServerError)
	}
}

func (s *Server) fetchFootballMatchesForLeague(league string) ([]models.Fixture, error) {
	var fixtures []models.Fixture
	url := fmt.Sprintf("https://v3.football.api-sports.io/fixtures?league=%s&season=2023", league)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("x-apisports-key", config.ApiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %w", err)
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
	return fixtures, nil
}

func (s *Server) saveLeagueMatchesToDB(fixtures []models.Fixture) error {
	for _, fixture := range fixtures {
		matchDate := fixture.FixtureDetails.Date

		// Check if the match already exists in the database
		var matchExists bool
		err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM matches WHERE home_team_id=$1 AND away_team_id=$2 AND match_date=$3)",
			fixture.Teams.Home.ID, fixture.Teams.Away.ID, matchDate).Scan(&matchExists)
		if err != nil {
			log.Println("Error checking if match exists:", err)
			return err
		}

		if !matchExists {
			// Insert the match into the matches table
			_, err = s.db.Exec(`
				INSERT INTO matches (home_team_id, away_team_id, home_goals, away_goals, match_date, league_id, referee, venue_name, venue_city)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			`, fixture.Teams.Home.ID, fixture.Teams.Away.ID,
				fixture.Goals.Home, fixture.Goals.Away, matchDate,
				fixture.League.ID, fixture.FixtureDetails.Referee,
				fixture.FixtureDetails.Venue.Name, fixture.FixtureDetails.Venue.City)
			if err != nil {
				log.Println("Error inserting match data:", err)
				return err
			}
			log.Printf("Inserted match: %s vs %s on %s", fixture.Teams.Home.Name, fixture.Teams.Away.Name, matchDate)
		} else {
			log.Printf("Match already exists: %s vs %s on %s", fixture.Teams.Home.Name, fixture.Teams.Away.Name, matchDate)
		}
	}

	return nil
}
