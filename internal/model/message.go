package model

type Message struct {
	ID        int    `json:"id"`
	ChannelID int    `json:"channelId`
	UserID    int    `json:"userId"`
	Title     string `json:"title"`
}

type FullMessage struct {
	ID          int    `json:"ID`
	Title       string `json:"title"`
	ChannelName string `json:"chanellName"`
	FullName    string `json:"fullName"`
	PhotoURL    string `json:"PhotoURL"`
	ReplieCount int    `json:"replieCount"`
}
