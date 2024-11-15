package repositories

import (
	"beetle-quest/pkg/models"
)

type MarketRepo interface {
	Create(*models.Auction) bool
	Delete(*models.Auction) bool

	GetAll() ([]models.Auction, bool)
	GetUserAuctions(models.UUID) ([]models.Auction, bool)
	GetBidListOfAuction(models.UUID) ([]models.Bid, bool)

	FindByID(models.UUID) (*models.Auction, bool)
}
