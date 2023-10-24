package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Adeithe/go-twitch"

	"github.com/KoLLlaka/poma-botv2.0/internal/config"
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
)

func main() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	// config
	conf = config.NewConfig()

	// load playlist
	myPlaylist = playlist.LoadMyPlaylist(conf.AudioPath)

	// connect to PubSub
	ps := twitch.PubSub()
	ps.OnShardMessage(onMessage)
	ps.Listen("community-points-channel-v1", conf.UserID)

	// Start Server
	server := router.New(conf, auguryChannel, musicChanel, myPlaylist)
	router := server.Start()

	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	router.Handle("/audio/", http.StripPrefix("/audio/", http.FileServer(http.Dir(conf.AudioPath))))
	go http.ListenAndServe(":"+conf.Port, router)
	log.Printf("server start on port :%s\n", conf.Port)

	<-sc
	log.Println("Stopping...")
	ps.Close()
}

func onMessage(shard int, topic string, data []byte) {
	reward := model.Reward{}
	json.Unmarshal(data, &reward)

	if reward.Type == "reward-redeemed" {
		title := strings.ToLower(reward.Data.Redemption.Reward.Title)

		confReward, ok := conf.Rewards[title]
		if !ok {
			log.Printf("команда не зарегистрирована %s\n", title)

			return
		}

		name := reward.Data.Redemption.User.DisplayName

		switch confReward.RewardType {
		case "augury":
			log.Printf("%s заказал %s\n", name, confReward.RewardName)

			auguryChannel <- fmt.Sprintf(confReward.RewardTitle, name)
		case "music":
			music := reward.Data.Redemption.UserInput
			log.Printf("[info] %s заказал %s с текстом: %s\n", name, confReward.RewardName, music)

			songUrl, err := playlist.SongRequest(music)
			if err != nil {
				log.Println("[error]", err)

				return
			}

			songsInfo := playlist.ReqSongInfo(songUrl, conf.YoutubeKey)
			for _, songInfo := range songsInfo {
				songInfo.IsReward = true
				songInfo.Name = name
				log.Println("[info]", *songInfo)
				if songInfo.Duration > confReward.Duration {
					log.Printf("[warn] request слишком длинный: %dc а должен быть меньше %dc\n", songInfo.Duration, confReward.Duration)

					return
				}

				musicChanel <- *songInfo
			}
		}

		return
	}
}
