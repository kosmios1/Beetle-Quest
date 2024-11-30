package repository

import (
	"beetle-quest/pkg/models"
	"time"
)

type EventRepo struct {
	closeEventChan chan models.UUID
}

func NewEventRepo() (*EventRepo, error) {
	return &EventRepo{
		closeEventChan: make(chan models.UUID),
	}, nil
}

func (r *EventRepo) AddEndAuctionEvent(auction *models.Auction) error {
	go func() {
		time.Sleep(time.Until(auction.EndTime))
		r.closeEventChan <- auction.AuctionID
	}()
	return nil
}

// NOTE: We don't call errCallback because we don't have any error possible in this implementation
func (r *EventRepo) StartSubscriber(callback func(aid models.UUID), errCallback func(err error)) {
	for {
		aid := <-r.closeEventChan
		callback(aid)
	}
}
