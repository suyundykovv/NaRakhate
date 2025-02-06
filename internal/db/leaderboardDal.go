package db

import (
	"Aitu-Bet/internal/models"
	"database/sql"
	"log"
)

// Получаем топ игроков из базы данных
func GetTopPlayers(db *sql.DB, limit int) ([]models.Player, error) {
	query := `SELECT id, username, total_winnings FROM players ORDER BY total_winnings DESC LIMIT $1`
	rows, err := db.Query(query, limit)
	if err != nil {
		log.Printf("Error fetching top players: %v", err)
		return nil, err
	}
	defer rows.Close()

	var players []models.Player
	for rows.Next() {
		var player models.Player
		if err := rows.Scan(&player.ID, &player.Username, &player.TotalWinnings); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		players = append(players, player)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	return players, nil
}
