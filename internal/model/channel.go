package model

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
