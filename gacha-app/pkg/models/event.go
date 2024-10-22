package models

type EventId []byte

type EventType uint8

const (
	AuctionEv EventType = iota
	MarketEv
	GameEv
)
