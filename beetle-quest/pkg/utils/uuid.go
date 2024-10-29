package utils

import (
	"beetle-quest/pkg/models"

	"github.com/google/uuid"
)

func Parse(id string) (models.UUID, error) {
	if uuid, err := uuid.Parse(id); err != nil {
		return models.UUID{}, models.ErrCouldNotFindResourceByUUID
	} else {
		return models.UUID(uuid), nil
	}
}

func GenerateUUID() models.UUID {
	return models.UUID(uuid.New())
}
