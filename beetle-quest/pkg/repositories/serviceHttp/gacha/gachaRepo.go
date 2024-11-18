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
	createEndpoint = utils.FindEnv("CREATE_GACHA_ENDPOINT")
	updateEndpoint = utils.FindEnv("UPDATE_GACHA_ENDPOINT")
	deleteEndpoint = utils.FindEnv("DELETE_GACHA_ENDPOINT")

	getAllEndpoint        = utils.FindEnv("GET_ALL_GACHA_ENDPOINT")
	findGachaByIDEndpoint = utils.FindEnv("FIND_GACHA_BY_ID_ENDPOINT")

	getUserGachasEndpoint    = utils.FindEnv("GET_USER_GACHAS_ENDPOINT")
	removeUserGachasEndpoint = utils.FindEnv("REMOVE_USER_GACHAS_ENDPOINT")

	addGachaToUserEndpoint      = utils.FindEnv("ADD_GACHA_TO_USER_ENDPOINT")
	removeGachaFromUserEndpoint = utils.FindEnv("REMOVE_GACHA_FROM_USER_ENDPOINT")
)

type GachaRepo struct {
	client *http.Client
	cb     *gobreaker.CircuitBreaker[*http.Response]
}

func NewGachaRepo() *GachaRepo {
	return &GachaRepo{
		client: utils.SetupHTTPSClient(),
		cb:     gobreaker.NewCircuitBreaker[*http.Response](gobreaker.Settings{}),
	}
}

func (r *GachaRepo) Create(g *models.Gacha) bool {
	jsonData, err := json.Marshal(g)
	if err != nil {
		return false
	}

	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Post(
			createEndpoint,
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
		return false
	}

	if resp.StatusCode != http.StatusOK {
		return false
	}

	return true
}

func (r *GachaRepo) Update(g *models.Gacha) bool {
	jsonData, err := json.Marshal(g)
	if err != nil {
		return false
	}

	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Post(
			updateEndpoint,
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
		return false
	}

	if resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

func (r *GachaRepo) Delete(g *models.Gacha) bool {
	jsonData, err := json.Marshal(g)
	if err != nil {
		return false
	}

	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Post(
			deleteEndpoint,
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
		return false
	}

	if resp.StatusCode != http.StatusOK {
		return false
	}

	return true
}

func (r *GachaRepo) GetAll() ([]models.Gacha, bool) {
	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Get(getAllEndpoint)

		if err != nil {
			return nil, err
		}

		return resp, nil
	})
	defer resp.Body.Close()

	if err != nil {
		return []models.Gacha{}, false
	}

	if resp.StatusCode != http.StatusOK {
		return []models.Gacha{}, false
	}

	var result models.GetAllGachasDataResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return []models.Gacha{}, false
	}
	return result.GachaList, true
}

func (r *GachaRepo) FindByID(gid models.UUID) (*models.Gacha, bool) {
	requestData := models.FindGachaByIDData{
		GachaID: gid.String(),
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, false
	}

	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Post(
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

func (r *GachaRepo) AddGachaToUser(uid models.UUID, gid models.UUID) bool {
	requestData := models.AddGachaToUserData{
		UserID:  uid,
		GachaID: gid,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return false
	}
	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Post(
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

	if err != nil {
		return false
	}

	if resp.StatusCode != http.StatusOK {
		return false
	}

	return true
}

func (r *GachaRepo) RemoveGachaFromUser(uid models.UUID, gid models.UUID) bool {
	requestData := models.RemoveGachaFromUserData{
		UserID:  uid,
		GachaID: gid,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return false
	}
	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Post(
			removeGachaFromUserEndpoint,
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
		return false
	}

	if resp.StatusCode != http.StatusOK {
		return false
	}

	return true
}

func (r *GachaRepo) GetUserGachas(uid models.UUID) ([]models.Gacha, bool) {
	requestData := models.GetUserGachasData{
		UserID: uid,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return []models.Gacha{}, false
	}
	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Post(
			getUserGachasEndpoint,
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
		return []models.Gacha{}, false
	}

	if resp.StatusCode != http.StatusOK {
		return []models.Gacha{}, false
	}

	var result models.GetUserGachasDataResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return []models.Gacha{}, false
	}

	return result.GachaList, true
}

func (r *GachaRepo) RemoveUserGachas(uid models.UUID) bool {
	requestData := models.RemoveUserGachasData{
		UserID: uid,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return false
	}
	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Post(
			removeUserGachasEndpoint,
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
		return false
	}

	if resp.StatusCode != http.StatusOK {
		return false
	}

	return true
}
