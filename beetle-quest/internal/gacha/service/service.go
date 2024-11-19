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

func (s *GachaService) GetUserGachaDetails(userId, gachaId string) (models.Gacha, bool) {
	gachas, ok := s.GetUserGachasStr(userId)
	if !ok {
		return models.Gacha{}, false
	}

	gid, err := utils.ParseUUID(gachaId)
	if err != nil {
		return models.Gacha{}, false
	}

	for _, gacha := range gachas {
		if gacha.GachaID == gid {
			return gacha, true
		}
	}
	return models.Gacha{}, false
}

func (s *GachaService) GetUserGachasStr(userId string) ([]models.Gacha, bool) {
	uid, err := utils.ParseUUID(userId)
	if err != nil {
		return []models.Gacha{}, false
	}

	gachas, ok := s.gachaRepo.GetUserGachas(uid)
	if !ok {
		return []models.Gacha{}, false
	}
	return gachas, true
}

func (s *GachaService) GetUserGachas(uid models.UUID) ([]models.Gacha, bool) {
	gachas, ok := s.gachaRepo.GetUserGachas(uid)
	if !ok {
		return []models.Gacha{}, false
	}
	return gachas, true
}

func (s *GachaService) RemoveUserGachas(uid models.UUID) bool {
	return s.gachaRepo.RemoveUserGachas(uid)
}

func (s *GachaService) CreateGacha(g *models.Gacha) error {
	if ok := s.gachaRepo.Create(g); !ok {
		return models.ErrCouldNotCreateGacha
	}
	return nil
}

func (s *GachaService) UpdateGacha(g *models.Gacha) error {
	if ok := s.gachaRepo.Update(g); !ok {
		return models.ErrCouldNotUpdateGacha
	}
	return nil
}

func (s *GachaService) DeleteGacha(g *models.Gacha) error {
	if ok := s.gachaRepo.Delete(g); !ok {
		return models.ErrCouldNotDeleteGacha
	}
	return nil
}
