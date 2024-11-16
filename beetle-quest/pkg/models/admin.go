package models

type Admin struct {
	AdminId      UUID   `json:"admin_id"`
	PasswordHash []byte `json:"password_hash"`
	OtpSecret    string `json:"otp_secret"`
	Email        string `json:"email"`
}

type AdminSafeData struct {
	AdminId      UUID   `json:"admin_id"`
	PasswordHash []byte `json:"password"`
	Otp          string `json:"otp"`
	Email        string `json:"email"`
}
