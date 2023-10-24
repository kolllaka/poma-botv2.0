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
	"time"

	"github.com/Adeithe/go-twitch"
	"github.com/Adeithe/go-twitch/irc"

	"github.com/KoLLlaka/__augury/internal/config"
	"github.com/KoLLlaka/__augury/internal/model"
	"github.com/KoLLlaka/__augury/internal/moderate"
	"github.com/KoLLlaka/__augury/internal/playlist"
	"github.com/KoLLlaka/__augury/internal/router"
)

const (
	MESSAGE = "message"
	TWITCH  = "twitch"

	BotId    = "g3mziksa4gqhlbeh926iuh4uqv230p"
	BotToken = "Bearer xio4ieig0ka569y4j6hakvvi8mp08m"
)

// 28140745 yoburg
var (
	conf *model.Config

	auguryChannel chan string         = make(chan string)
	musicChanel   chan model.Playlist = make(chan model.Playlist)
	banChan       chan model.BanUser  = make(chan model.BanUser)

	myPlaylist []model.Playlist
)

func main() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	// config
	conf = config.NewConfig()

	// load playlist
	myPlaylist = playlist.LoadMyPlaylist(conf.AudioPath)

	// writer to IRC
	writer := &irc.Conn{}
	writer.SetLogin(conf.BotName, conf.BotAuth)
	if err := writer.Connect(); err != nil {
		panic("failed to start writer")
	}

	// moderate
	api := moderate.NewApiClient(
		conf.UserID,
		conf.UserID,
		BotId,
		BotToken,
	)

	// api.GetUserId("bobcehekolliaka")
	log.Printf("[broadcasterId] %s\n", conf.UserID)
	// vips, _ := api.GetVips(conf.UserID, "")
	// log.Printf("%+v\n", vips)
	// moderators, _ := api.GetModerators(conf.UserID, []string{"42879324", "78053207"})
	// log.Printf("%+v\n", moderators)
	subscriptions, _ := api.GetSubscriptions(conf.UserID, nil)
	log.Printf("%+v\n", subscriptions)
	// !
	go func() {
		for {
			ban := <-banChan

			userId, _ := api.GetUserId(ban.UserName)
			api.Ban(userId, ban.Duration, ban.Reason)

			time.Sleep(30 * time.Second)
			api.UnBan(userId)
		}
	}()

	// reader from IRC
	reader := twitch.IRC()
	reader.OnShardReconnect(onShardReconnect)
	reader.OnShardLatencyUpdate(onShardLatencyUpdate)
	reader.OnShardMessage(onShardMessage)
	if err := reader.Join(conf.UserName); err != nil {
		panic(err)
	}
	fmt.Printf("Connecting to IRC to channel %s.....\n", conf.UserName)

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
	fmt.Printf("server start on port :%s\n", conf.Port)

	<-sc
	fmt.Println("Stopping...")
	ps.Close()
	reader.Close()
	writer.Close()
}

func onShardReconnect(shardID int) {
	log.Printf("Shard #%d reconnected\n", shardID)
}
func onShardLatencyUpdate(shardID int, latency time.Duration) {
	log.Printf("Shard #%d has %dms ping\n", shardID, latency.Milliseconds())
}
func onShardMessage(shardID int, msg irc.ChatMessage) {
	log.Printf("%s: %s\n", msg.Sender.DisplayName, msg.Text)

	// msgFromIRC := model.MessageFromIRC{
	// 	TypeOfMessage: MESSAGE,
	// 	Platform:      TWITCH,
	// 	Channel:       msg.Channel,
	// 	ChannelId:     string(msg.ChannelID),
	// 	Sender: model.Sender{
	// 		From:          msg.Sender.DisplayName,
	// 		FromId:        string(msg.Sender.ID),
	// 		Badges:        msg.Sender.Badges,
	// 		BadgesInfo:    msg.Sender.BadgeInfo,
	// 		Color:         msg.Sender.Color,
	// 		IsBroadcaster: msg.Sender.IsBroadcaster,
	// 		IsModerator:   msg.Sender.IsModerator,
	// 		IsVIP:         msg.Sender.IsVIP,
	// 		IsSubscriber:  msg.Sender.IsSubscriber,
	// 	},
	// 	Time:     time.Now(),
	// 	Text:     msg.Text,
	// 	EmoteMap: msg.IRCMessage.Tags["emotes"],
	// }

	// channelToChat <- msgFromIRC
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
		case "text":
			//log.Printf("%s заказал %s\n", name, conf.Rewards[2].RewardName)
		case "ban":
			text := reward.Data.Redemption.UserInput

			name := strings.Fields(text)[0]
			if strings.HasPrefix(name, "@") {
				name = strings.Replace(name, "@", "", 1)
			}
			log.Printf("%s заказал %s\n с текстом: %s", name, confReward.RewardName, text)

			banChan <- model.BanUser{
				UserName: strings.ToLower(name),
				Duration: confReward.Duration,
				Reason:   confReward.RewardTitle,
			}
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
