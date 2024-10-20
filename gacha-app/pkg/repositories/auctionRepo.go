package repositories

import (
	"gacha-app/pkg/models"
)

type AuctionRepo interface {
	AddAuction(*models.Auction) error
	GetAuction(models.AuctionId) (*models.Auction, error)

	VaildateAuctionID(models.AuctionId) bool
}
