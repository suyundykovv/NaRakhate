package services

import (
	"Aitu-Bet/internal/models"
	"database/sql"
	"log"
)

func winRateService(teamID int, db *sql.DB) (float64, error) {
	var winCount int
	var totalMatches int

	rows, err := db.Query(`
		SELECT 
			CASE 
				WHEN (home_team_id = $1 AND home_goals > away_goals) OR (away_team_id = $1 AND away_goals > home_goals) THEN 1 
				ELSE 0 
			END AS win
		FROM matches
		WHERE (home_team_id = $1 OR away_team_id = $1)
		ORDER BY match_date DESC
		LIMIT 10
	`, teamID)
	if err != nil {
		log.Println("Error fetching matches for team:", err)
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var win int
		if err := rows.Scan(&win); err != nil {
			log.Println("Error scanning match result:", err)
			return 0, err
		}
		winCount += win
		totalMatches++
	}

	if totalMatches == 0 {
		return 0, nil
	}

	winRate := float64(winCount) / float64(totalMatches)
	return winRate, nil
}

func CreateOddsService(fixture models.Fixture, db *sql.DB) (models.Odds, error) {
	homeTeamID := fixture.Teams.Home.ID
	awayTeamID := fixture.Teams.Away.ID

	homeWinRate, err := winRateService(homeTeamID, db)
	if err != nil {
		return models.Odds{}, err
	}

	awayWinRate, err := winRateService(awayTeamID, db)
	if err != nil {
		return models.Odds{}, err
	}

	if homeWinRate == 0 && awayWinRate == 0 {
		homeWinRate = 0.5
		awayWinRate = 0.5
	}
	totalWinRate := homeWinRate + awayWinRate

	var homeOdds, awayOdds float64
	if homeWinRate == 0 {
		homeOdds = 3.0
	} else {
		homeOdds = 1 / (homeWinRate / totalWinRate)
	}

	if awayWinRate == 0 {
		awayOdds = 3.0
	} else {
		awayOdds = 1 / (awayWinRate / totalWinRate)
	}

	drawOdds := 3.0

	return models.Odds{
		HomeWin: homeOdds,
		AwayWin: awayOdds,
		Draw:    drawOdds,
	}, nil
}
