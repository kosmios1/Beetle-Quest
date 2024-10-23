package models

type ApiUUID [16]byte
type EventId []byte

type EventType uint8

const (
	AuctionEv EventType = iota
	MarketEv
	GameEv
)
