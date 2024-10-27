package models

import "time"

type TransactionType uint8

const (
	Deposit TransactionType = iota
	Withdraw
)

type Transaction struct {
	TransactionID UUID            `json:"transaction_id"`
	Type          TransactionType `json:"transaction_type"`
	UserID        UUID            `json:"user_id"`
	Amount        uint64          `json:"amount"`
	DateTime      time.Time       `json:"date_time"`
	EventType     EventType       `json:"event_type"`
	EventID       UUID            `json:"event_id"`
}
