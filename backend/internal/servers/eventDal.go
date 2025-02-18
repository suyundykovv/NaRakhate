package servers

import (
	"Aitu-Bet/internal/models"
	"database/sql"
	"errors"
)

func (s *Server) createEvent(event models.Event) (*models.Event, error) {
	err := s.db.QueryRow(
		"INSERT INTO events (name, description, start_time, category) VALUES ($1, $2, $3, $4) RETURNING id",
		event.Name, event.Description, event.StartTime, event.Category).Scan(&event.ID)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (s *Server) readAllEvents() ([]models.Event, error) {
	var events []models.Event

	rows, err := s.db.Query(`
        SELECT id, name, description, start_time, category, home_win_odds, away_win_odds, draw_odds, match_status
        FROM events
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event models.Event
		var homeWinOdds, awayWinOdds, drawOdds float64

		if err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Description,
			&event.StartTime,
			&event.Category,
			&homeWinOdds,
			&awayWinOdds,
			&drawOdds,
			&event.Status,
		); err != nil {
			return nil, err
		}

		event.Odd = models.Odds{
			HomeWin: homeWinOdds,
			AwayWin: awayWinOdds,
			Draw:    drawOdds,
		}

		events = append(events, event)
	}

	return events, rows.Err()
}

func (s *Server) readEventByID(id int) (*models.Event, error) {
	var event models.Event

	err := s.db.QueryRow(
		"SELECT id, name, description, start_time, category, match_status, home_goals, away_goals FROM events WHERE id = $1", id).
		Scan(
			&event.ID,
			&event.Name,
			&event.Description,
			&event.StartTime,
			&event.Category,
			&event.Status,
			&event.Home_goals,
			&event.Away_goals,
		)

	if err == sql.ErrNoRows {
		return nil, errors.New("event not found")
	} else if err != nil {
		return nil, err
	}

	return &event, nil
}

func (s *Server) updateEvent(event models.Event) error {
	_, err := s.db.Exec(
		"UPDATE events SET name = $1, description = $2, start_time = $3, category = $4 WHERE id = $5",
		event.Name, event.Description, event.StartTime, event.Category, event.ID)
	return err
}

func (s *Server) deleteEventData(id int) error {
	_, err := s.db.Exec("DELETE FROM events WHERE id = $1", id)
	return err
}

func (s *Server) readAllFixtures() ([]models.Fixture, error) {
	var fixtures []models.Fixture

	rows, err := s.db.Query(`
		SELECT id, home_team_id, away_team_id, home_goals, away_goals, match_date, league_id, referee, venue_name, venue_city
		FROM matches
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var fixture models.Fixture
		var fixtureDetails models.FixtureDetails
		var teams models.Teams
		var goals models.Goals
		var league models.League
		var venue models.Venue

		if err := rows.Scan(
			&fixtureDetails.ID, &teams.Home.ID, &teams.Away.ID,
			&goals.Home, &goals.Away, &fixtureDetails.Date,
			&league.ID, &fixtureDetails.Referee,
			&venue.Name, &venue.City); err != nil {
			return nil, err
		}

		fixture.FixtureDetails = fixtureDetails
		fixture.Teams = teams
		fixture.Goals = goals
		fixture.League = league
		fixture.FixtureDetails.Venue = venue

		fixtures = append(fixtures, fixture)
	}

	return fixtures, rows.Err()
}
