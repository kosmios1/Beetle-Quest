package service

import (
	"beetle-quest/internal/auth/repository"
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/repositories"
	"beetle-quest/pkg/utils"
)

type AuthService struct {
	UserRepo repositories.UserRepo
	repository.Oauth2Repo
}

func (s *AuthService) Register(email, username, password string) error {
	if email == "" || username == "" || password == "" {
		return models.ErrInvalidUsernameOrPassOrEmail
	}

	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// if err != nil {
	// 	return models.ErrCheckingPassword
	// }

	hashedPassword, err := utils.GenerateHashFromPassword([]byte(password))
	if err != nil {
		return models.ErrCheckingPassword
	}

	if ok := s.UserRepo.Create(email, username, hashedPassword, 200); !ok {
		return models.ErrUserParametersNotValid
	}

	return nil
}

func (s *AuthService) Login(username, password string) (*models.User, error) {
	if username == "" || password == "" {
		return nil, models.ErrInvalidUsernameOrPass
	}

	user, ok := s.UserRepo.FindByUsername(username)
	if !ok {
		return nil, models.ErrInvalidUsernameOrPass
	}

	// if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
	// 	return nil, models.ErrInvalidUsernameOrPass
	// }

	if err := utils.CompareHashPassword([]byte(password), user.PasswordHash); err != nil {
		return nil, models.ErrInvalidUsernameOrPass
	}

	return user, nil
}
