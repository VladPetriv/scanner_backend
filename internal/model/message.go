package model

type Message struct {
	ID        int    `json:"ID"`
	ChannelID int    `json:"ChannelID"`
	Title     string `json:"Title"`
}
