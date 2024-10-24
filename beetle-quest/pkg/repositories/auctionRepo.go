package repositories

import (
	"beetle-quest/pkg/models"
)

type AuctionRepo interface {
	AddAuction(*models.Auction) error
	FindByID(models.AuctionId) (*models.Auction, bool)
	FindByUUID(models.ApiUUID) (*models.Auction, bool)
}
