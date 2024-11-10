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

type GachaRepo struct {
	DB *gorm.DB
}

func NewGachaRepo() *GachaRepo {
	if dbHost == "" || dbUserName == "" || dbUserPassword == "" || dbName == "" || dbPort == "" || dbSSLMode == "" || dbTimeZone == "" {
		log.Fatalf("Either POSTGRES_HOST, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, POSTGRES_PORT, POSTGRES_SSLMODE or POSTGRES_TIMEZONE is not set")
	}

	var repo = &GachaRepo{}
	for {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s timezone=%s", dbHost, dbUserName, dbUserPassword, dbName, dbPort, dbSSLMode, dbTimeZone)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("Failed to connect to the Database: %v", err)
			time.Sleep(1 * time.Second)
		} else {
			repo.DB = db
			break
		}
	}

	// This will create the table if it does not exist and will keep the schema updated
	err := repo.DB.AutoMigrate(&models.Gacha{})
	if err != nil {
		log.Printf("Failed to migrate the database: %v", err)
	}

	return repo
}

func (r GachaRepo) FindByID(id models.UUID) (*models.Gacha, bool) {
	var gacha models.Gacha
	result := r.DB.First(&gacha, models.Gacha{GachaID: id})
	if result.Error != nil {
		return nil, false
	}
	return &gacha, true
}

func (r GachaRepo) GetAll() ([]models.Gacha, bool) {
	var gachas []models.Gacha
	result := r.DB.Find(&gachas)
	if result.Error != nil {
		return nil, false
	}
	return gachas, true
}
