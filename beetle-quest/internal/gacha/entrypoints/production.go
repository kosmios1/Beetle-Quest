//go:build !beetleQuestTest

package entrypoint

import (
	"beetle-quest/internal/gacha/controller"
	"beetle-quest/internal/gacha/service"

	grepo "beetle-quest/internal/gacha/repository"
)

func NewGachaController() *controller.GachaController {
	return controller.NewGachaController(
		service.NewGachaService(
			grepo.NewGachaRepo(),
		),
	)
}
