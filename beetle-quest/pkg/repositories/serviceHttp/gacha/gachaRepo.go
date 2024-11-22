package repository

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sony/gobreaker/v2"
)

var (
	createEndpoint = os.Getenv("CREATE_GACHA_ENDPOINT")
	updateEndpoint = os.Getenv("UPDATE_GACHA_ENDPOINT")
	deleteEndpoint = os.Getenv("DELETE_GACHA_ENDPOINT")

	getAllEndpoint        = os.Getenv("GET_ALL_GACHA_ENDPOINT")
	findGachaByIDEndpoint = os.Getenv("FIND_GACHA_BY_ID_ENDPOINT")

	getUserGachasEndpoint    = os.Getenv("GET_USER_GACHAS_ENDPOINT")
	removeUserGachasEndpoint = os.Getenv("REMOVE_USER_GACHAS_ENDPOINT")

	addGachaToUserEndpoint      = os.Getenv("ADD_GACHA_TO_USER_ENDPOINT")
	removeGachaFromUserEndpoint = os.Getenv("REMOVE_GACHA_FROM_USER_ENDPOINT")
)

type GachaRepo struct {
	client *http.Client
	cb     *gobreaker.CircuitBreaker[*http.Response]
}

func NewGachaRepo() *GachaRepo {
	return &GachaRepo{
		client: utils.SetupHTTPSClient(),
		// closed: ok, open: service does not respond
		cb: gobreaker.NewCircuitBreaker[*http.Response](gobreaker.Settings{
			MaxRequests: 5,
			Interval:    5 * time.Second,  // When to flush counters int the Closed state
			Timeout:     10 * time.Second, // Time to switch from open -> half-open
			ReadyToTrip: func(counts gobreaker.Counts) bool { // When to switch from closed -> open
				failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
				return (counts.Requests >= 10 && failureRatio >= 0.6) || counts.ConsecutiveFailures > 10
			},
			OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
				log.Printf("[INFO] Circuit breaker changed from %s to %s", from, to)
			},
		}),
	}
}

func (r *GachaRepo) Create(g *models.Gacha) error {
	jsonData, err := json.Marshal(g)
	if err != nil {
		return models.ErrInternalServerError
	}

	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Post(
			createEndpoint,
			"application/json",
			bytes.NewBuffer(jsonData),
		)

		if err != nil {
			return nil, models.ErrInternalServerError
		}

		return resp, nil
	})
	if err != nil {
		return models.ErrInternalServerError
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	} else if resp.StatusCode == http.StatusBadRequest {
		return models.ErrInvalidData
	} else if resp.StatusCode == http.StatusInternalServerError {
		return models.ErrInternalServerError
	} else if resp.StatusCode == http.StatusConflict {
		return models.ErrGachaAlreadyExists
	}

	panic("Unreachable code")
}

func (r *GachaRepo) Update(g *models.Gacha) error {
	jsonData, err := json.Marshal(g)
	if err != nil {
		return models.ErrInternalServerError
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
	if err != nil {
		return models.ErrInternalServerError
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	} else if resp.StatusCode == http.StatusNotFound {
		return models.ErrGachaNotFound
	} else if resp.StatusCode == http.StatusBadRequest {
		return models.ErrInvalidData
	} else if resp.StatusCode == http.StatusInternalServerError {
		return models.ErrInternalServerError
	} else if resp.StatusCode == http.StatusConflict {
		return models.ErrGachaAlreadyExists
	}

	panic("Unreachable code")
}

func (r *GachaRepo) Delete(g *models.Gacha) error {
	jsonData, err := json.Marshal(g)
	if err != nil {
		return models.ErrInternalServerError
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
	if err != nil {
		return models.ErrInternalServerError
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	} else if resp.StatusCode == http.StatusNotFound {
		return models.ErrGachaNotFound
	} else if resp.StatusCode == http.StatusInternalServerError {
		return models.ErrInternalServerError
	} else if resp.StatusCode == http.StatusBadRequest {
		return models.ErrInvalidData
	}

	panic("Unreachable code")
}

func (r *GachaRepo) GetAll() ([]models.Gacha, error) {
	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Get(getAllEndpoint)

		if err != nil {
			return nil, err
		}

		return resp, nil
	})
	if err != nil {
		return nil, models.ErrInternalServerError
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, models.ErrGachaNotFound
	} else if resp.StatusCode == http.StatusInternalServerError {
		return nil, models.ErrInternalServerError
	}

	if resp.StatusCode == http.StatusOK {
		var result models.GetAllGachasDataResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return nil, models.ErrInternalServerError
		}

		return result.GachaList, nil
	}

	panic("Unreachable code")
}

func (r *GachaRepo) FindByID(gid models.UUID) (*models.Gacha, error) {
	requestData := models.FindGachaByIDData{
		GachaID: gid.String(),
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, models.ErrInternalServerError
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
	if err != nil {
		return nil, models.ErrInternalServerError
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, models.ErrGachaNotFound
	} else if resp.StatusCode == http.StatusInternalServerError {
		return nil, models.ErrInternalServerError
	} else if resp.StatusCode == http.StatusBadRequest {
		return nil, models.ErrInvalidData
	}

	if resp.StatusCode == http.StatusOK {
		var gacha models.Gacha
		err = json.NewDecoder(resp.Body).Decode(&gacha)
		if err != nil {
			return nil, models.ErrInternalServerError
		}

		return &gacha, nil
	}
	panic("Unreachable code")
}

func (r *GachaRepo) AddGachaToUser(uid models.UUID, gid models.UUID) error {
	requestData := models.AddGachaToUserData{
		UserID:  uid,
		GachaID: gid,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return models.ErrInternalServerError
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
	if err != nil {
		return models.ErrInternalServerError
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		return models.ErrUserAlreadyHasGacha
	} else if resp.StatusCode == http.StatusInternalServerError {
		return models.ErrInternalServerError
	} else if resp.StatusCode == http.StatusBadRequest {
		return models.ErrInvalidData
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}
	panic("Unreachable code")
}

func (r *GachaRepo) RemoveGachaFromUser(uid models.UUID, gid models.UUID) error {
	requestData := models.RemoveGachaFromUserData{
		UserID:  uid,
		GachaID: gid,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return models.ErrInternalServerError
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
	if err != nil {
		return models.ErrInternalServerError
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return models.ErrRetalationGachaUserNotFound
	} else if resp.StatusCode == http.StatusInternalServerError {
		return models.ErrInternalServerError
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}
	panic("Unreachable code")
}

func (r *GachaRepo) RemoveUserGachas(uid models.UUID) error {
	requestData := models.RemoveUserGachasData{
		UserID: uid,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return models.ErrInternalServerError
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
	if err != nil {
		return models.ErrInternalServerError
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return models.ErrRetalationGachaUserNotFound
	} else if resp.StatusCode == http.StatusInternalServerError {
		return models.ErrInternalServerError
	} else if resp.StatusCode == http.StatusBadRequest {
		return models.ErrInvalidData
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}
	panic("Unreachable code")
}

func (r *GachaRepo) GetUserGachas(uid models.UUID) ([]models.Gacha, error) {
	requestData := models.GetUserGachasData{
		UserID: uid,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, models.ErrInternalServerError
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
	if err != nil {
		return nil, models.ErrInternalServerError
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusInternalServerError {
		return nil, models.ErrInternalServerError
	} else if resp.StatusCode == http.StatusBadRequest {
		return nil, models.ErrInvalidData
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, models.ErrUserNotFound
	}

	if resp.StatusCode == http.StatusOK {
		var result models.GetUserGachasDataResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return nil, models.ErrInternalServerError
		}

		return result.GachaList, nil
	}

	panic("Unreachable code")
}
