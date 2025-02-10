package servers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// SpinWheelHandler - API обработчик вращения
func (s *Server) SpinWheelHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Преобразуем userID в int
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	// Проверяем, играл ли уже сегодня
	var count int
	err = s.db.QueryRow("SELECT COUNT(*) FROM user_wheel_spins WHERE user_id = $1 AND spin_time >= NOW() - INTERVAL '1 day'", userID).Scan(&count)
	if err == nil && count > 0 {
		http.Error(w, "Вы уже крутили колесо сегодня!", http.StatusForbidden)
		return
	}

	// Получаем случайный приз
	reward, err := s.GetRandomReward()
	if err != nil {
		http.Error(w, "Ошибка выбора приза", http.StatusInternalServerError)
		return
	}

	// Сохраняем результат в БД
	err = s.SaveWheelSpin(userID, reward.ID)
	if err != nil {
		http.Error(w, "Ошибка сохранения результата", http.StatusInternalServerError)
		return
	}

	// Начисляем бонус
	err = s.ApplyReward(userID, reward)
	if err != nil {
		http.Error(w, "Ошибка начисления бонуса", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(reward)
}
