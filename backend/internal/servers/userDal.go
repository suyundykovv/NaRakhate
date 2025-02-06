package servers

import (
	"Aitu-Bet/internal/models"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) createUser(username, email, password, role string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = s.db.QueryRow(
		"INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id, username, email, role",
		username, email, hashedPassword, role).Scan(&user.ID, &user.Username, &user.Email, &user.Role)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Server) getUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow("SELECT id, username, email, password, role FROM users WHERE email = $1", email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Server) verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *Server) ReadAllUsers() ([]models.User, error) {
	var users []models.User

	rows, err := s.db.Query("SELECT id, username, email, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (s *Server) deleteUserData(id string) error {
	_, err := s.db.Exec(`DELETE FROM users WHERE "id" = $1`, id)
	return err
}
func (s *Server) GetTableUsers() ([]models.User, error) {
	var users []models.User

	rows, err := s.db.Query("SELECT id, username, Wincash FROM users ORDER BY Wincash DESC ")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Wincash); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}
