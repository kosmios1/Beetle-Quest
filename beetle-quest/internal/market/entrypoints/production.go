//go:build !beetleQuestTest

package entrypoint

import (
	"beetle-quest/internal/market/controller"
	"beetle-quest/internal/market/service"
	"beetle-quest/pkg/models"

	evrepo "beetle-quest/internal/market/repository"
	mrepo "beetle-quest/internal/market/repository"
	grepo "beetle-quest/pkg/repositories/impl/http/gacha"
	urepo "beetle-quest/pkg/repositories/impl/http/user"
)

func NewMarketController() (*controller.MarketController, error) {
 	evrepo, err := evrepo.NewEventRepo()
	if err != nil {
		return nil, models.ErrInternalServerError
	}
	return controller.NewMarketController(
		service.NewMarketService(
			urepo.NewUserRepo(),
			grepo.NewGachaRepo(),
			mrepo.NewMarketRepo(),
			evrepo,
		),
	), nil
}
