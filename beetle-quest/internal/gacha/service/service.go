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
		return nil, models.ErrInvalidUUID
	}

	gacha, err := s.gachaRepo.FindByID(gachaID)
	if err != nil {
		return nil, err
	}
	return gacha, nil
}

func (s *GachaService) GetAll() ([]models.Gacha, error) {
	return s.gachaRepo.GetAll()
}

func (s *GachaService) AddGachaToUser(userID, gachaID models.UUID) error {
	return s.gachaRepo.AddGachaToUser(userID, gachaID)
}

func (s *GachaService) RemoveGachaFromUser(userID models.UUID, gachaID models.UUID) error {
	return s.gachaRepo.RemoveGachaFromUser(userID, gachaID)
}

func (s *GachaService) GetUserGachaDetails(userId, gachaId string) (*models.Gacha, error) {
	gachas, err := s.GetUserGachasStr(userId)
	if err != nil {
		return nil, err
	}

	gid, err := utils.ParseUUID(gachaId)
	if err != nil {
		return nil, models.ErrInvalidUUID
	}

	for _, gacha := range gachas {
		if gacha.GachaID == gid {
			return &gacha, nil
		}
	}
	return nil, models.ErrGachaNotFound
}

func (s *GachaService) GetUserGachasStr(userId string) ([]models.Gacha, error) {
	uid, err := utils.ParseUUID(userId)
	if err != nil {
		return nil, models.ErrInvalidUUID
	}
	return s.gachaRepo.GetUserGachas(uid)
}

func (s *GachaService) GetUserGachas(uid models.UUID) ([]models.Gacha, error) {
	return s.gachaRepo.GetUserGachas(uid)
}

func (s *GachaService) RemoveUserGachas(uid models.UUID) error {
	return s.gachaRepo.RemoveUserGachas(uid)
}

func (s *GachaService) CreateGacha(g *models.Gacha) error {
	return s.gachaRepo.Create(g)
}

func (s *GachaService) UpdateGacha(g *models.Gacha) error {
	return s.gachaRepo.Update(g)
}

func (s *GachaService) DeleteGacha(g *models.Gacha) error {
	return s.gachaRepo.Delete(g)
}
