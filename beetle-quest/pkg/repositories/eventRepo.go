package repositories

import "beetle-quest/pkg/models"

type EventRepo interface {
	AddEndAuctionEvent(auction *models.Auction) error
	StartSubscriber(callback func(aid models.UUID), errCallback func(err error))
}
