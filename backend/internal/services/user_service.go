package services

import (
	"Aitu-Bet/internal/models"
	"database/sql"
	"fmt"
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
	var outcome string
	var income float64
	var userID int
	query := `
		SELECT outcome, income, user_id 
		FROM bets 
		WHERE id = $1
	`
	err := db.QueryRow(query, betID).Scan(&outcome, &income, &userID)
	if err != nil {
		return fmt.Errorf("failed to query bet: %w", err)
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
	}
	return nil
}
