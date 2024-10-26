package repositories

import (
	"beetle-quest/pkg/models"
)

type UserRepo interface {
	Create(email, username string, hashedPassword []byte) bool
	FindByUsername(username string) (*models.User, bool)
	FindByID(id models.UserId) (*models.User, bool)
	FindByUUID(uuid models.ApiUUID) (*models.User, bool)
}
