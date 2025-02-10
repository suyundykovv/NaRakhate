package servers

import (
	"Aitu-Bet/internal/models"
	"math/rand"
	"time"
)

// Получаем список призов из БД
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

// Выбираем случайный приз с учетом вероятностей
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

// Обновляем данные пользователя после спина
func (s *Server) UpdateUserAfterSpin(userID int, reward models.WheelReward) error {
	_, err := s.db.Exec(
		"UPDATE users SET last_spin_time = NOW(), spin_count = spin_count + 1, wincash = wincash + $1 WHERE id = $2",
		reward.RewardValue, userID,
	)
	return err
}
