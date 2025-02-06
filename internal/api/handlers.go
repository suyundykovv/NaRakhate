package api

import (
	"Aitu-Bet/internal/services"
	"encoding/json"
	"net/http"
)

// ProtectedHandler — обработчик для защищенного маршрута
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	// Для примера, просто возвращаем статус 200 OK
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Protected content"})
}

// GetTopPlayersHandler — обработчик для получения топ-игроков
func GetTopPlayersHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем топ-10 игроков
	players, err := services.GetTopPlayers(10)
	if err != nil {
		http.Error(w, "Error getting top players", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(players); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
