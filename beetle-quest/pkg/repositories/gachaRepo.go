package repositories

import "beetle-quest/pkg/models"

type GachaRepo interface {
	Create(*models.Gacha) error
	Update(*models.Gacha) error
	Delete(*models.Gacha) error

	GetAll() ([]models.Gacha, error)
	FindByID(models.UUID) (*models.Gacha, error)

	AddGachaToUser(models.UUID, models.UUID) error
	RemoveGachaFromUser(models.UUID, models.UUID) error

	RemoveUserGachas(models.UUID) error
	GetUserGachas(models.UUID) ([]models.Gacha, error)
}
