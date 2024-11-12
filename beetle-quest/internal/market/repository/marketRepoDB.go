package repository

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
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
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("Failed to connect to the Database: %v", err)
			time.Sleep(1 * time.Second)
		} else {
			repo.db = db
			break
		}
	}

	repo.db.AutoMigrate(&models.Auction{})
	return repo
}

func (r *MarketRepo) Create(auction *models.Auction) bool {
	result := r.db.Table("auctions").Create(auction)
	if result.Error != nil {
		return false
	}
	return true
}

func (r *MarketRepo) FindByID(id models.UUID) (*models.Auction, bool) {
	var auction models.Auction
	result := r.db.Table("auctions").First(&auction, id)
	if result.Error != nil {
		return nil, false
	}
	return &auction, true
}

func (r *MarketRepo) GetAll() ([]models.Auction, bool) {
	var auctions []models.Auction
	result := r.db.Table("auctions").Find(&auctions)
	if result.Error != nil {
		return nil, false
	}
	return auctions, true
}
