package repositories

import "gacha-app/pkg/models"

type GachaRepo interface {
	ValidateGachaID(gachaID models.GachaId) bool
}
