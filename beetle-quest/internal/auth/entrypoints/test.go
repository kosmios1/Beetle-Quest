//go:build beetleQuestTest

package entrypoint

import (
	"beetle-quest/internal/auth/controller"
	"beetle-quest/internal/auth/service"

	arepo "beetle-quest/pkg/repositories/impl/mock/admin"
	srepo "beetle-quest/pkg/repositories/impl/mock/session"
	urepo "beetle-quest/pkg/repositories/impl/mock/user"
)

func NewAuthController() *controller.AuthController {
	return controller.NewAuthController(
		service.NewAuthService(
			urepo.NewUserRepo(),
			arepo.NewAdminRepo(),
			srepo.NewSessionRepo(),
		),
	)
}
