package mock_repositories

import (
	"encoding/hex"
	"gacha-app/pkg/models"
)

type MockGachaRepo struct {
	Gachas map[string]*models.Gacha
}

func NewMockGachaRepo() *MockGachaRepo {
	return &MockGachaRepo{
		Gachas: make(map[string]*models.Gacha),
	}
}

func (m MockGachaRepo) ValidateGachaID(id *models.GachaId) bool {
	hexId := hex.EncodeToString(*id)
	if _, ok := m.Gachas[hexId]; !ok {
		return false
	}
	return true
}
