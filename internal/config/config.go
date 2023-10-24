package config

import (
	"os"
	"strings"

	"github.com/KoLLlaka/poma-botv2.0/internal/logging"
	"github.com/KoLLlaka/poma-botv2.0/internal/model"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

const (
	HOST       = "HOST"
	PORT       = "PORT"
	USERID     = "USERID"
	YOUTUBEKEY = "YOUTUBEKEY"
	AUDIOPATH  = "AUDIOPATH"
)

var logger logging.Logger

// loads values from .env into the system
func init() {
	logger = logging.GetLogger()
	if err := godotenv.Load(); err != nil {
		logger.Fatalln("No .env file found")
	}
}

// loads values from .env to Config
func NewConfig() *model.Config {
	yamlFile, err := os.ReadFile("./config.yaml")
	yamlConf := &model.YAMLConfig{}
	if err != nil {
		logger.Fatalln("No .yaml file found")
	}
	err = yaml.Unmarshal(yamlFile, yamlConf)
	if err != nil {
		logger.Fatalln(err)
	}

	rewards := make(map[string]model.ConfigReward)
	for _, reward := range yamlConf.Rewards {
		name := strings.ToLower(reward.RewardName)

		if reward.Duration == 0 && reward.RewardType == "music" {
			reward.Duration = 600
		}

		rewards[name] = reward
	}

	return &model.Config{
		Host:       getEnv(HOST, "localhost"),
		Port:       getEnv(PORT, "8080"),
		UserID:     getEnv(USERID, ""),
		YoutubeKey: getEnv(YOUTUBEKEY, ""),
		AudioPath:  getEnv(AUDIOPATH, "./static/playlist"),
		Rewards:    rewards,
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
