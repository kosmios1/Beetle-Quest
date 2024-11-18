package repositories

import "beetle-quest/pkg/models"

type GachaRepo interface {
	Create(*models.Gacha) bool
	Update(*models.Gacha) bool
	Delete(*models.Gacha) bool

	GetAll() ([]models.Gacha, bool)
	FindByID(models.UUID) (*models.Gacha, bool)

	AddGachaToUser(models.UUID, models.UUID) bool
	RemoveGachaFromUser(models.UUID, models.UUID) bool

	RemoveUserGachas(models.UUID) bool
	GetUserGachas(models.UUID) ([]models.Gacha, bool)
}
