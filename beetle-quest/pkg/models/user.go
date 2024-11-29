package models

type User struct {
	UserID       UUID   `json:"user_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Currency     int64  `json:"currency"`
	PasswordHash []byte `json:"password_hash"`
}
