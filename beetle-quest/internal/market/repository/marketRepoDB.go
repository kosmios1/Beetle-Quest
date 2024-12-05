package repository

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbHost         string = utils.FindEnv("POSTGRES_HOST")
	dbUserName     string = utils.FindEnv("POSTGRES_USER")
	dbUserPassword string = utils.FindEnv("POSTGRES_PASSWORD")
	dbName         string = utils.FindEnv("POSTGRES_DB")
	dbPort         string = utils.FindEnv("POSTGRES_PORT")
	dbSSLMode      string = utils.FindEnv("POSTGRES_SSLMODE")
	dbTimeZone     string = utils.FindEnv("POSTGRES_TIMEZONE")
)

type MarketRepo struct {
	db *gorm.DB
}

func NewMarketRepo() *MarketRepo {
	var repo = &MarketRepo{}
	for {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s timezone=%s", dbHost, dbUserName, dbUserPassword, dbName, dbPort, dbSSLMode, dbTimeZone)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
		if err != nil {
			log.Printf("Failed to connect to the Database: %v", err)
			time.Sleep(1 * time.Second)
		} else {
			repo.db = db
			break
		}
	}
	return repo
}

func (r *MarketRepo) Create(auction *models.Auction) error {
	result := r.db.Table("auctions").Create(auction)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return models.ErrAuctionAltreadyExists
		}
		return models.ErrInternalServerError
	}

	if result.RowsAffected == 0 {
		return models.ErrInternalServerError
	}
	return nil
}

func (r *MarketRepo) Update(auction *models.Auction) error {
	result := r.db.Table("auctions").Where("auction_id = ?", auction.AuctionID).Select("*").Updates(auction)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.ErrAuctionNotFound
		} else if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return models.ErrGachaAlreadyExists
		}
		return models.ErrInternalServerError
	}

	if result.RowsAffected == 0 {
		return models.ErrInternalServerError
	}
	return nil
}

func (r *MarketRepo) Delete(auction *models.Auction) error {
	result := r.db.Table("auctions").Delete(auction, models.Auction{AuctionID: auction.AuctionID})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.ErrAuctionNotFound
		}
		return models.ErrInternalServerError
	}

	// if result.RowsAffected == 0 {
	// 	return models.ErrInternalServerError
	// }
	return nil
}

func (r *MarketRepo) FindByID(id models.UUID) (*models.Auction, error) {
	var auction models.Auction
	result := r.db.Table("auctions").First(&auction, models.Auction{AuctionID: id})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, models.ErrAuctionNotFound
		}
		return nil, models.ErrInternalServerError
	}
	return &auction, nil
}

func (r *MarketRepo) DeleteUserTransactionHistory(uid models.UUID) error {
	result := r.db.Table("transactions").Delete(models.Transaction{}, models.Transaction{UserID: uid})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) || errors.Is(result.Error, gorm.ErrEmptySlice) {
			return models.ErrUserTransactionNotFound
		}
		return models.ErrInternalServerError
	}

	// if result.RowsAffected == 0 {
	// 	return models.ErrInternalServerError
	// }
	return nil
}

func (r *MarketRepo) GetUserTransactionHistory(uid models.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	result := r.db.Table("transactions").Where("user_id = ?", uid).Find(&transactions)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) || errors.Is(result.Error, gorm.ErrEmptySlice) {
			return nil, models.ErrTransactionNotFound
		}
		return nil, models.ErrInternalServerError
	}
	return transactions, nil
}

func (r *MarketRepo) GetUserAuctions(uid models.UUID) ([]models.Auction, error) {
	var auctions []models.Auction
	result := r.db.Table("auctions").Where("owner_id = ?", uid).Find(&auctions)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) || errors.Is(result.Error, gorm.ErrEmptySlice) {
			return nil, models.ErrAuctionNotFound
		}
		return nil, models.ErrInternalServerError
	}
	return auctions, nil
}

func (r *MarketRepo) GetAll() ([]models.Auction, error) {
	var auctions []models.Auction
	result := r.db.Table("auctions").Find(&auctions)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) || errors.Is(result.Error, gorm.ErrEmptySlice) {
			return nil, models.ErrAuctionNotFound
		}
		return nil, models.ErrInternalServerError
	}
	return auctions, nil
}

func (r *MarketRepo) BidToAuction(bid *models.Bid) error {
	result := r.db.Table("bids").Create(bid)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return models.ErrCouldNotBidToAuction
		}
		return models.ErrInternalServerError
	}

	if result.RowsAffected == 0 {
		return models.ErrInternalServerError
	}
	return nil
}

func (r *MarketRepo) GetBidListOfAuction(aid models.UUID) ([]models.Bid, error) {
	var bids []models.Bid
	result := r.db.Table("bids").Where("auction_id = ?", aid).Find(&bids)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) || errors.Is(result.Error, gorm.ErrEmptySlice) {
			return nil, models.ErrBidsNotFound
		}
		return nil, models.ErrInternalServerError
	}
	return bids, nil
}

func (r *MarketRepo) AddTransaction(transaction *models.Transaction) error {
	result := r.db.Table("transactions").Create(transaction)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return models.ErrCouldNotAddTransaction
		}
		return models.ErrInternalServerError
	}

	if result.RowsAffected == 0 {
		return models.ErrInternalServerError
	}
	return nil
}

func (r *MarketRepo) GetAllTransactions() ([]models.Transaction, error) {
	var transactions []models.Transaction
	result := r.db.Table("transactions").Find(&transactions)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) || errors.Is(result.Error, gorm.ErrEmptySlice) {
			return nil, models.ErrTransactionNotFound
		}
		return nil, models.ErrInternalServerError
	}
	return transactions, nil
}
