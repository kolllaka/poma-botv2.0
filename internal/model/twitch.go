package model

import "time"

type MessageToChat struct {
	Platform   string       `json:"platform,omitempty"`
	Sender     Sender       `json:"sender,omitempty"`
	Time       time.Time    `json:"time,omitempty"`
	TextToChat []TextToChat `json:"text_to_chat,omitempty"`
}

type TextToChat struct {
	Text    string `json:"text"`
	IsEmote bool   `json:"is_emote"`
}

type Emotes struct {
	Data     []Emote `json:"data,omitempty"`
	Template string  `json:"template,omitempty"`
}

type Emote struct {
	Id        string   `json:"id,omitempty"`
	Name      string   `json:"name,omitempty"`
	Images    Images   `json:"images,omitempty"`
	Format    []string `json:"format,omitempty"`
	Scale     []string `json:"scale,omitempty"`
	ThemeMode []string `json:"theme_mode,omitempty"`
}

type Images struct {
	Url1x string `json:"url_1x,omitempty"`
	Url2x string `json:"url_2x,omitempty"`
	Url4x string `json:"url_4x,omitempty"`
}

type Message struct {
	TypeOfMsg string   `json:"type_of_msg"`
	Position  Position `json:"position,omitempty"`
	KeyBoard  string   `json:"key_board,omitempty"`
}

type Position struct {
	X float64 `json:"x,omitempty"`
	Y float64 `json:"y,omitempty"`
}

// Reward
type Reward struct {
	Type string `json:"type"`
	Data Data   `json:"data"`
}

type Data struct {
	Timestamp  string     `json:"timestamp"`
	Redemption Redemption `json:"redemption"`
}

type Redemption struct {
	ID         string      `json:"id"`
	User       User        `json:"user"`
	ChannelID  string      `json:"channel_id"`
	RedeemedAt string      `json:"redeemed_at"`
	Reward     RewardClass `json:"reward"`
	UserInput  string      `json:"user_input"`
	Status     string      `json:"status"`
	Cursor     string      `json:"cursor"`
}

type RewardClass struct {
	ID                                string              `json:"id"`
	ChannelID                         string              `json:"channel_id"`
	Title                             string              `json:"title"`
	Prompt                            string              `json:"prompt"`
	Cost                              int64               `json:"cost"`
	IsUserInputRequired               bool                `json:"is_user_input_required"`
	IsSubOnly                         bool                `json:"is_sub_only"`
	Image                             interface{}         `json:"image"`
	DefaultImage                      DefaultImage        `json:"default_image"`
	BackgroundColor                   string              `json:"background_color"`
	IsEnabled                         bool                `json:"is_enabled"`
	IsPaused                          bool                `json:"is_paused"`
	IsInStock                         bool                `json:"is_in_stock"`
	MaxPerStream                      MaxPerStream        `json:"max_per_stream"`
	ShouldRedemptionsSkipRequestQueue bool                `json:"should_redemptions_skip_request_queue"`
	TemplateID                        interface{}         `json:"template_id"`
	UpdatedForIndicatorAt             string              `json:"updated_for_indicator_at"`
	MaxPerUserPerStream               MaxPerUserPerStream `json:"max_per_user_per_stream"`
	GlobalCooldown                    GlobalCooldown      `json:"global_cooldown"`
	RedemptionsRedeemedCurrentStream  interface{}         `json:"redemptions_redeemed_current_stream"`
	CooldownExpiresAt                 interface{}         `json:"cooldown_expires_at"`
}

type DefaultImage struct {
	URL1X string `json:"url_1x"`
	URL2X string `json:"url_2x"`
	URL4X string `json:"url_4x"`
}

type GlobalCooldown struct {
	IsEnabled             bool  `json:"is_enabled"`
	GlobalCooldownSeconds int64 `json:"global_cooldown_seconds"`
}

type MaxPerStream struct {
	IsEnabled    bool  `json:"is_enabled"`
	MaxPerStream int64 `json:"max_per_stream"`
}

type MaxPerUserPerStream struct {
	IsEnabled           bool  `json:"is_enabled"`
	MaxPerUserPerStream int64 `json:"max_per_user_per_stream"`
}

type User struct {
	ID          string `json:"id"`
	Login       string `json:"login"`
	DisplayName string `json:"display_name"`
}

// Analitics from streams
// https://api.twitch.tv/helix/streams
type Streams struct {
	Data       []Datum    `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Datum struct {
	ID           string        `json:"id"`
	UserID       string        `json:"user_id"`
	UserLogin    string        `json:"user_login"`
	UserName     string        `json:"user_name"`
	GameID       string        `json:"game_id"`
	GameName     string        `json:"game_name"`
	Type         string        `json:"type"`
	Title        string        `json:"title"`
	Tags         []string      `json:"tags"`
	ViewerCount  int64         `json:"viewer_count"`
	StartedAt    string        `json:"started_at"`
	Language     string        `json:"language"`
	ThumbnailURL string        `json:"thumbnail_url"`
	TagIDS       []interface{} `json:"tag_ids"`
	IsMature     bool          `json:"is_mature"`
}

type Pagination struct {
	Cursor string `json:"cursor"`
}
