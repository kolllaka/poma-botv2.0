package model

type Config struct {
	Port       string
	Host       string
	UserID     string
	YoutubeKey string
	AudioPath  string
	Rewards    map[string]ConfigReward
}

type YAMLConfig struct {
	Rewards []ConfigReward `yaml:"rewards"`
}

type ConfigReward struct {
	RewardType  string `yaml:"type"`
	RewardTitle string `yaml:"rewardTitle"`
	RewardName  string `yaml:"rewardName"`
	Duration    int    `yaml:"duration,omitempty"`
}
