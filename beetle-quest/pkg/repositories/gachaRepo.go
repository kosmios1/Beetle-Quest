package repositories

import "beetle-quest/pkg/models"

type GachaRepo interface {
	FindByID(models.UUID) (*models.Gacha, bool)
}
