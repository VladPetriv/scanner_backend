package model

type Saved struct {
	ID        int `json:"id" db:"id"`
	WebUserID int `json:"webUserId" db:"user_id"`
	MessageID int `json:"messageId" db:"message_id"`
}
