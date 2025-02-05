package models

type Event struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	Category    string `json:"category"`
}
