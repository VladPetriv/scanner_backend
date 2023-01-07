package model

type DBReply struct {
	ID        int    `db:"id"`
	MessageID int    `db:"message_id"`
	UserID    int    `db:"user_id"`
	Title     string `db:"title"`
	ImageURL  string `db:"image_url"`
}

type FullReply struct {
	ID           int    `json:"id" db:"id"`
	UserID       int    `json:"userId" db:"user_id"`
	Title        string `json:"title" db:"title"`
	ImageURL     string `json:"ImageUrl" db:"image_url"`
	FullName     string `json:"Fullname" db:"fullname"`
	UserImageURL string `json:"userImageUrl" db:"user_image_url"`
}
