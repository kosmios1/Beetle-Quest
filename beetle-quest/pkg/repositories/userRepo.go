package repositories

import (
	"beetle-quest/pkg/models"
)

type UserRepo interface {
	GetAll() ([]models.User, error)

	Create(email, username string, hashedPassword []byte, currency int64) error
	Update(user *models.User) error
	Delete(user *models.User) error

	FindByID(id models.UUID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
}
