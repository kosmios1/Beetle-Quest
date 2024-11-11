package service

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/repositories"
	"beetle-quest/pkg/utils"
)

type GachaService struct {
	gachaRepo repositories.GachaRepo
}

func NewGachaService(repo repositories.GachaRepo) *GachaService {
	return &GachaService{
		gachaRepo: repo,
	}
}

func (s *GachaService) FindByID(id string) (*models.Gacha, error) {
	gachaID, err := utils.ParseUUID(id)
	if err != nil {
		return nil, models.ErrInvalidGachaID
	}

	gacha, ok := s.gachaRepo.FindByID(gachaID)
	if !ok {
		return nil, models.ErrCouldNotFindResourceByUUID
	}

	return gacha, nil
}

func (s *GachaService) GetAll() ([]models.Gacha, bool) {
	return s.gachaRepo.GetAll()
}
