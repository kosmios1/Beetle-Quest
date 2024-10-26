package models

type UserId []byte

type User struct {
	UserID UserId  `json:"user_id"`
	UUID   ApiUUID `json:"uuid"`
	// Salt         []byte        `json:"salt"` // NOTE: Bcrypt already includes salt
	Username     string        `json:"username"`
	Email        string        `json:"email"`
	PasswordHash []byte        `json:"password_hash"`
	Gachas       []Gacha       `json:"gachas"`
	Transactions []Transaction `json:"transactions"`
}
