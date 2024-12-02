//go:build beetleQuestTest

package entrypoint

import (
	"beetle-quest/internal/auth/controller"
	"beetle-quest/internal/auth/service"

	arepo "beetle-quest/pkg/repositories/impl/mock/admin"
	urepo "beetle-quest/pkg/repositories/impl/mock/user"
)

func NewAuthController() (*controller.AuthController, error) {
	return controller.NewAuthController(
		service.NewAuthService(
			urepo.NewUserRepo(),
			arepo.NewAdminRepo(),
		),
	), nil
}
