package repositories

import "beetle-quest/pkg/models"

type AdminRepo interface {
	FindByID(models.UUID) (*models.Admin, bool)
}
