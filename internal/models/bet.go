package models

type Bet struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	EventID   int     `json:"event_id"`
	Amount    float64 `json:"amount"`
	Outcome   string  `json:"outcome"`
	CreatedAt string  `json:"created_at"`
}
