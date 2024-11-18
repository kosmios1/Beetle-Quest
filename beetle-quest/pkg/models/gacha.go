package models

import "errors"

type Rarity uint8

const (
	Common Rarity = iota
	Uncommon
	Rare
	Epic
	Legendary
)

func (r Rarity) String() string {
	return [...]string{"Common", "Uncommon", "Rare", "Epic", "Legendary"}[r]
}

func RarityFromString(r string) (Rarity, error) {
	switch r {
	case "Common":
		return Common, nil
	case "Uncommon":
		return Uncommon, nil
	case "Rare":
		return Rare, nil
	case "Epic":
		return Epic, nil
	case "Legendary":
		return Legendary, nil
	default:
		return Common, errors.New("Invalid rarity")
	}
}

type Gacha struct {
	GachaID   UUID   `json:"gacha_id"`
	Name      string `json:"name"`
	Rarity    Rarity `json:"rarity"`
	Price     int64  `json:"price"`
	ImagePath string `json:"image_path"`
}

type GachaUserRelation struct {
	UserID  UUID
	GachaID UUID
}
