package models

type Rarity string

const (
	Common    Rarity = "COMMON"
	Uncommon  Rarity = "UNCOMMON"
	Rare      Rarity = "RARE"
	Epic      Rarity = "EPIC"
	Legendary Rarity = "LEGENDARY"
)

type Gacha struct {
	GachaID   UUID   `json:"gacha_id"`
	Name      string `json:"name"`
	Rarity    Rarity `json:"rarity"`
	Price     uint64 `json:"price"`
	ImagePath string `json:"image_path"`
}
