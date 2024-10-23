package mock_repositories

import (
	"beetle-quest/pkg/models"
	"encoding/hex"
	"time"
)

type MockAuctionRepo struct{}

var Auctions = map[string]*models.Auction{
	"6175637469": { // "aucti" in hex
		AuctionID: models.AuctionId("6175637469"),
		OwnerID:   models.UserId("4b6974"),
		GachaID:   models.GachaId("6761636861"),
		StartTime: time.Now(),
		EndTime:   time.Now().Add(48 * time.Hour),
		WinnerID:  []byte(""),

		Blockchain: &models.Blockchain{
			Difficulty: 10,
			GenesyBlock: &models.Block{
				Hash:         []byte{},
				PreviousHash: []byte{},
				Timestamp:    time.Now(),
				Pow:          0,
			},
			Chain: []*models.Block{},
		},
	},
}

func (m MockAuctionRepo) GetAuction(auctionID models.AuctionId) (*models.Auction, error) {
	hexId := hex.EncodeToString(auctionID)
	if auction, ok := Auctions[hexId]; ok {
		return auction, nil
	}
	return nil, models.ErrAuctionNotFound
}

func (m MockAuctionRepo) VaildateAuctionID(auctionID models.AuctionId) bool {
	hexId := hex.EncodeToString(auctionID)
	_, ok := Auctions[hexId]
	return ok
}

func (m MockAuctionRepo) AddAuction(auction *models.Auction) error {
	hexId := hex.EncodeToString(auction.AuctionID)
	if _, ok := Auctions[hexId]; ok {
		return models.ErrAuctionAltreadyExists
	}

	Auctions[hexId] = auction
	return nil
}
