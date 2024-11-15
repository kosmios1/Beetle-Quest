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

func (s *GachaService) FindByID(id string) (*models.Gacha, bool) {
	gachaID, err := utils.ParseUUID(id)
	if err != nil {
		return nil, false
	}

	gacha, ok := s.gachaRepo.FindByID(gachaID)
	if !ok {
		return nil, false
	}

	return gacha, true
}

func (s *GachaService) GetAll() ([]models.Gacha, bool) {
	return s.gachaRepo.GetAll()
}

func (s *GachaService) AddGachaToUser(userID, gachaID models.UUID) error {
	if ok := s.gachaRepo.AddGachaToUser(userID, gachaID); !ok {
		return models.ErrCouldNotAddGachaToUser
	}
	return nil
}

func (s *GachaService) RemoveGachaFromUser(userID models.UUID, gachaID models.UUID) error {
	if ok := s.gachaRepo.RemoveGachaFromUser(userID, gachaID); !ok {
		return models.ErrCouldNotRemoveGachaFromUser
	}
	return nil
}

func (s *GachaService) GetUserGachas(uid models.UUID) ([]models.Gacha, bool) {
	gachas, ok := s.gachaRepo.GetUserGachas(uid)
	if !ok {
		return []models.Gacha{}, false
	}
	return gachas, true
}
