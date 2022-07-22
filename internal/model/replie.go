package model

type FullReplie struct {
	ID           int    `json:"id"`
	UserID       int    `json:"userId"`
	Title        string `json:"title"`
	ImageURL     string `json:"ImageUrl"`
	FullName     string `json:"Fullname"`
	UserImageURL string `json:"userImageUrl"`
}
