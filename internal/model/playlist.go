package model

type Playlist struct {
	IsYouTube bool   `json:"isyoutube,omitempty"`
	IsReward  bool   `json:"isreward,omitempty"`
	Name      string `json:"name,omitempty"`
	Title     string `json:"title,omitempty"`
	Link      string `json:"link,omitempty"`
	Duration  int    `json:"duration,omitempty"`
}
