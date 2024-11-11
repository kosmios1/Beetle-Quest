package repository

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/sony/gobreaker/v2"
)

var (
	getAllEndpoint        = utils.FindEnv("GET_ALL_GACHA_ENDPOINT")
	findGachaByIDEndpoint = utils.FindEnv("FIND_GACHA_BY_ID_ENDPOINT")

	addGachaToUserEndpoint = utils.FindEnv("ADD_GACHA_TO_USER_ENDPOINT")
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

func (r GachaRepo) FindByID(gid models.UUID) (*models.Gacha, bool) {
	requestData := models.FindGachaByIDData{
		GachaID: gid.String(),
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, false
	}

	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := http.Post(
			findGachaByIDEndpoint,
			"application/json",
			bytes.NewBuffer(jsonData),
		)

		if err != nil {
			return nil, err
		}

		return resp, nil
	})
	defer resp.Body.Close()

	if err != nil {
		return nil, false
	}

	if resp.StatusCode != http.StatusOK {
		return nil, false
	}

	var gacha models.Gacha
	err = json.NewDecoder(resp.Body).Decode(&gacha)
	if err != nil {
		return nil, false
	}

	return &gacha, true
}

func (r GachaRepo) AddGachaToUser(uid models.UUID, gid models.UUID) bool {
	requestData := models.AddGachaToUserData{
		UserID:  uid,
		GachaID: gid,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return false
	}
	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := http.Post(
			addGachaToUserEndpoint,
			"application/json",
			bytes.NewBuffer(jsonData),
		)

		if err != nil {
			return nil, err
		}

		return resp, nil
	})
	defer resp.Body.Close()

	return true
}
