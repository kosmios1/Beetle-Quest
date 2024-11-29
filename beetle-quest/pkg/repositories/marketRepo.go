package repositories

import (
	"beetle-quest/pkg/models"
)

type MarketRepo interface {
	Create(*models.Auction) error
	Update(*models.Auction) error
	Delete(*models.Auction) error

	GetAll() ([]models.Auction, error)
	GetUserAuctions(uid models.UUID) ([]models.Auction, error)
	FindByID(aid models.UUID) (*models.Auction, error)

	GetBidListOfAuction(aid models.UUID) ([]models.Bid, error)
	BidToAuction(*models.Bid) error

	GetAllTransactions() ([]models.Transaction, error)
	DeleteUserTransactionHistory(uid models.UUID) error
	GetUserTransactionHistory(uid models.UUID) ([]models.Transaction, error)

	AddTransaction(*models.Transaction) error
}
