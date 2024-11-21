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

type UserRepo struct {
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

func NewUserRepo() *UserRepo {
	var repo = &UserRepo{}
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
	return repo
}

func (r *UserRepo) GetAll() ([]models.User, error) {
	var users []models.User
	result := r.db.Table("users").Find(&users)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) || errors.Is(result.Error, gorm.ErrEmptySlice) {
			return nil, models.ErrUserNotFound
		}
		return nil, models.ErrInternalServerError
	}
	return users, nil
}

func (r *UserRepo) Create(email, username string, hashedPassword []byte, currency int64) error {
	result := r.db.Table("users").Create(&models.User{
		UserID:       utils.GenerateUUID(),
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		Currency:     currency,
	})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return models.ErrUsernameOrEmailAlreadyExists
		}
		return models.ErrInternalServerError
	}

	if result.RowsAffected == 0 {
		return models.ErrInternalServerError
	}
	return nil
}

func (r *UserRepo) Update(user *models.User) error {
	result := r.db.Table("users").Where("user_id = ?", user.UserID).Updates(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.ErrUserNotFound
		} else if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return models.ErrUsernameOrEmailAlreadyExists
		}
		return models.ErrInternalServerError
	}

	if result.RowsAffected == 0 {
		return models.ErrInternalServerError
	}
	return nil
}

func (r *UserRepo) Delete(user *models.User) error {
	result := r.db.Table("users").Delete(user, models.User{UserID: user.UserID})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.ErrUserNotFound
		}
		return models.ErrInternalServerError
	}

	if result.RowsAffected == 0 {
		return models.ErrInternalServerError
	}
	return nil
}

func (r *UserRepo) FindByUsername(username string) (*models.User, error) {
	var user models.User
	result := r.db.Table("users").First(&user, models.User{Username: username})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, models.ErrUserNotFound
		}
		return nil, models.ErrInternalServerError
	}
	return &user, nil
}

func (r *UserRepo) FindByID(id models.UUID) (*models.User, error) {
	var user models.User
	result := r.db.Table("users").First(&user, models.User{UserID: id})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, models.ErrUserNotFound
		}
		return nil, models.ErrInternalServerError
	}
	return &user, nil
}

func (r *UserRepo) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Table("users").First(&user, models.User{Email: email})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, models.ErrUserNotFound
		}
		return nil, models.ErrInternalServerError
	}
	return &user, nil
}
