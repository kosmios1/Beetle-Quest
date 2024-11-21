package repository

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/sony/gobreaker/v2"
)

var (
	getAllUsersEndpoint string = os.Getenv("GET_ALL_USERS_ENDPOINT")
	createUserEndpoint  string = os.Getenv("CREATE_USER_ENDPOINT")
	updateUserEndpoint  string = os.Getenv("UPDATE_USER_ENDPOINT")
	deleteUserEndpoint  string = os.Getenv("DELETE_USER_ENDPOINT")

	findUserByIDEndpoint       string = os.Getenv("FIND_USER_BY_ID_ENDPOINT")
	findUserByUsernameEndpoint string = os.Getenv("FIND_USER_BY_USERNAME_ENDPOINT")
	findUserByEmailEndpoint    string = os.Getenv("FIND_USER_BY_EMAIL_ENDPOINT")
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

func (r *UserRepo) GetAll() ([]models.User, error) {
	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Get(getAllUsersEndpoint)
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
		return nil, models.ErrUserNotFound
	case http.StatusBadRequest:
		return nil, models.ErrInvalidData
	}

	if resp.StatusCode == http.StatusOK {
		var data models.GetAllUsersDataResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return nil, models.ErrInternalServerError
		}

		return data.UserList, nil
	}
	panic("unreachable code")
}

func (r *UserRepo) Create(email, username string, hashedPassword []byte, currency int64) error {
	requestData := models.CreateUserData{
		Email:          email,
		Username:       username,
		HashedPassword: hashedPassword,
		Currency:       currency,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return models.ErrInternalServerError
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
	if err != nil {
		return models.ErrInternalServerError
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		return models.ErrInternalServerError
	case http.StatusConflict:
		return models.ErrUsernameOrEmailAlreadyExists
	case http.StatusBadRequest:
		return models.ErrInvalidData
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	panic("unreachable code")
}

func (r *UserRepo) Update(user *models.User) error {
	jsonData, err := json.Marshal(user)
	if err != nil {
		return models.ErrInternalServerError
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
	if err != nil {
		return models.ErrInternalServerError
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		return models.ErrInternalServerError
	case http.StatusNotFound:
		return models.ErrUserNotFound
	case http.StatusConflict:
		return models.ErrUsernameOrEmailAlreadyExists
	case http.StatusBadRequest:
		return models.ErrInvalidData
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	panic("unreachable code")
}

func (r *UserRepo) FindByID(id models.UUID) (*models.User, error) {
	requestData := models.FindUserByIDData{
		UserID: id,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, models.ErrInternalServerError
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
	if err != nil {
		return nil, models.ErrInternalServerError
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		return nil, models.ErrInternalServerError
	case http.StatusNotFound:
		return nil, models.ErrUserNotFound
	case http.StatusBadRequest:
		return nil, models.ErrInvalidData
	}

	if resp.StatusCode == http.StatusOK {
		var user models.User
		err = json.NewDecoder(resp.Body).Decode(&user)
		if err != nil {
			return nil, models.ErrInternalServerError
		}
		return &user, nil
	}

	panic("unreachable code")
}

func (r *UserRepo) FindByUsername(username string) (*models.User, error) {
	requestData := models.FindUserByUsernameData{
		Username: username,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, models.ErrInternalServerError
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
	if err != nil {
		return nil, models.ErrInternalServerError
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		return nil, models.ErrInternalServerError
	case http.StatusNotFound:
		return nil, models.ErrUserNotFound
	case http.StatusBadRequest:
		return nil, models.ErrInvalidData
	}

	if resp.StatusCode == http.StatusOK {
		var user models.User
		err = json.NewDecoder(resp.Body).Decode(&user)
		if err != nil {
			return nil, models.ErrInternalServerError
		}
		return &user, nil
	}
	panic("unreachable code")
}

func (r *UserRepo) FindByEmail(email string) (*models.User, error) {
	log.Fatal("[ERROR] Not implemented!")
	return nil, errors.New("not implemented")
}

func (r *UserRepo) Delete(user *models.User) error {
	log.Fatal("[ERROR] Not implemented!")
	return errors.New("not implemented")
}
