package repositories

import (
	"beetle-quest/pkg/models"
)

type MarketRepo interface {
	GetAll() ([]models.Auction, bool)
	Create(*models.Auction) bool
	FindByID(models.UUID) (*models.Auction, bool)
}
