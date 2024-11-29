package repositories

import "beetle-quest/pkg/models"

type GachaRepo interface {
	Create(*models.Gacha) error
	Update(*models.Gacha) error
	Delete(*models.Gacha) error

	GetAll() ([]models.Gacha, error)
	FindByID(models.UUID) (*models.Gacha, error)

	AddGachaToUser(uid models.UUID, gid models.UUID) error
	RemoveGachaFromUser(uid models.UUID, gid models.UUID) error

	RemoveUserGachas(uid models.UUID) error
	GetUserGachas(uid models.UUID) ([]models.Gacha, error)
}
