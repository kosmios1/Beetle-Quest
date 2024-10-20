package mock_repositories

import (
	"encoding/hex"
	"gacha-app/pkg/models"
)

type MockAuctionRepo struct {
	Auctions map[string]*models.Auction
}

func NewMockAuctionRepo() *MockAuctionRepo {
	return &MockAuctionRepo{
		Auctions: make(map[string]*models.Auction),
	}
}

func (m MockAuctionRepo) AddAuction(auction *models.Auction) error {
	hexId := hex.EncodeToString(auction.AuctionID)
	if _, ok := m.Auctions[hexId]; ok {
		return models.ErrAuctionAltreadyExists
	}

	m.Auctions[hexId] = auction
	return nil
}
