package service

import (
	"beetle-quest/internal/auth/repository"
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/repositories"
	"beetle-quest/pkg/utils"
	"context"
	"encoding/hex"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
)

var (
	jwtSecretKey = utils.PanicIfError[[]byte](hex.DecodeString(utils.FindEnv("JWT_SECRET_KEY")))
)

type AuthService struct {
	userRepo   repositories.UserRepo
	oauth2Repo repository.Oauth2Repo
}

func NewAuthService(userRepo repositories.UserRepo, oauth2Repo repository.Oauth2Repo) AuthService {
	return AuthService{
		userRepo:   userRepo,
		oauth2Repo: oauth2Repo,
	}
}

func (s *AuthService) Register(email, username, password string) error {
	if email == "" || username == "" || password == "" {
		return models.ErrInvalidUsernameOrPassOrEmail
	}

	hashedPassword, err := utils.GenerateHashFromPassword([]byte(password))
	if err != nil {
		return models.ErrCheckingPassword
	}

	if ok := s.userRepo.Create(email, username, hashedPassword, 200); !ok {
		return models.ErrUserParametersNotValid
	}

	return nil
}

func (s *AuthService) Login(username, password string) (*models.User, error) {
	if username == "" || password == "" {
		return nil, models.ErrInvalidUsernameOrPass
	}

	user, ok := s.userRepo.FindByUsername(username)
	if !ok {
		return nil, models.ErrInvalidUsernameOrPass
	}

	if err := utils.CompareHashPassword([]byte(password), user.PasswordHash); err != nil {
		return nil, models.ErrInvalidUsernameOrPass
	}

	return user, nil
}

func (s *AuthService) MakeAuthRequest(userId string) (string, error) {
	state, err := utils.GenerateRandomSalt(32)
	if err != nil {
		return "", err
	}
	stateHex := hex.EncodeToString(state)

	url := s.oauth2Repo.AuthCodeURL(stateHex, userId)
	if url == "" {
		return "", err
	}
	return url, nil
}

func (s *AuthService) ExchangeCodeForToken(code string) (*oauth2.Token, jwt.MapClaims, error) {
	token, err := s.oauth2Repo.Exchange(context.Background(), code)
	if err != nil {
		return nil, nil, err
	}

	claims, err := utils.VerifyJWTToken(token.AccessToken, jwtSecretKey)
	if err != nil {
		return nil, nil, err
	}

	return token, claims, nil
}

func (s *AuthService) RevokeToken(token string) bool {
	if err := s.oauth2Repo.RevokeToken(token); err != nil {
		return false
	}

	return true
}

func (s *AuthService) VerifyToken(token string) (jwt.MapClaims, bool) {
	if err := s.oauth2Repo.VerifyToken(token); err != nil {
		return nil, false
	}

	claims, err := utils.VerifyJWTToken(token, jwtSecretKey)
	if err != nil {
		return nil, false
	}

	return claims, true
}
