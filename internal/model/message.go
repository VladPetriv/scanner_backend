package model

type TgMessage struct {
	Message    string `json:"Message"`
	MessageURL string `json:"MessageURL"`
	ImageURL   string `json:"ImageURL"`
	FromID     struct {
		Username string `json:"Username"`
		ImageURL string `json:"ImageURL"`
		Fullname string `json:"Fullname"`
	} `json:"FromID"`
	PeerID struct {
		Username string `json:"Username"`
	} `json:"PeerID"`
	Replies struct {
		Count    int `json:"Count"`
		Messages []struct {
			FromID struct {
				Username string `json:"Username"`
				Fullname string `json:"Fullname"`
				ImageURL string `json:"ImageURL"`
			} `json:"FromID"`
			Message  string `json:"Message"`
			ImageURL string `json:"ImageURL"`
		} `json:"Messages"`
	} `json:"Replies"`
}

type DBMessage struct {
	ID         int
	ChannelID  int
	UserID     int
	Title      string
	MessageURL string
	ImageURL   string
}

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
