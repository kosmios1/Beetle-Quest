package service

import "beetle-quest/pkg/repositories"

type MarketService struct {
	repositories.UserRepo
	repositories.GachaRepo
	repositories.AuctionRepo
}
