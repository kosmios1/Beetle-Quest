package repositories

import (
	"beetle-quest/pkg/models"
)

type MarketRepo interface {
	Create(*models.Auction) bool
	Update(*models.Auction) bool
	Delete(*models.Auction) bool

	GetAll() ([]models.Auction, bool)
	GetUserAuctions(models.UUID) ([]models.Auction, bool)
	FindByID(models.UUID) (*models.Auction, bool)

	GetBidListOfAuction(models.UUID) ([]models.Bid, bool)
	BidToAuction(*models.Bid) bool

	GetAllTransactions() ([]models.Transaction, bool)
	DeleteUserTransactionHistory(models.UUID) bool
	GetUserTransactionHistory(models.UUID) ([]models.Transaction, bool)

	AddTransaction(*models.Transaction) bool
}
