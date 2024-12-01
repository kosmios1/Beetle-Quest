package service

import (
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

	sesRepo repositories.SessionRepo
}

func NewAuthService(userRepo repositories.UserRepo, arepo repositories.AdminRepo, sessionRepo repositories.SessionRepo) *AuthService {
	return &AuthService{
		arepo:    arepo,
		userRepo: userRepo,
		sesRepo:  sessionRepo,
	}
}

func (s *AuthService) Register(email, username, password string) error {
	if email == "" || username == "" || password == "" {
		return models.ErrInvalidUsernameOrPassOrEmail
	}

	hashedPassword, err := utils.GenerateHashFromPassword([]byte(password))
	if err != nil {
		return models.ErrInternalServerError
	}

	user := models.User{
		UserID:       utils.GenerateUUID(),
		Email:        email,
		Username:     username,
		PasswordHash: hashedPassword,
		Currency:     200,
	}

	if err := s.userRepo.Create(&user); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(username, password string) (token *jwt.Token, tokenString string, err error) {
	if username == "" || password == "" {
		return nil, "", models.ErrInvalidUsernameOrPass
	}

	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, "", err
	}

	if err := utils.CompareHashPassword([]byte(password), user.PasswordHash); err != nil {
		return nil, "", models.ErrInvalidPassword
	}

	token, tokenString, err = utils.GenerateJWTToken(user.UserID.String(), false, jwtSecretKey)
	if err != nil {
		return nil, "", models.ErrInternalServerError
	}

	if err := s.sesRepo.CreateSession(tokenString); err != nil {
		return nil, "", models.ErrInternalServerError
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

	admin, err := s.arepo.FindByID(aid)
	if err != nil {
		// NOTE: Even if the server has the capability to know that the admin does not exist, it should not return
		// this information to the client.
		return nil, "", models.ErrInvalidAdminIDOrPassOrOTOP
	}

	if err := utils.CompareHashPassword([]byte(password), admin.PasswordHash); err != nil {
		return nil, "", models.ErrInvalidAdminIDOrPassOrOTOP
	}

	if ok := totp.Validate(otp, admin.OtpSecret); !ok {
		return nil, "", models.ErrInvalidAdminIDOrPassOrOTOP
	}

	token, tokenString, err = utils.GenerateJWTToken(admin.AdminId.String(), true, jwtSecretKey)
	if err != nil {
		return nil, "", models.ErrInternalServerError
	}

	if err := s.sesRepo.CreateSession(tokenString); err != nil {
		return nil, "", models.ErrInternalServerError
	}

	return token, tokenString, nil
}
