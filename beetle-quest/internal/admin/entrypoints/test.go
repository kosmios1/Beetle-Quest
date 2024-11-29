//go:build beetleQuestTest

package entrypoint

import (
	"beetle-quest/internal/admin/controller"
	"beetle-quest/internal/admin/service"

	arepo "beetle-quest/pkg/repositories/impl/mock/admin"
	grepo "beetle-quest/pkg/repositories/impl/mock/gacha"
	mrepo "beetle-quest/pkg/repositories/impl/mock/market"
	urepo "beetle-quest/pkg/repositories/impl/mock/user"
)

func NewAdminController() *controller.AdminController {
	return controller.NewAdminController(
		service.NewAdminService(
			arepo.NewAdminRepo(),
			mrepo.NewMarketRepo(),
			urepo.NewUserRepo(),
			grepo.NewGachaRepo(),
		),
	)
}
