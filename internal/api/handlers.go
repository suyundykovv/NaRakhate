package api

import (
	"Aitu-Bet/internal/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// LeaderboardHandler — обработчик для работы с лидербордом
type LeaderboardHandler struct {
	Service *services.LeaderboardService
}

// GetTopPlayers возвращает топ игроков по выигрышу
func (h *LeaderboardHandler) GetTopPlayers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10") // Получаем limit из запроса, по умолчанию 10
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	log.Printf("🔍 Received request to fetch top %d players", limit) // Логируем запрос на топ игроков

	// Вызов сервиса для получения топ игроков
	players, err := h.Service.GetTopPlayers(limit)
	if err != nil {
		log.Printf("❌ Error in Service.GetTopPlayers: %v", err) // Логируем ошибку сервиса
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaderboard"})
		return
	}

	// Логируем успешный ответ
	log.Printf("✅ Successfully fetched %d players", len(players))
	c.JSON(http.StatusOK, players)
}

// Регистрация маршрутов
func RegisterRoutes(router *gin.Engine, leaderboardService *services.LeaderboardService) {
	leaderboardHandler := &LeaderboardHandler{Service: leaderboardService}

	// Защищенный эндпоинт (пример)
	router.GET("/protected", ProtectedEndpoint)

	// Эндпоинт для лидерборда
	router.GET("/leaderboard", leaderboardHandler.GetTopPlayers)
}

// Защищенный эндпоинт, пример работы с userID
func ProtectedEndpoint(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Welcome authorized user", "user_id": userID})
}
