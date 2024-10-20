package mock_repositories

import (
	"encoding/hex"
	"gacha-app/pkg/models"
)

type MockUserRepo struct {
	Users map[string]*models.User
}

func NewMockUserRepo() *MockUserRepo {
	return &MockUserRepo{
		Users: map[string]*models.User{},
	}
}

func (m MockUserRepo) ValidateUserID(id *models.UserId) bool {
	hexId := hex.EncodeToString(*id)
	if _, ok := m.Users[hexId]; !ok {
		return false
	}
	return true
}
