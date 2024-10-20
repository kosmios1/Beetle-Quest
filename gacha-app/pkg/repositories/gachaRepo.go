package repositories

import "gacha-app/pkg/models"

type GachaRepo interface {
	ValidateGachaID(models.GachaId) bool
}
