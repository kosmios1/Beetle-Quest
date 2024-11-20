package service

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/repositories"
	"beetle-quest/pkg/utils"
	"strconv"
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

	currency, err := strconv.Atoi(data.Currency)
	if err != nil {
		return false
	}
	if currency >= 0 {
		user.Currency = int64(currency)
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

func (s *AdminService) GetUserAuctionList(userId string) ([]models.Auction, bool) {
	uid, err := utils.ParseUUID(userId)
	if err != nil {
		return []models.Auction{}, false
	}

	if auctions, ok := s.mrepo.GetUserAuctions(uid); ok {
		return auctions, true
	}

	return []models.Auction{}, false
}

// Gacha service functions =================================================

func (s *AdminService) AddGacha(data *models.AdminAddGachaRequest) error {
	price, err := strconv.Atoi(data.Price)
	if err != nil {
		return models.ErrInvalidIntValueAsString
	}

	rarity, err := models.RarityFromString(data.Rarity)
	if err != nil {
		return models.ErrInvalidRarityValue
	}

	gacha := &models.Gacha{
		GachaID:   utils.GenerateUUID(),
		Name:      data.Name,
		Price:     int64(price),
		Rarity:    rarity,
		ImagePath: data.ImagePath,
	}

	if ok := s.grepo.Create(gacha); !ok {
		return models.ErrGachaCreationFailed
	}
	return nil
}

func (s *AdminService) UpdateGacha(gachaId string, data *models.AdminUpdateGachaRequest) bool {
	gid, err := utils.ParseUUID(gachaId)
	if err != nil {
		return false
	}

	gacha, exists := s.grepo.FindByID(gid)
	if !exists {
		return false
	}

	if data.Name != "" {
		gacha.Name = data.Name
	}

	if data.ImagePath != "" {
		gacha.ImagePath = data.ImagePath
	}

	price, err := strconv.Atoi(data.Price)
	if err != nil {
		return false
	}
	if price >= 0 {
		gacha.Price = int64(price)
	}

	// FIXME: If rarity == "Common" it does not change
	rarity, err := models.RarityFromString(data.Rarity)
	if err != nil {
		return false
	}

	gacha.Rarity = rarity

	return s.grepo.Update(gacha)
}

func (s *AdminService) DeleteGacha(gachaId string) bool {
	gid, err := utils.ParseUUID(gachaId)
	if err != nil {
		return false
	}

	gacha, exists := s.grepo.FindByID(gid)
	if !exists {
		return false
	}
	return s.grepo.Delete(gacha)
}

func (s *AdminService) GetAllGachas() ([]models.Gacha, bool) {
	return s.grepo.GetAll()
}

func (s *AdminService) FindGachaByID(gachaId string) (*models.Gacha, bool) {
	gid, err := utils.ParseUUID(gachaId)
	if err != nil {
		return nil, false
	}

	return s.grepo.FindByID(gid)
}

// Market service functions =================================================

func (s *AdminService) GetMarketHistory() ([]models.Transaction, bool) {
	return s.mrepo.GetAllTransactions()
}

func (s *AdminService) GetAllAuctions() ([]models.Auction, bool) {
	return s.mrepo.GetAll()
}

func (s *AdminService) FindAuctionByID(auctionId string) (*models.Auction, bool) {
	aid, err := utils.ParseUUID(auctionId)
	if err != nil {
		return nil, false
	}

	return s.mrepo.FindByID(aid)
}
