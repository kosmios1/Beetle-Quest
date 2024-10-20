package repositories

import (
	"gacha-app/pkg/models"
)

type AuctionRepo interface {
	AddAuction(*models.Auction) error
}
