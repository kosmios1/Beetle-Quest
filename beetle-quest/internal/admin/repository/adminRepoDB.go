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

type AdminRepo struct {
	db *gorm.DB
}

var (
	dbHost         string = utils.FindEnv("POSTGRES_HOST")
	dbUserName     string = utils.FindEnv("POSTGRES_USER")
	dbUserPassword string = utils.FindEnv("POSTGRES_PASSWORD")
	dbName         string = utils.FindEnv("POSTGRES_DB")
	dbPort         string = utils.FindEnv("POSTGRES_PORT")
	dbSSLMode      string = utils.FindEnv("POSTGRES_SSLMODE")
	dbTimeZone     string = utils.FindEnv("POSTGRES_TIMEZONE")
)

func NewAdminRepo() *AdminRepo {
	var repo = &AdminRepo{}
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
	// repo.db.AutoMigrate(&models.User{})
	return repo
}

func (r *AdminRepo) FindByID(id models.UUID) (*models.Admin, bool) {
	var admin models.Admin
	result := r.db.Table("admins").First(&admin, models.Admin{AdminId: id})
	if result.Error != nil {
		return nil, false
	}
	return &admin, true
}
