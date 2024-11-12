package repositories

import "beetle-quest/pkg/models"

type GachaRepo interface {
	GetAll() ([]models.Gacha, bool)
	FindByID(models.UUID) (*models.Gacha, bool)

	AddGachaToUser(models.UUID, models.UUID) bool

	GetUserGachas(models.UUID) ([]models.Gacha, bool)
}
