package models

import "time"

type Bet struct {
	ID           int       `json:"id" `
	UserID       int       `json:"user_id"`
	EventID      int       `json:"event_id"`
	OddSelection string    `json:"odd_selection"`
	OddValue     float64   `json:"odd_value"`
	Amount       float64   `json:"amount"`
	Income       float64   `json:"income"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}
