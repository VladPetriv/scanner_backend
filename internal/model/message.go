package model

type Message struct {
	ID        int    `json:"id"`
	ChannelID int    `json:"channelId"`
	UserID    int    `json:"userId"`
	Title     string `json:"title"`
}

type FullMessage struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	ChannelID       int    `json:"channelId"`
	ChannelName     string `json:"chanellName"`
	ChannelTitle    string `json:"channelTitle"`
	ChannelPhotoURL string `json:"channelPhotoUrl"`
	UserID          int    `json:"userId"`
	FullName        string `json:"fullName"`
	PhotoURL        string `json:"photoURL"`
	ReplieCount     int    `json:"replieCount"`
	Replies         []FullReplie
}
