package router

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/KoLLlaka/poma-botv2.0/internal/model"
	"github.com/KoLLlaka/poma-botv2.0/internal/playlist"

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

	files, err := ioutil.ReadDir("./static/aug")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		augFiles = append(augFiles, file.Name())
	}
}

type Server struct {
	conf          *model.Config
	clients       map[string]*websocket.Conn
	handleMessage func(message []byte)

	augChan    chan string
	musicChan  chan model.Playlist
	myPlaylist []model.Playlist
}

func New(conf *model.Config,
	augChan chan string, musicChan chan model.Playlist,
	myPlaylist []model.Playlist,
) *Server {
	return &Server{
		conf:    conf,
		clients: make(map[string]*websocket.Conn),
		handleMessage: func(message []byte) {
			log.Printf("[message] %s", message)
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
	router.HandleFunc("/api/playlist", s.playlist)

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
		// fmt.Println(num)
		link := augFiles[num]

		// log.Println(link)
		augMsg := ReqAug{
			Name: name,
			Link: "./static/aug/" + link,
		}

		var network bytes.Buffer
		enc := json.NewEncoder(&network)
		err := enc.Encode(augMsg)
		if err != nil {
			log.Println(err)
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

	type ReqAug struct {
		IsReward bool   `json:"isreward,omitempty"`
		Name     string `json:"name,omitempty"`
		Title    string `json:"title,omitempty"`
		Link     string `json:"link,omitempty"`
		Duration int    `json:"duration,omitempty"`
	}

	go func() {
		for {
			mt, message, err := conn.ReadMessage()

			if err != nil || mt == websocket.CloseMessage {
				log.Printf("[error] %v", err)

				break // Выходим из цикла, если клиент пытается закрыть соединение или связь прервана
			}

			go s.handleMessage(message)
		}
	}()

	for {
		musicStruct := <-s.musicChan

		augMsg := ReqAug{
			IsReward: musicStruct.IsReward,
			Name:     musicStruct.Name,
			Title:    musicStruct.Title,
			Link:     musicStruct.Link,
			Duration: musicStruct.Duration,
		}

		var network bytes.Buffer
		enc := json.NewEncoder(&network)
		err := enc.Encode(augMsg)
		if err != nil {
			log.Println(err)

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
		log.Println("[error]", err)
	}
	log.Println("[info] request from server:", resp.Link)

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
func (s *Server) playlist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.myPlaylist)
}

// sss
func (s *Server) writeByteMsg(typeMsg string, message []byte) {
	conn := s.clients[typeMsg]
	conn.WriteMessage(websocket.TextMessage, message)
}
