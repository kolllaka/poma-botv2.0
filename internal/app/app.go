package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Adeithe/go-twitch"

	"github.com/KoLLlaka/poma-botv2.0/internal/config"
	"github.com/KoLLlaka/poma-botv2.0/internal/logging"
	"github.com/KoLLlaka/poma-botv2.0/internal/model"
	"github.com/KoLLlaka/poma-botv2.0/internal/playlist"
	"github.com/KoLLlaka/poma-botv2.0/internal/router"
)

const (
	MESSAGE = "message"
	TWITCH  = "twitch"
)

// 28140745 yoburg
var (
	conf *model.Config

	auguryChannel chan string         = make(chan string)
	musicChanel   chan model.Playlist = make(chan model.Playlist)

	myPlaylist []model.Playlist
	logger     logging.Logger
)

func StartApp() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	// logger
	logger = logging.GetLogger()

	// config
	conf = config.NewConfig()

	// load playlist
	myPlaylist = playlist.LoadMyPlaylist(conf.AudioPath)

	// connect to PubSub
	ps := twitch.PubSub()
	ps.OnShardMessage(onMessage)
	ps.Listen("community-points-channel-v1", conf.UserID)

	// Start Server
	server := router.New(logger, conf, auguryChannel, musicChanel, myPlaylist)
	router := server.Start()

	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	router.Handle("/audio/", http.StripPrefix("/audio/", http.FileServer(http.Dir(conf.AudioPath))))
	go http.ListenAndServe(":"+conf.Port, router)
	logger.Infof("server start on port :%s\n", conf.Port)

	<-sc
	logger.Infoln("Stopping...")
	ps.Close()
}

func onMessage(shard int, topic string, data []byte) {
	reward := model.Reward{}
	json.Unmarshal(data, &reward)

	if reward.Type == "reward-redeemed" {
		title := strings.ToLower(reward.Data.Redemption.Reward.Title)

		confReward, ok := conf.Rewards[title]
		if !ok {
			logger.Tracef("команда не зарегистрирована %s\n", title)

			return
		}

		name := reward.Data.Redemption.User.DisplayName

		switch confReward.RewardType {
		case "augury":
			logger.Infof("%s заказал %s\n", name, confReward.RewardName)

			auguryChannel <- fmt.Sprintf(confReward.RewardTitle, name)
		case "music":
			music := reward.Data.Redemption.UserInput
			logger.Infof("%s заказал %s с текстом: %s\n", name, confReward.RewardName, music)

			songUrl, err := playlist.SongRequest(music)
			if err != nil {
				logger.Errorln(err)

				return
			}

			songsInfo := playlist.ReqSongInfo(songUrl, conf.YoutubeKey)
			for _, songInfo := range songsInfo {
				songInfo.IsReward = true
				songInfo.Name = name
				logger.Infoln(*songInfo)
				if songInfo.Duration > confReward.Duration {
					logger.Infof("request слишком длинный: %dc а должен быть меньше %dc\n", songInfo.Duration, confReward.Duration)

					return
				}

				musicChanel <- *songInfo
			}
		}

		return
	}
}
