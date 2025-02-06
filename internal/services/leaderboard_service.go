package services

import (
	"Aitu-Bet/internal/db"
	"Aitu-Bet/internal/models"
	"log"
)

// Функция для получения топ игроков
func GetTopPlayers(limit int) ([]models.Player, error) {
	query := `SELECT id, username, total_winnings FROM players ORDER BY total_winnings DESC LIMIT $1`
	rows, err := db.DB.Query(query, limit) // Используем db.DB
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

// AddPlayer добавляет нового игрока в базу данных
func AddPlayer(name string, winnings float64) error {
	query := `INSERT INTO players (username, total_winnings) VALUES ($1, $2)`
	_, err := db.DB.Exec(query, name, winnings) // Вставляем игрока в таблицу
	if err != nil {
		log.Printf("Error adding player: %v", err)
		return err
	}

	log.Printf("Player added: %s, winnings: %.2f", name, winnings) // Логируем добавление игрока
	return nil
}
