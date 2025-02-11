package servers

import (
	"Aitu-Bet/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
)

// SpinWheelHandler - API обработчик вращения колеса фортуны (платного)
func (s *Server) SpinWheelHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем user_id из query параметров
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Преобразуем user_id в int
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	// Получаем данные пользователя
	var user models.User
	err = s.db.QueryRow("SELECT id, username, wincash FROM users WHERE id = $1", userID).
		Scan(&user.ID, &user.Username, &user.Wincash)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Устанавливаем стоимость вращения колеса
	spinCost := 50.0

	// Проверяем баланс игрока
	if user.Wincash < spinCost {
		http.Error(w, "Not enough balance to spin", http.StatusForbidden)
		return
	}

	// Списываем стоимость спина перед вращением
	_, err = s.db.Exec("UPDATE users SET wincash = wincash - $1 WHERE id = $2", spinCost, userID)
	if err != nil {
		http.Error(w, "Failed to deduct spin cost", http.StatusInternalServerError)
		return
	}

	// Получаем случайный приз
	reward, err := s.GetRandomReward()
	if err != nil {
		http.Error(w, "Ошибка выбора приза", http.StatusInternalServerError)
		return
	}

	// Обновляем данные пользователя (добавляем выигрыш)
	err = s.UpdateUserAfterSpin(userID, reward)
	if err != nil {
		http.Error(w, "Ошибка обновления пользователя", http.StatusInternalServerError)
		return
	}

	// Отправляем JSON-ответ пользователю
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reward)
}
