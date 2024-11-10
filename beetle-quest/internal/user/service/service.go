package service

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/repositories"
	"beetle-quest/pkg/utils"
)

type UserService struct {
	repo repositories.UserRepo
}

func NewUserService(repo repositories.UserRepo) UserService {
	return UserService{
		repo,
	}
}

func (s *UserService) GetUserAccountDetails(userID models.UUID) (*models.User, error) {
	if user, ok := s.repo.FindByID(userID); !ok {
		return nil, models.ErrUserNotFound
	} else {
		return user, nil
	}
}

func (s *UserService) DeleteUserAccount(userID models.UUID, password string) error {
	user, ok := s.repo.FindByID(userID)
	if !ok {
		return models.ErrUserNotFound
	}

	if err := utils.CompareHashPassword([]byte(password), user.PasswordHash); err != nil {
		return models.ErrInvalidPassword
	}

	if ok := s.repo.Delete(user); !ok {
		return models.ErrCouldNotDelete
	}
	return nil
}

func (s *UserService) UpdateUserAccountDetails(userID models.UUID, newEmail, newUsername, oldPassword, newPassword string) error {
	user, ok := s.repo.FindByID(userID)
	if !ok {
		return models.ErrUserNotFound
	}

	if newUsername != "" {
		if _, ok = s.repo.FindByUsername(newUsername); ok {
			return models.ErrUsernameAlreadyExists
		} else {
			user.Username = newUsername
		}
	}

	if newEmail != "" {
		if _, ok = s.repo.FindByEmail(newEmail); ok {
			return models.ErrEmailAlreadyExists
		} else {
			user.Email = newEmail
		}
	}

	if oldPassword != "" {
		if err := utils.CompareHashPassword([]byte(oldPassword), user.PasswordHash); err != nil {
			return models.ErrInvalidPassword
		}
	}

	if newPassword != "" {
		hashedPassword, err := utils.GenerateHashFromPassword([]byte(newPassword))
		if err != nil {
			return models.ErrCouldNotUseNewPassword
		} else {
			user.PasswordHash = hashedPassword
		}
	}

	if ok := s.repo.Update(user); !ok {
		return models.ErrCouldNotUpdate
	}
	return nil
}

func (s *UserService) Create(email, username string, hashedPassword []byte, currency int64) bool {
	return s.repo.Create(email, username, hashedPassword, currency)
}

func (s *UserService) FindByID(userId string) (*models.User, bool) {
	id, err := utils.ParseUUID(userId)
	if err != nil {
		return nil, false
	}

	user, exits := s.repo.FindByID(id)
	if !exits {
		return nil, false
	}

	return user, true
}

func (s *UserService) FindByUsername(username string) (*models.User, bool) {
	user, exits := s.repo.FindByUsername(username)
	if !exits {
		return nil, false
	}

	return user, true
}
