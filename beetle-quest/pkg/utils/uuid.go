package utils

import (
	"beetle-quest/pkg/models"

	"github.com/google/uuid"
)

func Parse(id string) (models.ApiUUID, error) {
	if uuid, err := uuid.Parse(id); err != nil {
		return models.ApiUUID{}, models.ErrCouldNotFindResourceByUUID
	} else {
		return models.ApiUUID(uuid), nil
	}
}
