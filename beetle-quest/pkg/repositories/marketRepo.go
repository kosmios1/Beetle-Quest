package repositories

import (
	"beetle-quest/pkg/models"
)

type MarketRepo interface {
	Create(*models.Auction) error
	Update(*models.Auction) error
	Delete(*models.Auction) error

	GetAll() ([]models.Auction, error)
	GetUserAuctions(models.UUID) ([]models.Auction, error)
	FindByID(models.UUID) (*models.Auction, error)

	GetBidListOfAuction(models.UUID) ([]models.Bid, error)
	BidToAuction(*models.Bid) error

	GetAllTransactions() ([]models.Transaction, error)
	DeleteUserTransactionHistory(models.UUID) error
	GetUserTransactionHistory(models.UUID) ([]models.Transaction, error)

	AddTransaction(*models.Transaction) error
}
