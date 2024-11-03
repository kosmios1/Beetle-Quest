package repository

import (
	"beetle-quest/pkg/models"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbHost         = os.Getenv("POSTGRES_HOST")
	dbUserName     = os.Getenv("POSTGRES_USER")
	dbUserPassword = os.Getenv("POSTGRES_PASSWORD")
	dbName         = os.Getenv("POSTGRES_DB")
	dbPort         = os.Getenv("POSTGRES_PORT")
	dbSSLMode      = os.Getenv("POSTGRES_SSLMODE")
	dbTimeZone     = os.Getenv("POSTGRES_TIMEZONE")
)

// User Repository ========================================================================================

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo() *UserRepo {
	if dbHost == "" || dbUserName == "" || dbUserPassword == "" || dbName == "" || dbPort == "" || dbSSLMode == "" || dbTimeZone == "" {
		log.Fatalf("Either POSTGRES_HOST, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, POSTGRES_PORT, POSTGRES_SSLMODE or POSTGRES_TIMEZONE is not set")
	}

	var repo = &UserRepo{}
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
	err := repo.DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Printf("Failed to migrate the database: %v", err)
	}

	return repo
}

func (r UserRepo) Create(email, username string, hashedPassword []byte, currency int64) bool {
	result := r.DB.Create(&models.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		Currency:     currency,
	})

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			// User with that email/Username already exists
			return false
		}
		// result.Error.Error()
		return false
	}
	return true
}

func (r UserRepo) Update(user *models.User) bool {
	result := r.DB.Save(user)
	if result.Error != nil {
		return false
	}
	return true
}

func (r UserRepo) Delete(user *models.User) bool {
	result := r.DB.Delete(user)
	if result.Error != nil {
		return false
	}
	return true
}

func (r UserRepo) FindByUsername(username string) (*models.User, bool) {
	var user models.User
	result := r.DB.First(&user, models.User{Username: username})
	if result.Error != nil {
		// No user with that username exists
		return nil, false
	}
	return &user, true
}

func (r UserRepo) FindByID(id models.UUID) (*models.User, bool) {
	var user models.User
	result := r.DB.First(&user, models.User{UserID: id})
	if result.Error != nil {
		return nil, false
	}
	return &user, true
}

func (r UserRepo) FindByEmail(email string) (*models.User, bool) {
	var user models.User
	result := r.DB.First(&user, models.User{Email: email})
	if result.Error != nil {
		// No user with that email exists
		return nil, false
	}
	return &user, true
}

// Gacha Repository =====================================================================================

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
