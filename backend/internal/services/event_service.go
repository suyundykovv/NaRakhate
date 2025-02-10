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

func UpdateOddsDuringMatch(fixture models.Fixture, db *sql.DB) error {
	elapsedTime := fixture.FixtureDetails.Status.Elapsed
	homeGoals := fixture.Goals.Home
	awayGoals := fixture.Goals.Away

	homeWinRate, awayWinRate := calculateDynamicWinRates(homeGoals, awayGoals, elapsedTime)
	totalWinRate := homeWinRate + awayWinRate

	var homeOdds, awayOdds, drawOdds float64
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

	drawOdds = 3.0

	updateQuery := `
		UPDATE events 
		SET home_win_odds = $1, away_win_odds = $2, draw_odds = $3 
		WHERE id = $4
	`
	_, err := db.Exec(updateQuery, homeOdds, awayOdds, drawOdds, fixture.FixtureDetails.ID)
	if err != nil {
		log.Println("Error updating odds in the database:", err)
		return err
	}

	log.Printf("Updated odds for match %d: home_odds=%.2f, away_odds=%.2f, draw_odds=%.2f",
		fixture.FixtureDetails.ID, homeOdds, awayOdds, drawOdds)

	return nil
}

func calculateDynamicWinRates(homeGoals, awayGoals, elapsedTime int) (float64, float64) {
	homeWinRate := 0.5
	awayWinRate := 0.5

	goalDifference := homeGoals - awayGoals
	if goalDifference > 0 {
		homeWinRate += 0.1 * float64(goalDifference)
		awayWinRate -= 0.1 * float64(goalDifference)
	} else if goalDifference < 0 {
		awayWinRate += 0.1 * float64(-goalDifference)
		homeWinRate -= 0.1 * float64(-goalDifference)
	}

	timeFactor := float64(elapsedTime) / 90.0
	if goalDifference > 0 {
		homeWinRate += 0.2 * timeFactor
		awayWinRate -= 0.2 * timeFactor
	} else if goalDifference < 0 {
		awayWinRate += 0.2 * timeFactor
		homeWinRate -= 0.2 * timeFactor
	}

	if homeWinRate < 0.1 {
		homeWinRate = 0.1
	}
	if awayWinRate < 0.1 {
		awayWinRate = 0.1
	}

	return homeWinRate, awayWinRate
}
