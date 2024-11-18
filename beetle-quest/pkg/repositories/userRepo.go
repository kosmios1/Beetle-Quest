package repositories

import (
	"beetle-quest/pkg/models"
)

type UserRepo interface {
	GetAll() ([]models.User, error)

	Create(email, username string, hashedPassword []byte, currency int64) bool
	Update(user *models.User) bool
	Delete(user *models.User) bool

	FindByID(id models.UUID) (*models.User, bool)
	FindByEmail(email string) (*models.User, bool)
	FindByUsername(username string) (*models.User, bool)
}
