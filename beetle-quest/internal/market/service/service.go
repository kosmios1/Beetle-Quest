package service

import "beetle-quest/pkg/repositories"

type MarketService struct {
	urepo repositories.UserRepo
	grepo repositories.GachaRepo
	arepo repositories.AuctionRepo
}

func NewMarketService(urepo repositories.UserRepo, grepo repositories.GachaRepo, arepo repositories.AuctionRepo) *MarketService {
	return &MarketService{urepo: urepo, grepo: grepo, arepo: arepo}
}
