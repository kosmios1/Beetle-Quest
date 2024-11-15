package models

import (
	"time"
)

type Auction struct {
	AuctionID UUID `json:"auction_id"`
	OwnerID   UUID `json:"owner_id"`
	GachaID   UUID `json:"gacha_id"`

	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`

	WinnerID UUID `json:"winner_id"`
}

type Bid struct {
	BidID       UUID      `json:"bid_id"`
	AuctionID   UUID      `json:"auction_id"`
	UserID      UUID      `json:"user_id"`
	AmountSpend int64     `json:"amount_spend"`
	TimeStamp   time.Time `json:"time_stamp"`
}
