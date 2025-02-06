package api

import (
	"Aitu-Bet/internal/services"
	"encoding/json"
	"net/http"
)

// AddPlayerHandler обрабатывает запрос на добавление игрока
func AddPlayerHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string  `json:"name"`
		Winnings float64 `json:"winnings"`
	}

	// Декодируем данные из запроса
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Добавляем нового игрока в базу данных
	err := services.AddPlayer(req.Name, req.Winnings)
	if err != nil {
		http.Error(w, "Error adding player", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
