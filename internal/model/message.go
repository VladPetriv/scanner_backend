package model

type FullMessage struct {
	ID              int    `json:"id"`
	MessageURL      string `json:"messageURL"`
	Title           string `json:"title"`
	ImageURL        string `json:"imageUrl"`
	ChannelID       int    `json:"channelId"`
	ChannelName     string `json:"chanellName"`
	ChannelTitle    string `json:"channelTitle"`
	ChannelImageURL string `json:"channelImageUrl"`
	UserID          int    `json:"userId"`
	FullName        string `json:"fullName"`
	UserImageURL    string `json:"userImageUrl"`
	ReplieCount     int    `json:"replieCount"`
	Replies         []FullReplie
	SavedID         int
	Status          bool
}
