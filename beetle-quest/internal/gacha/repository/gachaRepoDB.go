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

	return repo
}

func (r *GachaRepo) Create(gacha *models.Gacha) error {
	result := r.db.Table("gachas").Create(gacha)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return models.ErrGachaAlreadyExists
		}
		return models.ErrInternalServerError
	}

	if result.RowsAffected == 0 {
		return models.ErrInternalServerError
	}
	return nil
}

func (r *GachaRepo) Update(gacha *models.Gacha) error {
	result := r.db.Table("gachas").Where("gacha_id = ?", gacha.GachaID).Updates(gacha)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.ErrGachaNotFound
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

func (r *GachaRepo) Delete(gacha *models.Gacha) error {
	result := r.db.Table("gachas").Delete(gacha, models.Gacha{GachaID: gacha.GachaID})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.ErrGachaNotFound
		}
		return models.ErrInternalServerError
	}

	if result.RowsAffected == 0 {
		return models.ErrInternalServerError
	}
	return nil
}

func (r *GachaRepo) FindByID(id models.UUID) (*models.Gacha, error) {
	var gacha models.Gacha
	result := r.db.Table("gachas").First(&gacha, models.Gacha{GachaID: id})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, models.ErrGachaNotFound
		}
		return nil, models.ErrInternalServerError
	}
	return &gacha, nil
}

func (r *GachaRepo) GetAll() ([]models.Gacha, error) {
	var gachas []models.Gacha
	result := r.db.Table("gachas").Find(&gachas)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, models.ErrGachaNotFound
		}
		return nil, models.ErrInternalServerError
	}
	return gachas, nil
}

func (r *GachaRepo) AddGachaToUser(uid models.UUID, gid models.UUID) error {
	value := models.GachaUserRelation{UserID: uid, GachaID: gid}
	result := r.db.Table("user_gacha").Create(value)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return models.ErrUserAlreadyHasGacha
		}
		return models.ErrInternalServerError
	}

	if result.RowsAffected == 0 {
		return models.ErrInternalServerError
	}
	return nil
}

func (r *GachaRepo) RemoveGachaFromUser(uid models.UUID, gid models.UUID) error {
	value := models.GachaUserRelation{UserID: uid, GachaID: gid}
	result := r.db.Table("user_gacha").Delete(value, "user_id = ? AND gacha_id = ?", uid, gid)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.ErrRetalationGachaUserNotFound
		}
		return models.ErrInternalServerError
	}

	if result.RowsAffected == 0 {
		return models.ErrInternalServerError
	}
	return nil
}

func (r *GachaRepo) RemoveUserGachas(uid models.UUID) error {
	result := r.db.Table("user_gacha").Delete(models.GachaUserRelation{}, models.GachaUserRelation{UserID: uid})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.ErrRetalationGachaUserNotFound
		}
		return models.ErrInternalServerError
	}
	return nil
}

func (r *GachaRepo) GetUserGachas(uid models.UUID) ([]models.Gacha, error) {
	var gachas []models.Gacha
	result := r.db.Table("user_gacha").Select("gachas.*").Joins("JOIN gachas ON gachas.gacha_id = user_gacha.gacha_id").Where("user_gacha.user_id = ?", uid).Find(&gachas)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, models.ErrGachaNotFound
		}
		return nil, models.ErrInternalServerError
	}
	return gachas, nil
}
