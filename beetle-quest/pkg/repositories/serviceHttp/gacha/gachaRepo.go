package repository

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"net/http"

	"github.com/sony/gobreaker/v2"
)

var (
	getAllEndpoint        = utils.FindEnv("GET_ALL_GACHA_ENDPOINT")
	findGachaByIDEndpoint = utils.FindEnv("FIND_GACHA_BY_ID_ENDPOINT")
)

type GachaRepo struct {
	cb *gobreaker.CircuitBreaker[*http.Response]
}

func NewGachaRepo() *GachaRepo {
	return &GachaRepo{
		cb: gobreaker.NewCircuitBreaker[*http.Response](gobreaker.Settings{}),
	}
}

func (r GachaRepo) GetAll() ([]models.Gacha, bool) {
	return nil, false
}

func (r GachaRepo) FindByID(models.UUID) (*models.Gacha, bool) {
	return nil, false
}
