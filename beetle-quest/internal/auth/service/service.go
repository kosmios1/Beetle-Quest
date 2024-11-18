package service

import (
	"beetle-quest/internal/auth/repository"
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/repositories"
	"beetle-quest/pkg/utils"
	"encoding/hex"

	"github.com/golang-jwt/jwt"
	"github.com/pquerna/otp/totp"
)

var (
	jwtSecretKey = utils.PanicIfError[[]byte](hex.DecodeString(utils.FindEnv("JWT_SECRET_KEY")))
)

type AuthService struct {
	arepo    repositories.AdminRepo
	userRepo repositories.UserRepo

	sesRepo repository.SessionRepo
}

func NewAuthService(userRepo repositories.UserRepo, arepo repositories.AdminRepo) *AuthService {
	return &AuthService{
		arepo:    arepo,
		userRepo: userRepo,
		sesRepo:  *repository.NewSessionRepo(),
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

func (s *AuthService) Login(username, password string) (token *jwt.Token, tokenString string, err error) {
	if username == "" || password == "" {
		return nil, "", models.ErrInvalidUsernameOrPass
	}

	user, ok := s.userRepo.FindByUsername(username)
	if !ok {
		return nil, "", models.ErrInvalidUsernameOrPass
	}

	if err := utils.CompareHashPassword([]byte(password), user.PasswordHash); err != nil {
		return nil, "", models.ErrInvalidUsernameOrPass
	}

	token, tokenString, err = utils.GenerateJWTToken(user.UserID.String(), "user", jwtSecretKey)
	if err != nil {
		return nil, "", err
	}

	if err := s.sesRepo.CreateSession(tokenString); err != nil {
		return nil, "", err
	}

	return token, tokenString, nil
}

func (s *AuthService) RevokeToken(token string) bool {
	if err := s.sesRepo.RevokeToken(token); err != nil {
		return false
	}
	return true
}

func (s *AuthService) VerifyToken(token string) (jwt.MapClaims, bool) {
	if _, err := s.sesRepo.FindToken(token); err != nil {
		return nil, false
	}

	claims, err := utils.VerifyJWTToken(token, jwtSecretKey)
	if err != nil {
		return nil, false
	}

	return claims, true
}

// Admin ==============================================================================================================

func (s *AuthService) AdminLogin(id, password, otp string) (token *jwt.Token, tokenString string, err error) {
	if password == "" {
		return nil, "", models.ErrInvalidAdminIDOrPassOrOTOP
	}

	aid, err := utils.ParseUUID(id)
	if err != nil {
		return nil, "", models.ErrInvalidAdminIDOrPassOrOTOP
	}

	admin, ok := s.arepo.FindByID(aid)
	if !ok {
		return nil, "", models.ErrInvalidAdminIDOrPassOrOTOP
	}

	if err := utils.CompareHashPassword([]byte(password), admin.PasswordHash); err != nil {
		return nil, "", models.ErrInvalidAdminIDOrPassOrOTOP
	}

	if ok := totp.Validate(otp, admin.OtpSecret); !ok {
		return nil, "", models.ErrInvalidAdminIDOrPassOrOTOP
	}

	token, tokenString, err = utils.GenerateJWTToken(admin.AdminId.String(), "admin", jwtSecretKey)
	if err != nil {
		return nil, "", err
	}

	if err := s.sesRepo.CreateSession(tokenString); err != nil {
		return nil, "", err
	}

	return token, tokenString, nil
}
