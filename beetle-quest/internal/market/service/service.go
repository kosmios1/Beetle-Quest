package service

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/repositories"
	"beetle-quest/pkg/utils"
	"math/rand/v2"
	"time"
)

type MarketService struct {
	urepo repositories.UserRepo
	grepo repositories.GachaRepo
	arepo repositories.MarketRepo
}

func NewMarketService(urepo repositories.UserRepo, grepo repositories.GachaRepo, arepo repositories.MarketRepo) *MarketService {
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

func (s *MarketService) RollGacha(userId string) (string, error) {
	uid, err := utils.ParseUUID(userId)
	if err != nil {
		return "", models.ErrInvalidUserID
	}

	user, exists := s.urepo.FindByID(uid)
	if !exists {
		return "", models.ErrUserNotFound
	}

	if user.Currency < 1000 {
		return "", models.ErrNotEnoughMoneyToRollGacha
	}

	gachas, ok := s.grepo.GetAll()
	if !ok {
		return "", models.ErrInternalServerError
	}
	gid := gachas[rand.IntN(len(gachas))].GachaID

	user.Currency -= 1000
	if ok := s.urepo.Update(user); !ok {
		return "", models.ErrInternalServerError
	}

	gachas, ok = s.grepo.GetUserGachas(uid)
	if !ok {
		user.Currency += 1000
		_ = s.urepo.Update(user)
		// TODO: What do i do here if it fails?
		// - Report to admin
		return "", models.ErrInternalServerError
	}

	for _, gacha := range gachas {
		if gacha.GachaID == gid {
			return "Opps you already have this gacha!", nil
		}
	}

	if ok := s.grepo.AddGachaToUser(uid, gid); !ok {
		user.Currency += 1000
		_ = s.urepo.Update(user)
		// TODO: What do i do here?
		// - Report to admin
		return "", models.ErrCouldNotAddGachaToUser
	}

	return "Gacha successfully obtained, check your inventory!", nil
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

	userGacha, ok := s.grepo.GetUserGachas(uid)
	if !ok {
		return models.ErrUserAlreadyHasGacha
	}

	for _, gacha := range userGacha {
		if gid == gacha.GachaID {
			return models.ErrUserAlreadyHasGacha
		}
	}

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

	// TODO: Create transaction

	user.Currency -= gacha.Price
	if ok := s.urepo.Update(user); !ok {
		return models.ErrCouldNotUpdate
	}

	if ok := s.grepo.AddGachaToUser(uid, gid); !ok {
		// Compensating transaction
		user.Currency += gacha.Price
		if ok := s.urepo.Update(user); !ok {
			// TODO: What do i do here?
			// - Report to admin
		}
		return models.ErrCouldNotAddGachaToUser
	}
	return nil
}

func (s *MarketService) CreateAuction(userId, gachaId string, endTime time.Time) error {
	uid, err := utils.ParseUUID(userId)
	if err != nil {
		return models.ErrInvalidUserID
	}

	gid, err := utils.ParseUUID(gachaId)
	if err != nil {
		return models.ErrInvalidGachaID
	}

	user, exists := s.urepo.FindByID(uid)
	if !exists {
		return models.ErrUserNotFound
	}

	gacha, exists := s.grepo.FindByID(gid)
	if !exists {
		return models.ErrGachaNotFound
	}

	gachas, ok := s.grepo.GetUserGachas(uid)
	if !ok {
		return models.ErrCouldNotRetrieveUserGachas
	}

	var found bool
	for _, g := range gachas {
		if g.GachaID == gid {
			found = true
			break
		}
	}

	if !found {
		return models.ErrUserDoesNotOwnGacha
	}

	auctions, ok := s.arepo.GetUserAuctions(uid)
	if !ok {
		return models.ErrCouldNotRetrieveUserAuctions
	}

	for _, a := range auctions {
		if a.GachaID == gid {
			return models.ErrGachaAlreadyAuctioned
		}
	}

	// TODO: time is not correct inside containers
	startTime := time.Now()
	if endTime.Before(startTime) || endTime.After(startTime.Add(time.Hour*24)) {
		return models.ErrInvalidEndTime
	}

	var auction models.Auction = models.Auction{
		AuctionID: utils.GenerateUUID(),
		OwnerID:   user.UserID,
		GachaID:   gacha.GachaID,
		StartTime: time.Now(),
		EndTime:   endTime,
		WinnerID:  models.UUID{},
	}

	if ok := s.arepo.Create(&auction); !ok {
		return models.ErrCouldNotCreateAuction
	}

	return nil
}

func (s *MarketService) RetrieveAuctionTemplateList() ([]models.AuctionTemplate, error) {
	auctions, ok := s.arepo.GetAll()
	if !ok {
		return nil, models.ErrRetrievingAuctions
	}

	var data []models.AuctionTemplate = []models.AuctionTemplate{}
	for _, auction := range auctions {
		gacha, exists := s.grepo.FindByID(auction.GachaID)
		if !exists {
			return nil, models.ErrUserNotFound
		}

		owner, exists := s.urepo.FindByID(auction.OwnerID)
		if !exists {
			return nil, models.ErrGachaNotFound
		}

		data = append(data, models.AuctionTemplate{
			Auction:       auction,
			GachaName:     gacha.Name,
			ImagePath:     gacha.ImagePath,
			OwnerUsername: owner.Username,
		})
	}

	return data, nil
}

func (s *MarketService) FindByID(auctionId string) (*models.Auction, bool) {
	aid, err := utils.ParseUUID(auctionId)
	if err != nil {
		return &models.Auction{}, false
	}

	auction, exists := s.arepo.FindByID(aid)
	if !exists {
		return &models.Auction{}, false
	}
	return auction, true
}

func (s *MarketService) GetBidListOfAuctionID(auctionId string) ([]models.Bid, bool) {
	aid, err := utils.ParseUUID(auctionId)
	if err != nil {
		return []models.Bid{}, false
	}

	bids, ok := s.arepo.GetBidListOfAuction(aid)
	if !ok {
		return []models.Bid{}, false
	}
	return bids, true
}
