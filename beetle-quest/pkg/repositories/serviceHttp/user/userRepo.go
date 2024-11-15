package repository

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/sony/gobreaker/v2"
)

var (
	createUserEndpoint string = utils.FindEnv("CREATE_USER_ENDPOINT")
	updateUserEndpoint string = utils.FindEnv("UPDATE_USER_ENDPOINT")
	deleteUserEndpoint string = utils.FindEnv("DELETE_USER_ENDPOINT")

	findUserByIDEndpoint       string = utils.FindEnv("FIND_USER_BY_ID_ENDPOINT")
	findUserByUsernameEndpoint string = utils.FindEnv("FIND_USER_BY_USERNAME_ENDPOINT")
	findUserByEmailEndpoint    string = utils.FindEnv("FIND_USER_BY_EMAIL_ENDPOINT")
)

type UserRepo struct {
	client *http.Client
	cb     *gobreaker.CircuitBreaker[*http.Response]
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		client: utils.SetupHTTPSClient(),
		cb:     gobreaker.NewCircuitBreaker[*http.Response](gobreaker.Settings{}),
	}
}

func (r *UserRepo) Create(email, username string, hashedPassword []byte, currency int64) bool {
	requestData := models.CreateUserData{
		Email:          email,
		Username:       username,
		HashedPassword: hashedPassword,
		Currency:       currency,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return false
	}

	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Post(
			createUserEndpoint,
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
	} else {
		return true
	}
}

func (r *UserRepo) Update(user *models.User) bool {
	jsonData, err := json.Marshal(user)
	if err != nil {
		return false
	}

	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Post(
			updateUserEndpoint,
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

func (r *UserRepo) Delete(user *models.User) bool {
	log.Println("[ERROR] Not implemented!")
	return false
}

func (r *UserRepo) FindByID(id models.UUID) (*models.User, bool) {
	requestData := models.FindUserByIDData{
		UserID: id.String(),
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, false
	}

	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Post(
			findUserByIDEndpoint,
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

	var user models.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, false
	}

	return &user, true
}

func (r *UserRepo) FindByUsername(username string) (*models.User, bool) {
	requestData := models.FindUserByUsernameData{
		Username: username,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, false
	}

	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Post(
			findUserByUsernameEndpoint,
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

	var user models.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, false
	}

	return &user, true
}

func (r *UserRepo) FindByEmail(email string) (*models.User, bool) {
	requestData := models.FindUserByEmailData{
		Email: email,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, false
	}

	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Post(
			findUserByEmailEndpoint,
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

	var user models.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, false
	}

	return &user, true
}
