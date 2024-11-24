package models

type Admin struct {
	AdminId      UUID   `json:"admin_id"        gorm:"unique;primaryKey"`
	PasswordHash []byte `json:"password_hash"   gorm:"not null"`
	OtpSecret    string `json:"otp_secret"      gorm:"not null"`
	Email        string `json:"email"           gorm:"unique"`
}
