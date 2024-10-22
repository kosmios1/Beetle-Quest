package mock_repositories

import (
	"encoding/hex"
	"gacha-app/pkg/models"
)

type MockUserRepo struct{}

var Users = map[string]*models.User{
	"4b6974": { // "kit" in hex
		UserID:       []byte("4b6974"),
		Salt:         []byte("72616e"),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: []byte("68617368"),
		Gachas:       []models.Gacha{},
		Transactions: []models.Transaction{},
	},
}

func (m MockUserRepo) ValidateUserID(id models.UserId) bool {
	hexId := hex.EncodeToString(id)
	if _, ok := Users[hexId]; !ok {
		return false
	}
	return true
}
