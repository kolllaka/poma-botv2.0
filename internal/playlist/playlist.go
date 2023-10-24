package playlist

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/KoLLlaka/poma-botv2.0/internal/model"
)

var (
	trueFiles = map[string]bool{
		".mp4":  true,
		".mp3":  true,
		".webm": true,
	}
)

func LoadMyPlaylist(path string) []model.Playlist {
	playlist := []model.Playlist{}
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	i := 0
	for _, file := range files {
		i++
		ext := filepath.Ext(file.Name())
		if _, ok := trueFiles[ext]; !ok {
			continue
		}

		fmt.Printf("%s ", file.Name())
		playlist = append(playlist, model.Playlist{
			IsYouTube: false,
			IsReward:  false,
			Title:     file.Name(),
			Link:      fmt.Sprintf("./audio/%s", file.Name()),
		})
	}
	fmt.Println("общее количество файлов:", i)

	return playlist
}

// из текста выделить ссылку на youtube
func SongRequest(text string) (string, error) {
	urlReg := regexp.MustCompile(`https://www.youtube.com/watch\?v=([a-zA-Z0-9_-]*)|https://youtu.be/([a-zA-Z0-9_-]*)`)
	resUrl := urlReg.FindStringSubmatch(text)

	if len(resUrl) > 0 {
		if resUrl[1] == "" {
			return resUrl[2], nil
		}

		return resUrl[1], nil
	}

	return "", fmt.Errorf("unknown song: %s", text)
}

// проверка на аудио
func ListOfSongsFromPlaylist(playlistId string, key string, next string) []*model.Playlist {
	type ListResp struct {
		NextPageToken string `json:"nextPageToken,omitempty"`
		PrevPageToken string `json:"prevPageToken,omitempty"`
		Items         []struct {
			ContentDetails struct {
				VideoId string `json:"videoId,omitempty"`
			} `json:"contentDetails,omitempty"`
		} `json:"items,omitempty"`
	}

	songList := []*model.Playlist{}
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=contentDetails&playlistId=%s&key=%s&maxResults=50", playlistId, key)
	if next != "" {
		url = fmt.Sprintf("%s&pageToken=%s", url, next)
	}
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)

		return nil
	}
	defer resp.Body.Close()
	log.Println(resp.StatusCode)

	listResp := ListResp{}
	if resp.StatusCode != http.StatusOK {
		log.Println("Всё плохо!")

		return nil
	}
	json.NewDecoder(resp.Body).Decode(&listResp)

	if listResp.NextPageToken != "" {
		nextSongList := ListOfSongsFromPlaylist(playlistId, key, listResp.NextPageToken)

		songList = append(songList, nextSongList...)
	}

	for _, data := range listResp.Items {
		song := model.Playlist{
			Link: data.ContentDetails.VideoId,
		}

		songList = append(songList, &song)
	}

	return songList
}

func ReqSongInfo(song string, key string) []*model.Playlist {
	type SongResp struct {
		Items []struct {
			Id      string `json:"id"`
			Snippet struct {
				Title string `json:"title,omitempty"`
			} `json:"snippet,omitempty"`
			ContentDetails struct {
				Duration string `json:"duration,omitempty"`
			} `json:"contentDetails,omitempty"`
		} `json:"items,omitempty"`
	}
	songsList := []*model.Playlist{}

	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?id=%s&key=%s&part=snippet,contentDetails,statistics", song, key)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)

		return nil
	}
	defer resp.Body.Close()
	log.Println(resp.StatusCode)
	if resp.StatusCode != 200 {
		var response interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		log.Println(response)
	}

	songResp := SongResp{}
	json.NewDecoder(resp.Body).Decode(&songResp)

	for _, song := range songResp.Items {
		duration := song.ContentDetails.Duration
		songList := model.Playlist{
			IsYouTube: true,
			Title:     song.Snippet.Title,
			Link:      song.Id,
			Duration:  formatTimeFromYoutube(duration),
		}

		songsList = append(songsList, &songList)
	}

	return songsList
}

// перевод ютубовского формата в формат в секндах
func formatTimeFromYoutube(time string) int {
	var fotmatTime int = 0
	reg := regexp.MustCompile(`PT(\d*S)|PT(\d*M)(\d*S)|PT(\d*H)(\d*M)(\d*S)`)
	res := reg.FindStringSubmatch(time)

	if res == nil {
		return 0
	}

	for _, val := range res[1:] {
		if val != "" {
			sim := string(val[len(val)-1])
			chislo, _ := strconv.Atoi(val[:len(val)-1])

			switch sim {
			case "H":
				fotmatTime += chislo * 3600
			case "M":
				fotmatTime += chislo * 60
			case "S":
				fotmatTime += chislo
			}
		}
	}

	return fotmatTime
}
