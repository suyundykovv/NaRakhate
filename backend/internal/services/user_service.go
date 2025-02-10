package services

import (
	"Aitu-Bet/internal/models"
	"database/sql"
	"fmt"
	"log"
)

func DeductBetAmountFromUser(bet models.Bet, db *sql.DB) error {
	var currentBalance float64
	err := db.QueryRow("SELECT cash FROM users WHERE id = $1", bet.UserID).Scan(&currentBalance)
	if err != nil {
		return fmt.Errorf("failed to retrieve user balance: %w", err)
	}

	if currentBalance < bet.Amount {
		return fmt.Errorf("insufficient funds: available balance is %.2f, but bet amount is %.2f", currentBalance, bet.Amount)
	}

	newBalance := currentBalance - bet.Amount

	_, err = db.Exec("UPDATE users SET cash = $1 WHERE id = $2", newBalance, bet.UserID)
	if err != nil {
		return fmt.Errorf("failed to update user balance: %w", err)
	}
	return nil
}

func AddWinningsForClosedBet(betID int, db *sql.DB) error {
	var (
		oddSelection string
		oddValue     float64
		amount       float64
		income       float64
		userID       int
		eventID      int
	)
	betQuery := `
		SELECT odd_selection, odd_value, amount, income, user_id, event_id 
		FROM bets 
		WHERE id = $1
	`
	err := db.QueryRow(betQuery, betID).Scan(&oddSelection, &oddValue, &amount, &income, &userID, &eventID)
	if err != nil {
		return fmt.Errorf("failed to query bet: %w", err)
	}

	var (
		homeGoals int
		awayGoals int
	)
	eventQuery := `
		SELECT home_goals, away_goals 
		FROM events 
		WHERE id = $1
	`
	err = db.QueryRow(eventQuery, eventID).Scan(&homeGoals, &awayGoals)
	if err != nil {
		return fmt.Errorf("failed to query event: %w", err)
	}

	var outcome string
	switch oddSelection {
	case "home":
		if homeGoals > awayGoals {
			outcome = "win"
		} else {
			outcome = "lose"
		}
	case "away":
		if awayGoals > homeGoals {
			outcome = "win"
		} else {
			outcome = "lose"
		}
	case "draw":
		if homeGoals == awayGoals {
			outcome = "win"
		} else {
			outcome = "lose"
		}
	default:
		return fmt.Errorf("invalid odd_selection: %s", oddSelection)
	}

	if outcome == "win" {
		updateQuery := `
			UPDATE users 
			SET cash = cash + $1 
			WHERE id = $2
		`
		_, err = db.Exec(updateQuery, income, userID)
		if err != nil {
			return fmt.Errorf("failed to update user cash: %w", err)
		}
		log.Printf("Added winnings for bet %d: user %d received %f", betID, userID, income)
	} else {
		log.Printf("Bet %d lost: no winnings added for user %d", betID, userID)
	}

	return nil
}
