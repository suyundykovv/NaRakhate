package servers

import (
	"Aitu-Bet/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// SpinWheelHandler - API обработчик вращения
func (s *Server) SpinWheelHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	// Получаем данные пользователя
	var user models.User
	err = s.db.QueryRow("SELECT id, username, wincash, email, password, role, last_spin_time, spin_count FROM users WHERE id = $1", userID).
		Scan(&user.ID, &user.Username, &user.Wincash, &user.Email, &user.Password, &user.Role, &user.LastSpinTime, &user.SpinCount)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Проверяем, крутил ли пользователь колесо сегодня
	if time.Since(user.LastSpinTime).Hours() < 24 {
		http.Error(w, "Вы уже крутили колесо сегодня!", http.StatusForbidden)
		return
	}

	// Получаем случайный приз
	reward, err := s.GetRandomReward()
	if err != nil {
		http.Error(w, "Ошибка выбора приза", http.StatusInternalServerError)
		return
	}

	// Обновляем данные пользователя
	err = s.UpdateUserAfterSpin(userID, reward)
	if err != nil {
		http.Error(w, "Ошибка обновления пользователя", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(reward)
}
