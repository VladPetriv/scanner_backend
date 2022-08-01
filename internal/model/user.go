package model

type User struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	FullName string `json:"fullname" db:"fullname"`
	ImageURL string `json:"imageUrl" db:"imageurl"`
}

type WebUser struct {
	ID       int    `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}
