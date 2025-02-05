package models

import "time"

type Leaderboard struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	TotalWin  float64   `json:"total_win" db:"total_win"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
