//go:build !beetleQuestTest

package entrypoint

import (
	"beetle-quest/internal/auth/controller"
	"beetle-quest/internal/auth/service"

	srepo "beetle-quest/internal/auth/repository"
	arepo "beetle-quest/pkg/repositories/impl/http/admin"
	urepo "beetle-quest/pkg/repositories/impl/http/user"
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
