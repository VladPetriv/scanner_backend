package model

type DBChannel struct {
	ID       int    `json:"ID" db:"id"`
	Name     string `json:"Username" db:"name"`
	Title    string `json:"Title" db:"title"`
	ImageURL string `json:"ImageURL" db:"image_url"`
}

type Channel struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Title    string `json:"title" db:"title"`
	ImageURL string `json:"imageUrl" db:"image_url"`
	Stats    Stat
}

type Stat struct {
	MessagesCount int
	RepliesCount  int
}
