package services

import (
	"Aitu-Bet/internal/db"
	"Aitu-Bet/internal/models"
	"fmt"
	"log"
)

type LeaderboardService struct {
	LeaderboardDAL *db.LeaderboardDAL
}

// GetTopPlayers получает топ игроков
func (s *LeaderboardService) GetTopPlayers(limit int) ([]models.Leaderboard, error) {
	log.Printf("🔍 Calling GetTopPlayers with limit: %d", limit)

	// Вызываем метод DAL для получения топ игроков
	players, err := s.LeaderboardDAL.GetTopPlayers(limit)
	if err != nil {
		log.Printf("❌ Error in LeaderboardService.GetTopPlayers: %v", err)
		return nil, fmt.Errorf("failed to fetch leaderboard: %w", err)
	}

	log.Printf("✅ Fetched %d players", len(players))
	return players, nil
}
