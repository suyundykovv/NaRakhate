package dal

import (
	"Aitu-Bet/internal/models"
	"Aitu-Bet/internal/services"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

func CreateBet(bet models.Bet, db *sql.DB) (*models.Bet, error) {
	log.Println("Starting createBet function")

	// Step 1: Deduct bet amount from user
	log.Printf("Deducting bet amount %f from user %d", bet.Amount, bet.UserID)
	if err := services.DeductBetAmountFromUser(bet, db); err != nil {
		log.Printf("Failed to deduct bet amount from user %d: %v", bet.UserID, err)
		return nil, err
	}
	log.Printf("Successfully deducted bet amount from user %d", bet.UserID)

	// Step 2: Set CreatedAt timestamp if not already set
	if bet.CreatedAt.IsZero() {
		bet.CreatedAt = time.Now()
		log.Printf("Set CreatedAt timestamp to %v", bet.CreatedAt)
	}

	// Step 3: Calculate potential winnings
	log.Printf("Calculating potential winnings for event %d, selection %s, amount %f", bet.EventID, bet.OddSelection, bet.Amount)
	income, betValue, err := services.GetPotentialWinValue(bet.EventID, bet.OddSelection, bet.Amount, db)
	if err != nil {
		log.Printf("Failed to calculate potential winnings: %v", err)
		return nil, fmt.Errorf("failed to get potential win value: %w", err)
	}
	bet.OddValue = betValue
	bet.Income = income
	log.Printf("Calculated potential winnings: odd_value=%f, income=%f", betValue, income)

	// Step 4: Set bet status
	bet.Status = "open"
	log.Printf("Set bet status to 'open'")

	// Step 5: Insert bet into database
	query := `
		INSERT INTO bets (
			user_id, event_id, amount, odd_value, income, status, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	log.Printf("Inserting bet into database: user_id=%d, event_id=%d, amount=%f, odd_value=%f, income=%f, status=%s, created_at=%v",
		bet.UserID, bet.EventID, bet.Amount, bet.OddValue, bet.Income, bet.Status, bet.CreatedAt)
	err = db.QueryRow(query,
		bet.UserID,
		bet.EventID,
		bet.Amount,
		bet.OddValue,
		bet.Income,
		bet.Status,
		bet.CreatedAt,
	).Scan(&bet.ID)
	if err != nil {
		log.Printf("Failed to insert bet into database: %v", err)
		return nil, fmt.Errorf("failed to create bet: %w", err)
	}
	log.Printf("Successfully inserted bet into database with ID %d", bet.ID)

	// Step 6: Return the created bet
	log.Printf("Bet created successfully: %+v", bet)
	return &bet, nil
}

func ReadAllBets(db *sql.DB) ([]models.Bet, error) {
	var bets []models.Bet

	rows, err := db.Query("SELECT id, user_id, event_id, amount, income, status, created_at FROM bets")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var bet models.Bet
		if err := rows.Scan(&bet.ID, &bet.UserID, &bet.EventID, &bet.Amount, &bet.Income, &bet.Status, &bet.CreatedAt); err != nil {
			return nil, err
		}
		bets = append(bets, bet)
	}

	return bets, rows.Err()
}

func ReadBetByID(id int, db *sql.DB) (*models.Bet, error) {
	query := `
		SELECT id, user_id, event_id, odd_selection, odd_value, amount, income, status, created_at 
		FROM bets 
		WHERE id = $1
	`
	var bet models.Bet
	err := db.QueryRow(query, id).Scan(
		&bet.ID,
		&bet.UserID,
		&bet.EventID,
		&bet.OddSelection,
		&bet.OddValue,
		&bet.Amount,
		&bet.Income,
		&bet.Status,
		&bet.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("bet not found")
		}
		return nil, fmt.Errorf("failed to read bet by ID: %w", err)
	}
	return &bet, nil
}

func UpdateBet(bet models.Bet, db *sql.DB) error {
	query := `
		UPDATE bets 
		SET user_id = $1, event_id = $2, odd_selection = $3, odd_value = $4, amount = $5, income = $6, status = $7, created_at = $8 
		WHERE id = $9
	`
	_, err := db.Exec(query,
		bet.UserID,
		bet.EventID,
		bet.OddSelection,
		bet.OddValue,
		bet.Amount,
		bet.Income,
		bet.Status,
		bet.CreatedAt,
		bet.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update bet: %w", err)
	}
	return nil
}

func DeleteBetData(id int, db *sql.DB) error {
	query := `DELETE FROM bets WHERE id = $1`
	_, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete bet with id %d: %w", id, err)
	}
	return nil
}
