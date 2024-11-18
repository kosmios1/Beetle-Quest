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

func (s *UserService) GetUserAccountDetails(userID models.UUID) (*models.User, error) {
	if user, ok := s.urepo.FindByID(userID); !ok {
		return nil, models.ErrUserNotFound
	} else {
		return user, nil
	}
}

func (s *UserService) DeleteUserAccount(userID models.UUID, password string) error {
	user, ok := s.urepo.FindByID(userID)
	if !ok {
		return models.ErrUserNotFound
	}

	if err := utils.CompareHashPassword([]byte(password), user.PasswordHash); err != nil {
		return models.ErrInvalidPassword
	}

	if ok := s.grepo.RemoveUserGachas(user.UserID); !ok {
		return models.ErrCouldNotDelete
	}

	if ok := s.mrepo.DeleteUserTransactionHistory(user.UserID); !ok {
		return models.ErrCouldNotDelete
	}

	// NOTE: Even if we don't invalidate the user session, the user is invalid and can't do anything

	if ok := s.urepo.Delete(user); !ok {
		return models.ErrCouldNotDelete
	}
	return nil
}

func (s *UserService) UpdateUserAccountDetails(userID models.UUID, newEmail, newUsername, oldPassword, newPassword string) error {
	user, ok := s.urepo.FindByID(userID)
	if !ok {
		return models.ErrUserNotFound
	}

	if newUsername != "" {
		if _, ok = s.urepo.FindByUsername(newUsername); ok {
			return models.ErrUsernameAlreadyExists
		} else {
			user.Username = newUsername
		}
	}

	if newEmail != "" {
		if _, ok = s.urepo.FindByEmail(newEmail); ok {
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

	if ok := s.urepo.Update(user); !ok {
		return models.ErrCouldNotUpdate
	}
	return nil
}

// Internal service functions =================================================================================================

func (s *UserService) GetAllUsers() []models.User {
	if users, err := s.urepo.GetAll(); err == nil {
		return users
	}
	return []models.User{}
}

func (s *UserService) Create(email, username string, hashedPassword []byte, currency int64) bool {
	return s.urepo.Create(email, username, hashedPassword, currency)
}

func (s *UserService) Update(user *models.User) bool {
	return s.urepo.Update(user)
}

func (s *UserService) FindByID(userId string) (*models.User, bool) {
	id, err := utils.ParseUUID(userId)
	if err != nil {
		return nil, false
	}

	user, exits := s.urepo.FindByID(id)
	if !exits {
		return nil, false
	}

	return user, true
}

func (s *UserService) FindByUsername(username string) (*models.User, bool) {
	user, exits := s.urepo.FindByUsername(username)
	if !exits {
		return nil, false
	}

	return user, true
}

func (s *UserService) GetUserGachaList(userId string) []models.Gacha {
	uid, err := utils.ParseUUID(userId)
	if err != nil {
		return []models.Gacha{}
	}

	gachas, ok := s.grepo.GetUserGachas(uid)
	if !ok {
		return []models.Gacha{}
	}
	return gachas
}

func (s *UserService) GetUserTransactionHistory(userId string) []models.Transaction {
	uid, err := utils.ParseUUID(userId)
	if err != nil {
		return []models.Transaction{}
	}

	transactions, ok := s.mrepo.GetUserTransactionHistory(uid)
	if !ok {
		return []models.Transaction{}
	}
	return transactions
}
