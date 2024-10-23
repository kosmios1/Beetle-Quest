package repositories

import (
	"beetle-quest/pkg/models"
)

type UserRepo interface {
	ValidateUserID(id models.UserId) bool
}
