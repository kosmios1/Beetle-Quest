package service

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/repositories"
	"beetle-quest/pkg/utils"
)

type MarketService struct {
	urepo repositories.UserRepo
	grepo repositories.GachaRepo
	arepo repositories.AuctionRepo
}

func NewMarketService(urepo repositories.UserRepo, grepo repositories.GachaRepo, arepo repositories.AuctionRepo) *MarketService {
	return &MarketService{urepo: urepo, grepo: grepo, arepo: arepo}
}

func (s *MarketService) AddBugsCoin(userId string, amount int64) error {
	id, err := utils.ParseUUID(userId)
	if err != nil {
		return models.ErrInvalidUserID
	}

	if amount <= 0 {
		return models.ErrAmountNotValid
	}

	user, ok := s.urepo.FindByID(id)
	if !ok {
		return models.ErrUserNotFound
	}

	user.Currency += amount
	if ok := s.urepo.Update(user); !ok {
		return models.ErrCouldNotUpdate
	}
	return nil
}

func (s *MarketService) BuyGacha(userId string, gachaId string) error {
	uid, err := utils.ParseUUID(userId)
	if err != nil {
		return models.ErrInvalidUserID
	}

	gid, err := utils.ParseUUID(gachaId)
	if err != nil {
		return models.ErrInvalidGachaID
	}

	// TODO: Check if user has already bought the gacha

	user, ok := s.urepo.FindByID(uid)
	if !ok {
		return models.ErrUserNotFound
	}

	gacha, ok := s.grepo.FindByID(gid)
	if !ok {
		return models.ErrGachaNotFound
	}

	if user.Currency < gacha.Price {
		return models.ErrNotEnoughMoneyToBuyGacha
	}

	// TODO: How do i make this step consistent?
	// Compensating transaction if one fails?
	user.Currency -= gacha.Price
	if ok := s.urepo.Update(user); !ok {
		return models.ErrCouldNotUpdate
	}

	if ok := s.grepo.AddGachaToUser(uid, gid); !ok {
		return models.ErrCouldNotAddGachaToUser
	}
	return nil
}
