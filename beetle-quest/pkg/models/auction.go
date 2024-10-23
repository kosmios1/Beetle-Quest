package models

import (
	"time"
)

type AuctionId EventId

type Auction struct {
	AuctionID AuctionId `json:"auction_id"`
	UUID      ApiUUID   `json:"uuid"`
	OwnerID   UserId    `json:"owner_id"`
	GachaID   GachaId   `json:"gacha_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	WinnerID  UserId    `json:"winner_id"`

	Blockchain *Blockchain `json:"blockchain"`
}

type Blockchain struct {
	Difficulty  int      `json:"difficulty"`
	GenesyBlock *Block   `json:"genesy_bid"`
	Chain       []*Block `json:"chain"`
}

type Block struct {
	Hash         []byte    `json:"hash"`
	PreviousHash []byte    `json:"previous_hash"`
	Timestamp    time.Time `json:"timestamp"`
	Pow          int       `json:"pow"`
	Bids         []*Bid    `json:"bids"`
}

type Bid struct {
	UserID      UserId `json:"owner_id"`
	AmountSpend int64  `json:"amount_spend"`
}
