package repository

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"fmt"
	"log"
	"strings"
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
	db *gorm.DB
}

func NewGachaRepo() *GachaRepo {
	var repo = &GachaRepo{}
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

	// This will create the table if it does not exist and will keep the schema updated
	err := repo.db.AutoMigrate(&models.Gacha{})
	if err != nil {
		log.Printf("Failed to migrate the database: %v", err)
	}

	return repo
}

func (r GachaRepo) FindByID(id models.UUID) (*models.Gacha, bool) {
	var gacha models.Gacha
	result := r.db.Table("gachas").First(&gacha, models.Gacha{GachaID: id})
	if result.Error != nil {
		return nil, false
	}
	return &gacha, true
}

func (r GachaRepo) GetAll() ([]models.Gacha, bool) {
	var gachas []models.Gacha
	result := r.db.Table("gachas").Find(&gachas)
	if result.Error != nil {
		return nil, false
	}
	return gachas, true
}

func (r GachaRepo) AddGachaToUser(uid models.UUID, gid models.UUID) bool {
	result := r.db.Table("user_gacha").Create(struct {
		UserID  models.UUID
		GachaID models.UUID
	}{uid, gid})
	if result.Error != nil {
		// TODO: Maybe return error?
		if strings.Contains(result.Error.Error(), "duplicate key") {
			return false
		}
		return false
	}
	return true
}
