package models

import (
	"time"
)

type AuctionId []byte

type Auction struct {
	AuctionID AuctionId `json:"auction_id"`
	OwnerID   UserId    `json:"owner_id"`
	GachaID   GachaId   `json:"gacha_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	WinnerID  UserId    `json:"winner_id"`

	Difficulty int    `json:"difficulty"`
	GenesyBid  Bid    `json:"genesy_bid"`
	Biddings   []*Bid `json:"biddings"`
}

type Bid struct {
	Hash         []byte    `json:"hash"`
	PreviousHash []byte    `json:"previous_hash"`
	Timestamp    time.Time `json:"timestamp"`
	Pow          int       `json:"pow"`

	BidData BidData `json:"bid_data"`
}

type BidData struct {
	UserID      UserId `json:"user_id"`
	CurrencyBid int    `json:"currency_bid"`
}
