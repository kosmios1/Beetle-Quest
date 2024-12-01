//go:build beetleQuestTest

package entrypoint

import (
	"beetle-quest/internal/market/controller"
	"beetle-quest/internal/market/service"

	evrepo "beetle-quest/pkg/repositories/impl/mock/event"
	grepo "beetle-quest/pkg/repositories/impl/mock/gacha"
	mrepo "beetle-quest/pkg/repositories/impl/mock/market"
	urepo "beetle-quest/pkg/repositories/impl/mock/user"
	"beetle-quest/pkg/utils"
)

func NewMarketController() (*controller.MarketController, error) {
	return controller.NewMarketController(
		service.NewMarketService(
			urepo.NewUserRepo(),
			grepo.NewGachaRepo(),
			mrepo.NewMarketRepo(),
			utils.PanicIfError[*evrepo.EventRepo](evrepo.NewEventRepo()),
		),
	), nil
}
