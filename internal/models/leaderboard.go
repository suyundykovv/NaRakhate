package models

import "time"

// Player представляет игрока с его выигрышами
type Player struct {
	ID            int64     `json:"id"`
	Username      string    `json:"username"`
	TotalWinnings float64   `json:"total_winnings"`
	CreatedAt     time.Time `json:"created_at"`
}
