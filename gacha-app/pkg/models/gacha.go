package models

type GachaId []byte

type Rarity uint8

const (
	Common Rarity = iota
	Uncommon
	Rare
	Epic
	Legendary
)

type Gacha struct {
	GachaID GachaId `json:"gacha_id"`
	Name    string  `json:"name"`
	Rarity  Rarity  `json:"rarity"`
	Price   uint64  `json:"price"`
}
