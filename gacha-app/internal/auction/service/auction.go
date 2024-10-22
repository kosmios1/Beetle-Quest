package service

import (
	"gacha-app/pkg/models"
	"gacha-app/pkg/utils"
	"time"

	repo "gacha-app/pkg/repositories"
)

type AuctionService struct {
	UserRepo    repo.UserRepo
	GachaRepo   repo.GachaRepo
	AuctionRepo repo.AuctionRepo
}

const (
	AuctionDifficulty int = 10
)

func (s *AuctionService) CreateAuction(ownerId models.UserId, gachaId models.GachaId, endTime time.Time) (*models.Auction, error) {
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

	transactionId, err := utils.GenerateRandomID(8)
	if err != nil {
		return nil, models.ErrCouldNotGenerateAuction
	}

	genesy := &models.Block{
		Hash:         []byte{},
		PreviousHash: []byte{},
		Timestamp:    startTime,
		Pow:          0,
		Transactions: []*models.Transaction{
			{
				TransactionID: transactionId,
				Type:          models.Withdraw,
				UserID:        ownerId,
				Amount:        0,
				DateTime:      startTime,
				EventType:     models.AuctionEv,
				EventID:       auctionId,
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

func (s *AuctionService) GetAuction(auctionId models.AuctionId) (*models.Auction, error) {
	if !s.AuctionRepo.VaildateAuctionID(auctionId) {
		return &models.Auction{}, models.ErrInvalidAuctionID
	}

	return s.AuctionRepo.GetAuction(auctionId)
}
