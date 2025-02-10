package servers

import (
	"Aitu-Bet/internal/models"
	"math/rand"
	"time"
)

// GetWheelRewards - Получение списка призов
func (s *Server) GetWheelRewards() ([]models.WheelReward, error) {
	rows, err := s.db.Query("SELECT id, reward_name, reward_type, reward_value, probability FROM wheel_rewards")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rewards []models.WheelReward
	for rows.Next() {
		var r models.WheelReward
		if err := rows.Scan(&r.ID, &r.RewardName, &r.RewardType, &r.RewardValue, &r.Probability); err != nil {
			return nil, err
		}
		rewards = append(rewards, r)
	}
	return rewards, nil
}

// GetRandomReward - Выбор случайного приза
func (s *Server) GetRandomReward() (models.WheelReward, error) {
	rewards, err := s.GetWheelRewards()
	if err != nil {
		return models.WheelReward{}, err
	}

	rand.Seed(time.Now().UnixNano())
	randomVal := rand.Float64() * 100
	cumulativeProbability := 0.0

	for _, reward := range rewards {
		cumulativeProbability += reward.Probability
		if randomVal <= cumulativeProbability {
			return reward, nil
		}
	}
	return rewards[len(rewards)-1], nil
}

// SaveWheelSpin - Сохранение результата вращения
func (s *Server) SaveWheelSpin(userID, rewardID int) error {
	_, err := s.db.Exec("INSERT INTO user_wheel_spins (user_id, reward_id) VALUES ($1, $2)", userID, rewardID)
	return err
}

// ApplyReward - Начисление бонуса пользователю
func (s *Server) ApplyReward(userID int, reward models.WheelReward) error {
	if reward.RewardType == "bonus_money" {
		_, err := s.db.Exec("UPDATE users SET Wincash = Wincash + $1 WHERE id = $2", reward.RewardValue, userID)
		return err
	}
	return nil
}
