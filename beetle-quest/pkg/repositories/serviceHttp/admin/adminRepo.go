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
	findAdminByIDEndpoint string = utils.FindEnv("FIND_ADMIN_BY_ID_ENDPOINT")
)

type AdminRepo struct {
	client *http.Client
	cb     *gobreaker.CircuitBreaker[*http.Response]
}

func NewAdminRepo() *AdminRepo {
	return &AdminRepo{
		client: utils.SetupHTTPSClient(),
		cb:     gobreaker.NewCircuitBreaker[*http.Response](gobreaker.Settings{}),
	}
}

func (r *AdminRepo) FindByID(id models.UUID) (*models.Admin, bool) {
	requestData := models.FindAdminByIDData{
		AdminID: id,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, false
	}

	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Post(
			findAdminByIDEndpoint,
			"application/json",
			bytes.NewBuffer(jsonData),
		)

		if err != nil {
			return nil, err
		}

		return resp, nil
	})

	if err != nil {
		return nil, false
	}

	if resp.StatusCode != http.StatusOK {
		return nil, false
	}

	var admin models.Admin
	err = json.NewDecoder(resp.Body).Decode(&admin)
	if err != nil {
		return nil, false
	}

	return &admin, true
}
