//go:build beetleQuestTest

package entrypoint

import (
	"beetle-quest/internal/user/controller"
	"beetle-quest/internal/user/service"

	grepo "beetle-quest/pkg/repositories/impl/mock/gacha"
	mrepo "beetle-quest/pkg/repositories/impl/mock/market"
	urepo "beetle-quest/pkg/repositories/impl/mock/user"
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
