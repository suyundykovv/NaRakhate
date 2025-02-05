package db

import (
	"Aitu-Bet/internal/models"
	"database/sql"
	"fmt"
	"log"
)

type LeaderboardDAL struct {
	DB *sql.DB
}

// GetTopPlayers получает топ игроков по выигрышу
func (dal *LeaderboardDAL) GetTopPlayers(limit int) ([]models.Leaderboard, error) {
	query := `
        SELECT id, user_id, total_win, updated_at 
        FROM leaderboard
        ORDER BY total_win DESC
        LIMIT $1;
    `
	log.Printf("🔍 Executing SQL query: %s with limit = %d", query, limit)

	rows, err := dal.DB.Query(query, limit)
	if err != nil {
		log.Printf("❌ SQL Error: %v", err)
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}
	defer rows.Close()

	var leaderboard []models.Leaderboard
	for rows.Next() {
		var entry models.Leaderboard
		if err := rows.Scan(&entry.ID, &entry.UserID, &entry.TotalWin, &entry.UpdatedAt); err != nil {
			log.Printf("❌ Row Scan Error: %v", err)
			return nil, err
		}
		leaderboard = append(leaderboard, entry)
	}

	if err := rows.Err(); err != nil {
		log.Printf("❌ Row Iteration Error: %v", err)
		return nil, err
	}

	log.Printf("✅ Successfully fetched %d leaderboard records", len(leaderboard))
	return leaderboard, nil
}

// UpdateLeaderboard обновляет сумму выигрыша игрока
func (dal *LeaderboardDAL) UpdateLeaderboard(userID int, winAmount float64) error {
	query := `
        INSERT INTO leaderboard (user_id, total_win, updated_at)
        VALUES ($1, $2, NOW())
        ON CONFLICT (user_id)
        DO UPDATE SET total_win = leaderboard.total_win + EXCLUDED.total_win, updated_at = NOW();
    `
	log.Printf("🔍 Updating leaderboard for user %d with win amount: %.2f", userID, winAmount)

	_, err := dal.DB.Exec(query, userID, winAmount)
	if err != nil {
		log.Printf("❌ SQL Error (UpdateLeaderboard): %v", err)
		return fmt.Errorf("failed to update leaderboard: %w", err)
	}

	log.Printf("✅ Successfully updated leaderboard for user %d", userID)
	return nil
}
