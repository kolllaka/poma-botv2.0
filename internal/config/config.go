package config

import (
	"log"
	"os"
	"strings"

	"github.com/KoLLlaka/__augury/internal/model"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

const (
	HOST       = "HOST"
	PORT       = "PORT"
	USERNAM    = "USERNAM"
	USERID     = "USERID"
	BOTNAME    = "BOTNAME"
	BOTAUTH    = "BOTAUTH"
	YOUTUBEKEY = "YOUTUBEKEY"
	AUDIOPATH  = "AUDIOPATH"
)

// loads values from .env into the system
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

// loads values from .env to Config
func NewConfig() *model.Config {
	yamlFile, err := os.ReadFile("./config.yaml")
	yamlConf := &model.YAMLConfig{}
	if err != nil {
		log.Fatal("No .yaml file found")
	}
	err = yaml.Unmarshal(yamlFile, yamlConf)
	if err != nil {
		log.Fatal(err)
	}

	rewards := make(map[string]model.ConfigReward)
	for _, reward := range yamlConf.Rewards {
		name := strings.ToLower(reward.RewardName)

		if reward.Duration == 0 && reward.RewardType == "music" {
			reward.Duration = 600
		}

		if reward.Duration == 0 && reward.RewardType == "ban" {
			reward.Duration = 60
		}

		rewards[name] = reward
	}

	userName := strings.ToLower(getEnv(USERNAM, ""))
	log.Printf("from config: %s\n", userName)

	return &model.Config{
		Host:       getEnv(HOST, "localhost"),
		Port:       getEnv(PORT, "8080"),
		UserName:   userName,
		UserID:     getEnv(USERID, ""),
		YoutubeKey: getEnv(YOUTUBEKEY, ""),
		BotName:    strings.ToLower(getEnv(BOTNAME, "")),
		BotAuth:    getEnv(BOTAUTH, ""),
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
