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

type UserRepo struct {
	DB *gorm.DB
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
