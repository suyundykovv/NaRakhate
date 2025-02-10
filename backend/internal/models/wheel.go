package models

type WheelReward struct {
	ID          int     `json:"id"`
	RewardName  string  `json:"reward_name"`
	RewardType  string  `json:"reward_type"`
	RewardValue int     `json:"reward_value"`
	Probability float64 `json:"probability"`
}

type UserWheelSpin struct {
	ID       int `json:"id"`
	UserID   int `json:"user_id"`
	RewardID int `json:"reward_id"`
}
