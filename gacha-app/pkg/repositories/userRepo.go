package repositories

import (
	"gacha-app/pkg/models"
)

type UserRepo interface {
	ValidateUserID(userID models.UserId) bool
}
