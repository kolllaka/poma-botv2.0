package model

import "time"

type MessageFromIRC struct {
	TypeOfMessage string
	Platform      string    `json:"platform,omitempty"`
	Channel       string    `json:"channel,omitempty"`
	ChannelId     string    `json:"channel_id,omitempty"`
	Sender        Sender    `json:"sender,omitempty"`
	Time          time.Time `json:"time,omitempty"`
	Text          string    `json:"text,omitempty"`
	EmoteMap      string
}

type Sender struct {
	From          string            `json:"from,omitempty"`
	FromId        string            `json:"from_id,omitempty"`
	Badges        map[string]string `json:"badges,omitempty"`
	BadgesInfo    map[string]string `json:"badges_info,omitempty"`
	Color         string            `json:"color,omitempty"`
	IsBroadcaster bool              `json:"is_broadcaster,omitempty"`
	IsModerator   bool              `json:"is_moderator,omitempty"`
	IsVIP         bool              `json:"is_vip,omitempty"`
	IsSubscriber  bool              `json:"is_subscriber,omitempty"`
}
