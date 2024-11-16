package models

type Admin struct {
	AdminId      UUID   `json:"admin_id"`
	PasswordHash []byte `json:"password_hash"`
	OtpSecret    string `json:"otp_secret"`
	Email        string `json:"email"`
}
