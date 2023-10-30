package model

type MsgFromSocket struct {
	Reason string   `json:"reason,omitempty"`
	Song   Playlist `json:"song,omitempty"`
}
