package model

type DBChannel struct {
	ID       int    `json:"ID"`
	Name     string `json:"Username"`
	Title    string `json:"Title"`
	ImageURL string `json:"ImageURL"`
}

type Channel struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Title    string `json:"title"`
	ImageURL string `json:"imageUrl"`
	Stats    Stat
}

type Stat struct {
	MessagesCount int
	RepliesCount  int
}
