package models

import (
	"fmt"

	"github.com/google/uuid"
)

type UUID uuid.UUID

func (uid *UUID) Scan(value interface{}) error {
	var id uuid.UUID
	var err error
	switch v := value.(type) {
	case string:
		id, err = uuid.Parse(v)
	case []byte:
		id, err = uuid.FromBytes(v)
	default:
		return fmt.Errorf("unsupported type for UserId: %T", value)
	}
	if err != nil {
		return err
	}
	*uid = UUID(id)
	return nil
}

type EventType uint8

const (
	AuctionEv EventType = iota
	MarketEv
	GameEv
)
