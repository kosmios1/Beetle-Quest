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

func (s *AuctionService) CreateAuction(ownerId models.UserId, gachaId models.GachaId, endTime time.Time) (models.Auction, error) {
	if !s.UserRepo.ValidateUserID(&ownerId) {
		return models.Auction{}, models.ErrInvalidUserID
	}

	if !s.GachaRepo.ValidateGachaID(&gachaId) {
		return models.Auction{}, models.ErrInvalidGachaID
	}

	auctionId, err := utils.GenerateRandomID(16)
	if err != nil {
		return models.Auction{}, err
	}

	startTime := time.Now()
	if startTime.After(endTime) {
		return models.Auction{}, models.ErrInvalidEndTime
	}

	genesy := models.Bid{
		Hash:         []byte{},
		PreviousHash: []byte{},
		BidData: models.BidData{
			UserID:      ownerId,
			CurrencyBid: 0,
		},
		Timestamp: time.Now(),
	}

	auction := models.Auction{
		AuctionID: auctionId,
		OwnerID:   ownerId,
		GachaID:   gachaId,
		StartTime: startTime,
		EndTime:   endTime,
		WinnerID:  models.UserId{},

		Difficulty: AuctionDifficulty,
		GenesyBid:  genesy,
		Biddings:   []*models.Bid{&genesy},
	}

	if err = s.AuctionRepo.AddAuction(&auction); err != nil {
		return models.Auction{}, err
	}

	return auction, err
}
