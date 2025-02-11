package models

import "time"

type SlotSpin struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	BetAmount float64   `json:"bet_amount"`
	Reel1     string    `json:"reel1"`
	Reel2     string    `json:"reel2"`
	Reel3     string    `json:"reel3"`
	Payout    float64   `json:"payout"`
	CreatedAt time.Time `json:"created_at"`
}
