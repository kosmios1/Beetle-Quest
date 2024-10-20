package repositories

import (
	"gacha-app/pkg/models"
)

type UserRepo interface {
	ValidateUserID(id *models.UserId) bool
}
