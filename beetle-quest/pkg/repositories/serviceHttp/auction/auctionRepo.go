package repository

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"net/http"

	"github.com/sony/gobreaker/v2"
)

var (
	findAuctionByIDEndpoint = utils.FindEnv("FIND_AUCTION_BY_ID_ENDPOINT")
)

type AuctionRepo struct {
	cb *gobreaker.CircuitBreaker[*http.Response]
}

func NewGachaRepo() *AuctionRepo {
	return &AuctionRepo{
		cb: gobreaker.NewCircuitBreaker[*http.Response](gobreaker.Settings{}),
	}
}

func (r AuctionRepo) FindByID(aid models.UUID) (*models.Auction, bool) {
	return nil, false
}

func (r AuctionRepo) AddAuction(auction *models.Auction) error {
	return nil
}
