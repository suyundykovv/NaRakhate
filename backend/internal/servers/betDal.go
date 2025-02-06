package servers

import (
	"Aitu-Bet/internal/models"
	"database/sql"
	"errors"
)

func (s *Server) createBet(bet models.Bet) (*models.Bet, error) {
	err := s.db.QueryRow(
		"INSERT INTO bets (user_id, event_id, amount, outcome, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		bet.UserID, bet.EventID, bet.Amount, bet.Outcome, bet.CreatedAt).Scan(&bet.ID)

	if err != nil {
		return nil, err
	}
	return &bet, nil
}

func (s *Server) readAllBets() ([]models.Bet, error) {
	var bets []models.Bet

	rows, err := s.db.Query("SELECT id, user_id, event_id, amount, outcome, created_at FROM bets")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var bet models.Bet
		if err := rows.Scan(&bet.ID, &bet.UserID, &bet.EventID, &bet.Amount, &bet.Outcome, &bet.CreatedAt); err != nil {
			return nil, err
		}
		bets = append(bets, bet)
	}

	return bets, rows.Err()
}

func (s *Server) readBetByID(id int) (*models.Bet, error) {
	var bet models.Bet

	err := s.db.QueryRow(
		"SELECT id, user_id, event_id, amount, outcome, created_at FROM bets WHERE id = $1", id).
		Scan(&bet.ID, &bet.UserID, &bet.EventID, &bet.Amount, &bet.Outcome, &bet.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("bet not found")
	} else if err != nil {
		return nil, err
	}

	return &bet, nil
}

func (s *Server) updateBet(bet models.Bet) error {
	_, err := s.db.Exec(
		"UPDATE bets SET user_id = $1, event_id = $2, amount = $3, outcome = $4, created_at = $5 WHERE id = $6",
		bet.UserID, bet.EventID, bet.Amount, bet.Outcome, bet.CreatedAt, bet.ID)
	return err
}

func (s *Server) deleteBetData(id int) error {
	_, err := s.db.Exec("DELETE FROM bets WHERE id = $1", id)
	return err
}
