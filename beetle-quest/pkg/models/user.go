package models

type User struct {
	UserID       UUID   `json:"user_id"       gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Username     string `json:"username"      gorm:"uniqueIndex;type:varchar(255);not null"`
	Email        string `json:"email"         gorm:"uniqueIndex;type:varchar(255);not null"`
	PasswordHash []byte `json:"password_hash" gorm:"not null"`
	// Gachas       []Gacha       `json:"gachas"`
	// Transactions []Transaction `json:"transactions"`
}
