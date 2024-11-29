package repository

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
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

func (r *UserRepo) Create(email, username string, hashedPassword []byte, currency int64) error {
	r.mux.RLock()
	uuid := utils.GenerateUUID()

	for _, user := range r.users {
		if user.Email == email || user.Username == username {
			return models.ErrUsernameOrEmailAlreadyExists
		}
	}

	if _, ok := r.users[uuid]; ok {
		return models.ErrInternalServerError
	}
	r.mux.RUnlock()

	r.mux.Lock()
	defer r.mux.Unlock()
	r.users[uuid] = models.User{
		UserID:       uuid,
		Username:     username,
		Email:        email,
		Currency:     currency,
		PasswordHash: hashedPassword,
	}
	return nil
}

func (r *UserRepo) Update(user *models.User) error {
	r.mux.RLock()
	for _, u := range r.users {
		if u.Username == user.Username || u.Email == user.Email {
			return models.ErrUsernameOrEmailAlreadyExists
		}
	}
	r.mux.RUnlock()

	r.mux.Lock()
	defer r.mux.Unlock()
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
			Username:     "Alice",
			Email:        "alice@example.com",
			Currency:     200,
			PasswordHash: utils.PanicIfError(utils.GenerateHashFromPassword([]byte("password"))),
		},
		{
			Username:     "Bob",
			Email:        "bob@example.com",
			Currency:     200,
			PasswordHash: utils.PanicIfError(utils.GenerateHashFromPassword([]byte("password"))),
		},
		{
			Username:     "Charlie",
			Email:        "charlie@example.com",
			Currency:     200,
			PasswordHash: utils.PanicIfError(utils.GenerateHashFromPassword([]byte("password"))),
		},
	}

	for _, user := range mockUsers {
		repo.Create(user.Email, user.Email, user.PasswordHash, user.Currency)
	}
}
