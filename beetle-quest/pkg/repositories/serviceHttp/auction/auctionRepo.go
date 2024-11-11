package repository

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"errors"
	"net/http"

	"github.com/sony/gobreaker/v2"
)

var (
	findAuctionByIDEndpoint = utils.FindEnv("FIND_AUCTION_BY_ID_ENDPOINT")
)

type AuctionRepo struct {
	cb *gobreaker.CircuitBreaker[*http.Response]
}

func NewAuctionRepo() *AuctionRepo {
	return &AuctionRepo{
		cb: gobreaker.NewCircuitBreaker[*http.Response](gobreaker.Settings{}),
	}
}

func (a AuctionRepo) FindByID(models.UUID) (*models.Auction, bool) {
	return nil, false
}

func (a AuctionRepo) AddAuction(*models.Auction) error {
	return errors.New("Not implemented")
}
