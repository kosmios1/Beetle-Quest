package service

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/repositories"
	"beetle-quest/pkg/utils"
)

type UserService struct {
	repositories.UserRepo
}

func (s *UserService) GetUserAccountDetails(userID models.UUID) (*models.User, error) {
	if user, ok := s.UserRepo.FindByID(userID); !ok {
		return nil, models.ErrUserNotFound
	} else {
		return user, nil
	}
}

func (s *UserService) DeleteUserAccount(userID models.UUID, password string) error {
	user, ok := s.UserRepo.FindByID(userID)
	if !ok {
		return models.ErrUserNotFound
	}

	if err := utils.CompareHashPassword([]byte(password), user.PasswordHash); err != nil {
		return models.ErrInvalidPassword
	}

	if ok := s.UserRepo.Delete(user); !ok {
		return models.ErrCouldNotDelete
	}
	return nil
}

func (s *UserService) UpdateUserAccountDetails(userID models.UUID, newEmail, newUsername, oldPassword, newPassword string) error {
	user, ok := s.UserRepo.FindByID(userID)
	if !ok {
		return models.ErrUserNotFound
	}

	if newUsername != "" {
		if _, ok = s.UserRepo.FindByUsername(newUsername); ok {
			return models.ErrUsernameAlreadyExists
		} else {
			user.Username = newUsername
		}
	}

	if newEmail != "" {
		if _, ok = s.UserRepo.FindByEmail(newEmail); ok {
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

	if ok := s.UserRepo.Update(user); !ok {
		return models.ErrCouldNotUpdate
	}
	return nil

}
