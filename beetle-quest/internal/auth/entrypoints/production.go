//go:build !beetleQuestTest

package entrypoint

import (
	"beetle-quest/internal/auth/controller"
	"beetle-quest/internal/auth/service"
	"beetle-quest/pkg/models"

	srepo "beetle-quest/internal/auth/repository"
	arepo "beetle-quest/pkg/repositories/impl/http/admin"
	urepo "beetle-quest/pkg/repositories/impl/http/user"
)

func NewAuthController() (*controller.AuthController, error) {
  sesRepo, err := srepo.NewSessionRepo()
 	if err != nil {
		return nil, models.ErrInternalServerError
	}
	return controller.NewAuthController(
		service.NewAuthService(
			urepo.NewUserRepo(),
			arepo.NewAdminRepo(),
			sesRepo),
	), nil
}
