//go:build beetleQuestTest

package entrypoint

import (
	"beetle-quest/internal/auth/controller"
	"beetle-quest/internal/auth/service"

	arepo "beetle-quest/pkg/repositories/impl/mock/admin"
	urepo "beetle-quest/pkg/repositories/impl/mock/user"

	"github.com/go-session/session/v3"
)

func NewAuthController() (*controller.AuthController, error) {
	session.InitManager(
		session.SetStore(session.NewMemoryStore()),
	)
	return controller.NewAuthController(
		service.NewAuthService(
			urepo.NewUserRepo(),
			arepo.NewAdminRepo(),
		),
	), nil
}
