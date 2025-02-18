package dal

import (
	"Aitu-Bet/internal/models"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(username, email, password, role string, db *sql.DB) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = db.QueryRow(
		"INSERT INTO users (username, email, password, role, cash) VALUES ($1, $2, $3, $4, $5) RETURNING id, username, email, role, cash",
		username, email, hashedPassword, role, 1000).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Cash)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(email string, db *sql.DB) (*models.User, error) {
	var user models.User
	err := db.QueryRow("SELECT id, username, email, password, role, cash FROM users WHERE email = $1", email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.Cash)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func VerifyPassword(hashedPassword, password string, db *sql.DB) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func ReadAllUsers(db *sql.DB) ([]models.User, error) {
	var users []models.User

	rows, err := db.Query("SELECT id, username, email, role, cash FROM users")
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

func DeleteUserData(id string, db *sql.DB) error {
	_, err := db.Exec(`DELETE FROM users WHERE "id" = $1`, id)
	return err
}
