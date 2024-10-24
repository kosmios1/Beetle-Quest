package service

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"time"

	repo "beetle-quest/pkg/repositories"
)

type AuctionService struct {
	UserRepo    repo.UserRepo
	GachaRepo   repo.GachaRepo
	AuctionRepo repo.AuctionRepo
}

const (
	AuctionDifficulty int = 10
)

func (s *AuctionService) CreateAuction(ownerUUID models.ApiUUID, gachaUUID models.ApiUUID, endTime time.Time) (*models.Auction, error) {
	user, ok := s.UserRepo.FindByUUID(ownerUUID)
	if !ok {
		return nil, models.ErrInvalidUserID
	}

	gacha, ok := s.GachaRepo.FindByUUID(gachaUUID)
	if !ok {
		return nil, models.ErrInvalidGachaID
	}

	auctionID, err := utils.GenerateRandomID(16)
	if err != nil {
		return nil, err
	}

	startTime := time.Now()
	if startTime.After(endTime) {
		return nil, models.ErrInvalidEndTime
	}

	auctionUUID := utils.GenerateUUID()

	genesy := &models.Block{
		Hash:         []byte{},
		PreviousHash: []byte{},
		Timestamp:    startTime,
		Pow:          0,
		Bids: []*models.Bid{
			{
				UserID:      user.UserID,
				AmountSpend: 0,
			},
		},
	}

	auction := &models.Auction{
		AuctionID: auctionID,
		UUID:      auctionUUID,
		OwnerID:   user.UserID,
		GachaID:   gacha.GachaID,
		StartTime: startTime,
		EndTime:   endTime,
		WinnerID:  nil,

		Blockchain: &models.Blockchain{
			Difficulty:  AuctionDifficulty,
			GenesyBlock: genesy,
			Chain:       []*models.Block{genesy},
		},
	}

	if err = s.AuctionRepo.AddAuction(auction); err != nil {
		return nil, err
	}

	return auction, err
}

func (s *AuctionService) GetAuction(auctionUUID models.ApiUUID) (*models.Auction, error) {
	auction, ok := s.AuctionRepo.FindByUUID(auctionUUID)
	if !ok {
		return &models.Auction{}, models.ErrInvalidAuctionID
	}
	return auction, nil
}
