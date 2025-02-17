package servers

import (
	"Aitu-Bet/internal/models"
	"Aitu-Bet/utils"
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
		"INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id, username, email, role, cash",
		username, email, hashedPassword, role).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Cash)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Server) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow("SELECT id, username, email, password, role, cash FROM users WHERE email = $1", email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.Cash)

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

	rows, err := s.db.Query("SELECT id, username, email, role, cash FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Cash); err != nil {
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

func (s *Server) FindUserByToken(tokenString string) (*models.User, error) {
	// Validate the JWT token
	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		return nil, err
	}

	// Extract email from claims
	email := claims.Email

	// Fetch the user by email
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
