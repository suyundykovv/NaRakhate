package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Wincash      float64   `json:"wincash"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Role         string    `json:"role"`
	LastSpinTime time.Time `json:"last_spin_time"` // Last spin time
	SpinCount    int       `json:"spin_count"`     // Amount of spins
}
