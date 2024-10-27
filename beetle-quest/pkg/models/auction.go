package models

import (
	"time"
)

type Auction struct {
	AuctionID UUID      `json:"auction_id"`
	OwnerID   UUID      `json:"owner_id"`
	GachaID   UUID      `json:"gacha_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	WinnerID  UUID      `json:"winner_id"`

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
	UserID      UUID  `json:"owner_id"`
	AmountSpend int64 `json:"amount_spend"`
}
