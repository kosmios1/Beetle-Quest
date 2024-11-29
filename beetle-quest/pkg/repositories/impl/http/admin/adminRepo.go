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
	findAdminByIDEndpoint string = os.Getenv("FIND_ADMIN_BY_ID_ENDPOINT")
)

type AdminRepo struct {
	client *http.Client
	cb     *gobreaker.CircuitBreaker[*http.Response]
}

func NewAdminRepo() *AdminRepo {
	return &AdminRepo{
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
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		return nil, models.ErrInternalServerError
	case http.StatusNotFound:
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

	log.Panicf("Unreachable code, status code received: %d", resp.StatusCode)
	return nil, models.ErrInternalServerError
}
