package model

type Channel struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Title    string `json:"title"`
	PhotoURL string `json:"photoUrl"`
	Stats    Stat
}

type Stat struct {
	MessagesCount int
	RepliesCount  int
}
