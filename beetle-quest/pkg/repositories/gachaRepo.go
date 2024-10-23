package repositories

import "beetle-quest/pkg/models"

type GachaRepo interface {
	ValidateGachaID(models.GachaId) bool
}
