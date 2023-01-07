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
	ID         int    `db:"id"`
	ChannelID  int    `db:"channel_id"`
	UserID     int    `db:"user_id"`
	Title      string `db:"title"`
	MessageURL string `db:"message_url"`
	ImageURL   string `db:"imageurl"`
}

type FullMessage struct {
	ID              int    `db:"id"`
	MessageURL      string `db:"message_url"`
	Title           string `db:"title"`
	ImageURL        string `db:"imageurl"`
	ChannelID       int    `db:"channelid"`
	ChannelName     string `db:"channelname"`
	ChannelTitle    string `db:"channeltitle"`
	ChannelImageURL string `db:"channelimageurl"`
	UserID          int    `db:"userid"`
	FullName        string `db:"fullname"`
	UserImageURL    string `db:"userimageurl"`
	ReplieCount     int    `db:"count"`
	Replies         []FullReply
	SavedID         int
	Status          bool
}
