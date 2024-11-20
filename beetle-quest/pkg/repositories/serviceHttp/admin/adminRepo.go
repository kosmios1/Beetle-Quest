package repository

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/sony/gobreaker/v2"
)

var (
	findAdminByIDEndpoint string = os.Getenv("FIND_ADMIN_BY_ID_ENDPOINT")
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

func (r *AdminRepo) FindByID(id models.UUID) (*models.Admin, error) {
	requestData := models.FindAdminByIDData{
		AdminID: id,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, models.ErrInternalServerError
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
		return nil, models.ErrInternalServerError
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return nil, models.ErrInternalServerError
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, models.ErrAdminNotFound
	}

	if resp.StatusCode == http.StatusOK {
		var admin models.Admin
		err = json.NewDecoder(resp.Body).Decode(&admin)
		if err != nil {
			return nil, models.ErrInternalServerError
		}

		return &admin, nil
	}

	panic("unreachable code")
}
