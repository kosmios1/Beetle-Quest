package repositories

import (
	"beetle-quest/pkg/models"
)

type AuctionRepo interface {
	AddAuction(*models.Auction) error
	FindByID(models.UUID) (*models.Auction, bool)
}
