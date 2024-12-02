package repository

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"log"
	"sync"
)

type UserRepo struct {
	mux   sync.RWMutex
	users map[models.UUID]models.User
}

func NewUserRepo() *UserRepo {
	repo := &UserRepo{
		mux:   sync.RWMutex{},
		users: make(map[models.UUID]models.User),
	}

	populateMockRepo(repo)

	return repo
}

func (r *UserRepo) GetAll() ([]models.User, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	var values []models.User
	for _, value := range r.users {
		values = append(values, value)
	}
	return values, nil
}

func (r *UserRepo) Create(user *models.User) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, u := range r.users {
		if user.Email == u.Email || user.Username == u.Username {
			return models.ErrUsernameOrEmailAlreadyExists
		}
	}

	if _, ok := r.users[user.UserID]; ok {
		return models.ErrInternalServerError
	}
	r.users[user.UserID] = *user
	return nil
}

func (r *UserRepo) Update(user *models.User) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	for _, u := range r.users {
		if (u.Username == user.Username || u.Email == user.Email) && u.UserID != user.UserID {
			return models.ErrUsernameOrEmailAlreadyExists
		}
	}

	if _, ok := r.users[user.UserID]; !ok {
		return models.ErrUserNotFound
	}
	r.users[user.UserID] = *user
	return nil
}

func (r *UserRepo) Delete(user *models.User) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	if _, ok := r.users[user.UserID]; !ok {
		return models.ErrUserNotFound
	}
	delete(r.users, user.UserID)
	return nil
}

func (r *UserRepo) FindByID(id models.UUID) (*models.User, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()
	if user, ok := r.users[id]; ok {
		return &user, nil
	}
	return nil, models.ErrUserNotFound
}

func (r *UserRepo) FindByEmail(email string) (*models.User, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()
	for _, user := range r.users {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, models.ErrUserNotFound
}

func (r *UserRepo) FindByUsername(username string) (*models.User, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()
	for _, user := range r.users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, models.ErrUserNotFound
}

// Utils ===============================================================================================================
func populateMockRepo(repo *UserRepo) {
	mockUsers := []models.User{
		{
			UserID:       utils.PanicIfError[models.UUID](utils.ParseUUID("02b84c2f-6b7d-48fd-9850-35610a1d4373")),
			Username:     "Alice",
			Email:        "alice@example.com",
			Currency:     200,
			PasswordHash: utils.PanicIfError(utils.GenerateHashFromPassword([]byte("password"))),
		},
		{
			UserID:       utils.PanicIfError[models.UUID](utils.ParseUUID("744a2f4d-a693-4352-916e-64f4ef94b709")),
			Username:     "Bob",
			Email:        "bob@example.com",
			Currency:     200,
			PasswordHash: utils.PanicIfError(utils.GenerateHashFromPassword([]byte("password"))),
		},
		{
			UserID:       utils.PanicIfError[models.UUID](utils.ParseUUID("7712a483-1202-4225-8dca-0d2d2b60f403")),
			Username:     "Charlie",
			Email:        "charlie@example.com",
			Currency:     200,
			PasswordHash: utils.PanicIfError(utils.GenerateHashFromPassword([]byte("password"))),
		},
	}

	for _, user := range mockUsers {
		if err := repo.Create(&user); err != nil {
			log.Fatal("[FATAL] Could not create user in mock repo!")
		}
	}
}
