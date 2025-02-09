package services

import (
	"database/sql"
	"fmt"
	"log"
)

func GetPotentialWinValue(eventID int, oddSelection string, betAmount float64, db *sql.DB) (float64, float64, error) {
	var homeWinOdds, awayWinOdds, drawOdds float64

	query := `
		SELECT home_win_odds, away_win_odds, draw_odds
		FROM events 
		WHERE id = $1
	`
	err := db.QueryRow(query, eventID).Scan(&homeWinOdds, &awayWinOdds, &drawOdds)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get event odds: %w", err)
	}

	var oddValue float64
	switch oddSelection {
	case "home":
		oddValue = homeWinOdds
	case "away":
		oddValue = awayWinOdds
	case "draw":
		oddValue = drawOdds
	default:
		return 0, 0, fmt.Errorf("invalid odd selection: %s", oddSelection)
	}

	potentialWin := betAmount * oddValue
	return potentialWin, oddValue, nil
}

func UpdateAllBetsIfMatchFinished(db *sql.DB) error {
	rows, err := db.Query(`
		SELECT id, event_id, odd_selection 
		FROM bets 
		WHERE status = 'open'
	`)
	if err != nil {
		return fmt.Errorf("failed to query open bets: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var betID, eventID int
		var oddSelection string
		if err := rows.Scan(&betID, &eventID, &oddSelection); err != nil {
			return fmt.Errorf("failed to scan bet: %w", err)
		}

		if err := UpdateBetStatusIfMatchFinished(betID, eventID, oddSelection, db); err != nil {
			log.Printf("Error updating bet %d: %v", betID, err)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows iteration error: %w", err)
	}
	return nil
}

func UpdateBetStatusIfMatchFinished(betID int, eventID int, oddSelection string, db *sql.DB) error {
	var homeGoals, awayGoals int
	var matchStatus string

	query := `
		SELECT home_goals, away_goals, match_status
		FROM events 
		WHERE id = $1
	`
	err := db.QueryRow(query, eventID).Scan(&homeGoals, &awayGoals, &matchStatus)
	if err != nil {
		return fmt.Errorf("failed to query event details: %w", err)
	}

	if matchStatus != "Match Finished" {
		return nil
	}

	var outcome string
	if homeGoals > awayGoals {
		if oddSelection == "home" {
			outcome = "win"
		} else {
			outcome = "loss"
		}
	} else if awayGoals > homeGoals {
		if oddSelection == "away" {
			outcome = "win"
		} else {
			outcome = "loss"
		}
	} else {
		if oddSelection == "draw" {
			outcome = "win"
		} else {
			outcome = "loss"
		}
	}

	updateQuery := `
		UPDATE bets
		SET status = 'closed'
		WHERE id = $1
	`
	_, err = db.Exec(updateQuery, betID)
	if err != nil {
		return fmt.Errorf("failed to update bet status: %w", err)
	}

	if outcome == "win" {
		if err := AddWinningsForClosedBet(betID, db); err != nil {
			return fmt.Errorf("failed to add winnings: %w", err)
		}
	}

	return nil
}

func isEventFinished(eventID int, db *sql.DB) (bool, error) {
	var matchStatus string
	query := `SELECT match_status FROM events WHERE id = $1`
	err := db.QueryRow(query, eventID).Scan(&matchStatus)
	if err != nil {
		return false, fmt.Errorf("failed to get event status: %w", err)
	}
	return matchStatus == "Match Finished", nil
}

func CheckEventStatusForBet(eventID int, db *sql.DB) error {
	finished, err := isEventFinished(eventID, db)
	if err != nil {
		return err
	}
	if finished {
		return fmt.Errorf("cannot create bet for finished event")
	}
	return nil
}
