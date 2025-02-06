package models

type User struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Wincash  float64 `json:"wincash"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Role     string  `json:"role"`
}
