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

func (s *AuctionService) CreateAuction(ownerId models.ApiUUID, gachaId models.ApiUUID, endTime time.Time) (*models.Auction, error) {
	// TODO: Convert ApiUUID to UserID and gachaID
	if !s.UserRepo.ValidateUserID(ownerId) {
		return nil, models.ErrInvalidUserID
	}

	if !s.GachaRepo.ValidateGachaID(gachaId) {
		return nil, models.ErrInvalidGachaID
	}

	auctionId, err := utils.GenerateRandomID(16)
	if err != nil {
		return nil, err
	}

	startTime := time.Now()
	if startTime.After(endTime) {
		return nil, models.ErrInvalidEndTime
	}

	genesy := &models.Block{
		Hash:         []byte{},
		PreviousHash: []byte{},
		Timestamp:    startTime,
		Pow:          0,
		Bids: []*models.Bid{
			{
				UserID:      ownerId,
				AmountSpend: 0,
			},
		},
	}

	auction := &models.Auction{
		AuctionID: auctionId,
		OwnerID:   ownerId,
		GachaID:   gachaId,
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

func (s *AuctionService) GetAuction(auctionId models.ApiUUID) (*models.Auction, error) {
	// TODO: Convert ApiUUID to UserID and gachaID
	if !s.AuctionRepo.VaildateAuctionID(auctionId) {
		return &models.Auction{}, models.ErrInvalidAuctionID
	}

	return s.AuctionRepo.GetAuction(auctionId)
}
