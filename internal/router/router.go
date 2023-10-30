package router

import (
	"bytes"
	"encoding/json"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/KoLLlaka/poma-botv2.0/internal/logging"
	"github.com/KoLLlaka/poma-botv2.0/internal/model"
	"github.com/KoLLlaka/poma-botv2.0/internal/playlist"

	"github.com/KoLLlaka/poma-botv2.0/internal/db"

	"github.com/gorilla/websocket"
)

const (
	AUG   = "aug"
	MUSIC = "music"
)

var (
	tmpl     *template.Template
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Пропускаем любой запрос
		},
	}
	augFiles []string
)

func init() {
	tmpl = template.Must(template.ParseGlob("template/*.html"))

	if err := os.MkdirAll("static/aug", os.FileMode(0644)); err != nil {
		panic(err)
	}

	files, err := os.ReadDir("./static/aug")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		augFiles = append(augFiles, file.Name())
	}
}

type Server struct {
	logger     logging.Logger
	musicStore db.MusicStore

	conf          *model.Config
	clients       map[string]*websocket.Conn
	handleMessage func(message []byte)

	augChan    chan string
	musicChan  chan model.Playlist
	myPlaylist []model.Playlist
}

func New(logger logging.Logger, musicStore db.MusicStore, conf *model.Config,
	augChan chan string, musicChan chan model.Playlist,
	myPlaylist []model.Playlist,
) *Server {
	return &Server{
		logger:     logger,
		musicStore: musicStore,
		conf:       conf,
		clients:    make(map[string]*websocket.Conn),
		handleMessage: func(message []byte) {
			logger.Infof("[message from socket] %s", message)
		},
		augChan:    augChan,
		musicChan:  musicChan,
		myPlaylist: myPlaylist,
	}
}

func (s *Server) Start() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/"+AUG, s.aug)
	router.HandleFunc("/"+AUG+"/ws", s.augws)

	router.HandleFunc("/"+MUSIC, s.music)
	router.HandleFunc("/"+MUSIC+"/ws", s.musicws)
	router.HandleFunc("/api/yplaylist", s.yplaylist)

	return router
}

// augury
func (s *Server) aug(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, AUG+".html", nil)
}
func (s *Server) augws(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()
	s.clients[AUG] = conn
	defer delete(s.clients, AUG)

	type ReqAug struct {
		Name string `json:"name,omitempty"`
		Link string `json:"link,omitempty"`
	}

	go func() {
		for {
			mt, message, err := conn.ReadMessage()

			if err != nil || mt == websocket.CloseMessage {
				break // Выходим из цикла, если клиент пытается закрыть соединение или связь прервана
			}

			go s.handleMessage(message)
		}
	}()

	for {
		name := <-s.augChan

		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		num := r1.Intn(len(augFiles))
		link := augFiles[num]

		augMsg := ReqAug{
			Name: name,
			Link: "./static/aug/" + link,
		}

		var network bytes.Buffer
		enc := json.NewEncoder(&network)
		err := enc.Encode(augMsg)
		if err != nil {
			s.logger.Errorln(err)
			return
		}

		s.writeByteMsg(AUG, network.Bytes())

		time.Sleep(10 * time.Second)
	}
}

// music
func (s *Server) music(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, MUSIC+".html", nil)
}
func (s *Server) musicws(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()
	s.clients[MUSIC] = conn
	defer delete(s.clients, MUSIC)

	go func() {
		for {
			mt, message, err := conn.ReadMessage()

			if err != nil || mt == websocket.CloseMessage {
				s.logger.Errorf("[error from socket] %+v", err)

				break // Выходим из цикла, если клиент пытается закрыть соединение или связь прервана
			}

			data := model.MsgFromSocket{}
			if err := json.Unmarshal(message, &data); err != nil {
				s.logger.Errorf("[error from socket] %+v", err)
			}

			switch data.Reason {
			case "addDuration":
				s.musicStore.StoreDuration(&data.Song)
			}

			go s.handleMessage(message)
		}
	}()

	go func() {
		for _, song := range s.myPlaylist {
			s.musicStore.GetDuration(&song)

			var network bytes.Buffer
			enc := json.NewEncoder(&network)
			err := enc.Encode(song)
			if err != nil {
				s.logger.Errorln(err)

				return
			}

			s.writeByteMsg(MUSIC, network.Bytes())
		}
	}()

	for {
		musicStruct := <-s.musicChan

		var network bytes.Buffer
		enc := json.NewEncoder(&network)
		err := enc.Encode(musicStruct)
		if err != nil {
			s.logger.Errorln(err)

			return
		}

		s.writeByteMsg(MUSIC, network.Bytes())
	}
}
func (s *Server) yplaylist(w http.ResponseWriter, r *http.Request) {
	type Playlist struct {
		Link string `json:"link,omitempty"`
	}

	w.Header().Set("Content-Type", "application/json")
	resp := Playlist{}
	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		s.logger.Errorln("[error from music]", err)
	}
	s.logger.Infoln("[info] request from server:", resp.Link)

	songs := playlist.ListOfSongsFromPlaylist(resp.Link, s.conf.YoutubeKey, "")
	var songsLink []string
	listOfSongs := []*model.Playlist{}

	for i, song := range songs {
		songsLink = append(songsLink, song.Link)

		if (i+1)%10 == 0 {
			listSongsInfo := playlist.ReqSongInfo(strings.Join(songsLink, ","), s.conf.YoutubeKey)
			listOfSongs = append(listOfSongs, listSongsInfo...)

			songsLink = []string{}
		}

		if i+1 == len(songs) {
			listSongsInfo := playlist.ReqSongInfo(strings.Join(songsLink, ","), s.conf.YoutubeKey)
			listOfSongs = append(listOfSongs, listSongsInfo...)
		}
	}

	json.NewEncoder(w).Encode(listOfSongs)
}

// sss
func (s *Server) writeByteMsg(typeMsg string, message []byte) {
	conn := s.clients[typeMsg]
	conn.WriteMessage(websocket.TextMessage, message)
}
