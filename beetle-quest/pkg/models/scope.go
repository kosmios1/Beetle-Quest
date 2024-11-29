package models

type Scope string

const (
	UserScope   Scope = "user"
	GachaScope  Scope = "gacha"
	MarketScope Scope = "market"
	AdminScope  Scope = "admin"
)
