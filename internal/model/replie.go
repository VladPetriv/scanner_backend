package model

type DBReplie struct {
	ID        int
	MessageID int
	UserID    int
	Title     string
	ImageURL  string
}

type FullReplie struct {
	ID           int    `json:"id"`
	UserID       int    `json:"userId"`
	Title        string `json:"title"`
	ImageURL     string `json:"ImageUrl"`
	FullName     string `json:"Fullname"`
	UserImageURL string `json:"userImageUrl"`
}
