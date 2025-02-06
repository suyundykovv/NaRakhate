package api

import (
	"Aitu-Bet/internal/db"       // Импортируем для доступа к глобальной базе данных (если нужно)
	"Aitu-Bet/internal/services" // Импортируем пакет, где определен UserService
	"encoding/json"
	"log"
	"net/http"
)

// CreateUserHandler — обработчик для создания пользователя
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	// Чтение данных из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Создаем экземпляр UserService, передаем подключение к базе данных
	userService := &services.UserService{DB: db.DB} // db.DB — глобальная переменная для подключения к базе данных

	// Создание пользователя через метод CreateUser
	user, err := userService.CreateUser(userInput.Username, userInput.Email, userInput.Password, userInput.Role)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Ответ с созданным пользователем
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
