package models

import "time"

type TransactionId []byte

type TransactionType uint8

const (
	Deposit TransactionType = iota
	Withdraw
)

type Transaction struct {
	TransactionID TransactionId   `json:"transaction_id"`
	UUID          ApiUUID         `json:"uuid"`
	Type          TransactionType `json:"transaction_type"`
	UserID        UserId          `json:"user_id"`
	Amount        uint64          `json:"amount"`
	DateTime      time.Time       `json:"date_time"`
	EventType     EventType       `json:"event_type"`
	EventID       EventId         `json:"event_id"`
}
