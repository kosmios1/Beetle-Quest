package service

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/repositories"
	"beetle-quest/pkg/utils"
)

type UserService struct {
	urepo repositories.UserRepo
	grepo repositories.GachaRepo
	mrepo repositories.MarketRepo
}

func NewUserService(urepo repositories.UserRepo, grepo repositories.GachaRepo, mrepo repositories.MarketRepo) *UserService {
	return &UserService{
		urepo,
		grepo,
		mrepo,
	}
}

func (s *UserService) GetUserAccountDetails(userID string) (*models.User, []models.Gacha, []models.Transaction, error) {
	uid, err := utils.ParseUUID(userID)
	if err != nil {
		return nil, nil, nil, models.ErrInternalServerError

	}
	user, err := s.urepo.FindByID(uid)
	if err != nil {
		return nil, nil, nil, err
	}

	gachas, err := s.grepo.GetUserGachas(uid)
	if err != nil {
		if err == models.ErrInternalServerError {
			return nil, nil, nil, err
		}
		gachas = []models.Gacha{}
	}

	transactions, err := s.mrepo.GetUserTransactionHistory(uid)
	if err != nil {
		if err == models.ErrInternalServerError {
			return nil, nil, nil, err
		}
		transactions = []models.Transaction{}
	}

	return user, gachas, transactions, nil
}

func (s *UserService) DeleteUserAccount(userID string, password string) error {
	uid, err := utils.ParseUUID(userID)
	if err != nil {
		return models.ErrInternalServerError
	}

	user, err := s.urepo.FindByID(uid)
	if err != nil {
		return err
	}

	if err := utils.CompareHashPassword([]byte(password), user.PasswordHash); err != nil {
		return models.ErrInvalidPassword
	}

	if err := s.grepo.RemoveUserGachas(user.UserID); err != nil {
		return err
	}

	if err := s.mrepo.DeleteUserTransactionHistory(user.UserID); err != nil {
		return err
	}

	if err := s.urepo.Delete(user); err != nil {
		return models.ErrInternalServerError
	}
	return nil
}

func (s *UserService) UpdateUserAccountDetails(userID string, newEmail, newUsername, oldPassword, newPassword string) error {
	uid, err := utils.ParseUUID(userID)
	if err != nil {
		return models.ErrInternalServerError
	}

	user, err := s.urepo.FindByID(uid)
	if err != nil {
		return err
	}

	if newUsername != "" {
		if _, err = s.urepo.FindByUsername(newUsername); err == nil {
			return models.ErrUsernameAlreadyExists
		} else {
			user.Username = newUsername
		}
	}

	if newEmail != "" {
		if _, err = s.urepo.FindByEmail(newEmail); err == nil {
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
			return models.ErrInternalServerError
		} else {
			user.PasswordHash = hashedPassword
		}
	}

	if err := s.urepo.Update(user); err != nil {
		return err
	}
	return nil
}

// Internal service functions =================================================================================================

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.urepo.GetAll()
}

func (s *UserService) Create(email, username string, hashedPassword []byte, currency int64) error {
	return s.urepo.Create(email, username, hashedPassword, currency)
}

func (s *UserService) Update(user *models.User) error {
	return s.urepo.Update(user)
}

func (s *UserService) FindByID(uid models.UUID) (*models.User, error) {
	return s.urepo.FindByID(uid)
}

func (s *UserService) FindByUsername(username string) (*models.User, error) {
	return s.urepo.FindByUsername(username)
}
