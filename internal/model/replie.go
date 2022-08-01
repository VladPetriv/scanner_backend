package model

type DBReplie struct {
	ID        int
	MessageID int
	UserID    int
	Title     string
	ImageURL  string
}

type FullReplie struct {
	ID           int    `json:"id" db:"id"`
	UserID       int    `json:"userId" db:"user_id"`
	Title        string `json:"title" db:"title"`
	ImageURL     string `json:"ImageUrl" db:"imageurl"`
	FullName     string `json:"Fullname" db:"fullname"`
	UserImageURL string `json:"userImageUrl" db:"userimageurl"`
}
