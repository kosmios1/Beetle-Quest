package models

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

type Gacha struct {
	GachaID   UUID   `json:"gacha_id"`
	Name      string `json:"name"`
	Rarity    Rarity `json:"rarity"`
	Price     int64  `json:"price"`
	ImagePath string `json:"image_path"`
}
