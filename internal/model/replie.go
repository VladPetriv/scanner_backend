package model

type Replie struct {
	ID        int    `json:"id"`
	MessageID int    `json:"messageId"`
	UserID    int    `json:"userId"`
	Title     string `json:"title"`
}

type FullReplie struct {
	ID           int    `json:"id"`
	UserID       int    `json:"userId"`
	Title        string `json:"title"`
	FullName     string `json:"Fullname"`
	UserImageURL string `json:"userImageUrl"`
}
