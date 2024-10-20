package mock_repositories

import (
	"encoding/hex"
	"gacha-app/pkg/models"
)

type MockGachaRepo struct{}

var Gachas = map[string]*models.Gacha{
	"6761636861": { // "gacha" in hex
		GachaID: []byte("6761636861"),
		Name:    "Rare Gacha",
		Rarity:  models.Rare,
		Price:   500,
	},
}

func (m MockGachaRepo) ValidateGachaID(id models.GachaId) bool {
	hexId := hex.EncodeToString(id)
	if _, ok := Gachas[hexId]; !ok {
		return false
	}
	return true
}
