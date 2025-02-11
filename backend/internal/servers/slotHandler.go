package servers

import (
	"encoding/json"
	"log"
	"net/http"
)

// SpinSlotHandler - API для вращения слотов
func (s *Server) SpinSlotHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID    int     `json:"user_id"`
		BetAmount float64 `json:"bet_amount"`
	}

	// Разбираем JSON-запрос
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Ошибка парсинга JSON: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Обрабатываем вращение
	spin, err := s.SpinSlots(request.UserID, request.BetAmount)
	if err != nil {
		log.Printf("Ошибка во время спина: %v", err)
		http.Error(w, "Spin failed", http.StatusInternalServerError)
		return
	}

	// Отправляем JSON-ответ пользователю
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spin)
}
