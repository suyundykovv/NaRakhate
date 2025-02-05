package models

import "time"

// Root struct for the response
type FootballResponse struct {
	Get        string     `json:"get"`
	Parameters Parameters `json:"parameters"`
	Errors     []string   `json:"errors"`
	Results    int        `json:"results"`
	Paging     Paging     `json:"paging"`
	Response   []Fixture  `json:"response"`
}

type Parameters struct {
	League string `json:"league"`
	Season string `json:"season"`
}

type Paging struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}

type Fixture struct {
	FixtureDetails FixtureDetails `json:"fixture"`
	League         League         `json:"league"`
	Teams          Teams          `json:"teams"`
	Goals          Goals          `json:"goals"`
	Score          Score          `json:"score"`
}

type FixtureDetails struct {
	ID        int       `json:"id"`
	Referee   string    `json:"referee"`
	Timezone  string    `json:"timezone"`
	Date      time.Time `json:"date"`
	Timestamp int64     `json:"timestamp"`
	Periods   Periods   `json:"periods"`
	Venue     Venue     `json:"venue"`
	Status    Status    `json:"status"`
}

type Periods struct {
	First  int64 `json:"first"`
	Second int64 `json:"second"`
}

type Venue struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	City string `json:"city"`
}

type Status struct {
	Long    string `json:"long"`
	Short   string `json:"short"`
	Elapsed int    `json:"elapsed"`
	Extra   *int   `json:"extra,omitempty"`
}

type League struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Country   string `json:"country"`
	Logo      string `json:"logo"`
	Flag      string `json:"flag"`
	Season    int    `json:"season"`
	Round     string `json:"round"`
	Standings bool   `json:"standings"`
}

type Teams struct {
	Home Team `json:"home"`
	Away Team `json:"away"`
}

type Team struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Logo   string `json:"logo"`
	Winner bool   `json:"winner"`
}

type Goals struct {
	Home int `json:"home"`
	Away int `json:"away"`
}

type Score struct {
	Halftime  ScoreDetails  `json:"halftime"`
	Fulltime  ScoreDetails  `json:"fulltime"`
	Extratime *ScoreDetails `json:"extratime,omitempty"`
	Penalty   *ScoreDetails `json:"penalty,omitempty"`
}

type ScoreDetails struct {
	Home int `json:"home"`
	Away int `json:"away"`
}
