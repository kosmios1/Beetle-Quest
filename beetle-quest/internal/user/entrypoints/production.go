//go:build !beetleQuestTest

package entrypoint

import (
	"beetle-quest/internal/user/controller"
	"beetle-quest/internal/user/service"

	urepo "beetle-quest/internal/user/repository"
	grepo "beetle-quest/pkg/repositories/impl/http/gacha"
	mrepo "beetle-quest/pkg/repositories/impl/http/market"
)

func NewUserController() *controller.UserController {
	return controller.NewUserController(
		service.NewUserService(
			urepo.NewUserRepo(),
			grepo.NewGachaRepo(),
			mrepo.NewMarketRepo(),
		),
	)
}
