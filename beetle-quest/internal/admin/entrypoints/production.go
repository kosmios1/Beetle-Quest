//go:build !beetleQuestTest

package entrypoint

import (
	"beetle-quest/internal/admin/controller"
	"beetle-quest/internal/admin/service"

	arepo "beetle-quest/internal/admin/repository"
	grepo "beetle-quest/pkg/repositories/impl/http/gacha"
	mrepo "beetle-quest/pkg/repositories/impl/http/market"
	urepo "beetle-quest/pkg/repositories/impl/http/user"
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
