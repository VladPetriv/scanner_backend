package model

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullname"`
	PhotoURL string `json:"photoUrl"`
}

type WebUser struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
