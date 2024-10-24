package repositories

import "beetle-quest/pkg/models"

type GachaRepo interface {
	FindByID(models.GachaId) (*models.Gacha, bool)
	FindByUUID(models.ApiUUID) (*models.Gacha, bool)
}
