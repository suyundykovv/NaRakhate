package services

import (
	"Aitu-Bet/internal/models"
	"database/sql"
	_ "errors"
	"golang.org/x/crypto/bcrypt"
)

// UserService — сервис для работы с пользователями
type UserService struct {
	DB *sql.DB // Подключение к базе данных
}

// Хеширование пароля
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Создание нового пользователя
func (service *UserService) CreateUser(username, email, password, role string) (*models.User, error) {
	// Хешируем пароль
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Запрос в базу данных для добавления пользователя
	query := `INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id, username, email, role, created_at, updated_at`
	user := &models.User{}
	err = service.DB.QueryRow(query, username, email, hashedPassword, role).
		Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}
