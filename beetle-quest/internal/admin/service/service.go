package service

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/repositories"
	"beetle-quest/pkg/utils"
)

type AdminService struct {
	mrepo repositories.MarketRepo
	urepo repositories.UserRepo
	grepo repositories.GachaRepo

	arepo repositories.AdminRepo
}

func NewAdminService(arepo repositories.AdminRepo, mrepo repositories.MarketRepo, urepo repositories.UserRepo, grepo repositories.GachaRepo) *AdminService {
	return &AdminService{arepo: arepo, mrepo: mrepo, urepo: urepo, grepo: grepo}
}

func (s *AdminService) FindByID(aid models.UUID) (*models.Admin, bool) {
	return s.arepo.FindByID(aid)
}

// User service functions =================================================
func (s *AdminService) GetAllUsers() ([]models.User, error) {
	return s.urepo.GetAll()
}

func (s *AdminService) FindUserByID(userId string) (*models.UserDetailsTemplate, bool) {
	uid, err := utils.ParseUUID(userId)
	if err != nil {
		return nil, false
	}

	user, exists := s.urepo.FindByID(uid)
	if !exists {
		return nil, false
	}

	gachas, ok := s.grepo.GetUserGachas(uid)
	if !ok {
		gachas = []models.Gacha{}
	}

	transactions, ok := s.mrepo.GetUserTransactionHistory(uid)
	if !ok {
		transactions = []models.Transaction{}
	}

	auctions, ok := s.mrepo.GetUserAuctions(uid)
	if !ok {
		auctions = []models.Auction{}
	}

	return &models.UserDetailsTemplate{
		User:         user,
		Gachas:       gachas,
		Transactions: transactions,
		Auctions:     auctions,
	}, true
}

func (s *AdminService) UpdateUserProfile(userId string, data *models.AdminUpdateUserAccount) bool {
	uid, err := utils.ParseUUID(userId)
	if err != nil {
		return false
	}

	user, exists := s.urepo.FindByID(uid)
	if !exists {
		return false
	}

	if data.Email != "" {
		user.Email = data.Email
	}

	if data.Username != "" {
		user.Username = data.Username
	}

	if data.Currency >= 0 {
		user.Currency = data.Currency
	}

	return s.urepo.Update(user)
}

func (s *AdminService) GetUserTransactionHistory(userId string) ([]models.Transaction, bool) {
	uid, err := utils.ParseUUID(userId)
	if err != nil {
		return []models.Transaction{}, false
	}

	if transactions, ok := s.mrepo.GetUserTransactionHistory(uid); ok {
		return transactions, true
	}

	return []models.Transaction{}, false
}

// Gacha service functions =================================================

// Market service functions =================================================
