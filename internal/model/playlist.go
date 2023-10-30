package model

type Playlist struct {
	IsYouTube bool   `json:"isyoutube"`
	IsReward  bool   `json:"isreward"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	Link      string `json:"link"`
	Duration  int    `json:"duration"`
}
