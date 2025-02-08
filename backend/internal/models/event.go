package models

type Event struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	Category    string `json:"category"`
	Odd         Odds   `json:"odds"`
}

type Odds struct {
	HomeWin float64 `json:"home_win"`
	AwayWin float64 `json:"away_win"`
	Draw    float64 `json:"draw"`
}
