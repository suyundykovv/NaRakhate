package servers

import (
	"Aitu-Bet/internal/models"
	"Aitu-Bet/internal/services"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (s *Server) createBet(bet models.Bet) (*models.Bet, error) {
	if err := services.DeductBetAmountFromUser(bet, s.db); err != nil {
		return nil, err
	}

	if bet.CreatedAt.IsZero() {
		bet.CreatedAt = time.Now()
	}

	income, betValue, err := services.GetPotentialWinValue(bet.EventID, bet.OddSelection, bet.Amount, s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to get potential win value: %w", err)
	}
	bet.OddValue = betValue
	bet.Income = income
	bet.Status = "open"
	query := `
		INSERT INTO bets (
			user_id, event_id, amount, odd_value, income, status, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	err = s.db.QueryRow(query,
		bet.UserID,
		bet.EventID,
		bet.Amount,
		bet.OddValue,
		bet.Income,
		bet.Status,
		bet.CreatedAt,
	).Scan(&bet.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create bet: %w", err)
	}
	return &bet, nil
}

func (s *Server) readAllBets() ([]models.Bet, error) {
	query := `
		SELECT id, user_id, event_id, odd_selection, odd_value, amount, COALESCE(income, 0) as income, status, created_at 
		FROM bets
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query bets: %w", err)
	}
	defer rows.Close()

	var bets []models.Bet
	for rows.Next() {
		var bet models.Bet
		if err := rows.Scan(
			&bet.ID,
			&bet.UserID,
			&bet.EventID,
			&bet.OddSelection,
			&bet.OddValue,
			&bet.Amount,
			&bet.Income,
			&bet.Status,
			&bet.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan bet: %w", err)
		}
		bets = append(bets, bet)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return bets, nil
}

func (s *Server) readBetByID(id int) (*models.Bet, error) {
	query := `
		SELECT id, user_id, event_id, odd_selection, odd_value, amount, income, status, created_at 
		FROM bets 
		WHERE id = $1
	`
	var bet models.Bet
	err := s.db.QueryRow(query, id).Scan(
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

func (s *Server) updateBet(bet models.Bet) error {
	query := `
		UPDATE bets 
		SET user_id = $1, event_id = $2, odd_selection = $3, odd_value = $4, amount = $5, income = $6, status = $7, created_at = $8 
		WHERE id = $9
	`
	_, err := s.db.Exec(query,
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

func (s *Server) deleteBetData(id int) error {
	query := `DELETE FROM bets WHERE id = $1`
	_, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete bet with id %d: %w", id, err)
	}
	return nil
}
