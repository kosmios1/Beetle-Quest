package repositories

import (
	"beetle-quest/pkg/models"
)

type UserRepo interface {
	FindByID(id models.UserId) (*models.User, bool)
	FindByUUID(uuid models.ApiUUID) (*models.User, bool)
}
