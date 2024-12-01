package service

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/repositories"
	"beetle-quest/pkg/utils"
	"strconv"
	"time"
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

func (s *AdminService) FindByID(aid models.UUID) (*models.Admin, error) {
	return s.arepo.FindByID(aid)
}

// User service functions =================================================
func (s *AdminService) GetAllUsers() ([]models.User, error) {
	return s.urepo.GetAll()
}

func (s *AdminService) FindUserByID(userId string) (*models.UserDetailsTemplate, error) {
	uid, err := utils.ParseUUID(userId)
	if err != nil {
		return nil, models.ErrInternalServerError
	}

	user, err := s.urepo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	gachas, err := s.grepo.GetUserGachas(uid)
	if err != nil {
		if err == models.ErrInternalServerError {
			return nil, err
		}
		gachas = []models.Gacha{}
	}

	transactions, err := s.mrepo.GetUserTransactionHistory(uid)
	if err != nil {
		transactions = []models.Transaction{}
	}

	auctions, err := s.mrepo.GetUserAuctions(uid)
	if err != nil {
		if err == models.ErrInternalServerError {
			return nil, err
		}
		auctions = []models.Auction{}
	}

	return &models.UserDetailsTemplate{
		User:         user,
		Gachas:       gachas,
		Transactions: transactions,
		Auctions:     auctions,
	}, nil
}

func (s *AdminService) UpdateUserProfile(userId string, data *models.AdminUpdateUserAccount) error {
	uid, err := utils.ParseUUID(userId)
	if err != nil {
		return models.ErrInternalServerError
	}

	user, err := s.urepo.FindByID(uid)
	if err != nil {
		return err
	}

	if data.Email != "" {
		user.Email = data.Email
	}

	if data.Username != "" {
		user.Username = data.Username
	}

	currency, err := strconv.Atoi(data.Currency)
	if err != nil {
		return models.ErrInternalServerError
	}
	if currency >= 0 {
		user.Currency = int64(currency)
	}

	return s.urepo.Update(user)
}

func (s *AdminService) GetUserTransactionHistory(userId string) ([]models.Transaction, error) {
	uid, err := utils.ParseUUID(userId)
	if err != nil {
		return nil, models.ErrInternalServerError
	}

	if transactions, err := s.mrepo.GetUserTransactionHistory(uid); err == nil {
		return transactions, nil
	}

	return nil, err
}

func (s *AdminService) GetUserAuctionList(userId string) ([]models.Auction, error) {
	uid, err := utils.ParseUUID(userId)
	if err != nil {
		return nil, models.ErrInternalServerError
	}

	if auctions, err := s.mrepo.GetUserAuctions(uid); err == nil {
		return auctions, nil
	}

	return nil, err
}

// Gacha service functions =================================================

func (s *AdminService) AddGacha(data *models.AdminAddGachaRequest) error {
	price, err := strconv.Atoi(data.Price)
	if err != nil {
		return models.ErrInternalServerError
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

	return s.grepo.Create(gacha)
}

func (s *AdminService) UpdateGacha(gachaId string, data *models.AdminUpdateGachaRequest) error {
	gid, err := utils.ParseUUID(gachaId)
	if err != nil {
		return models.ErrInternalServerError
	}

	gacha, err := s.grepo.FindByID(gid)
	if err != nil {
		return err
	}

	if data.Name != "" {
		gacha.Name = data.Name
	}

	if data.ImagePath != "" {
		gacha.ImagePath = data.ImagePath
	}

	price, err := strconv.Atoi(data.Price)
	if err != nil {
		return models.ErrInternalServerError
	}
	if price >= 0 {
		gacha.Price = int64(price)
	}

	// FIXME: If rarity == "Common" it does not change
	rarity, err := models.RarityFromString(data.Rarity)
	if err != nil {
		return models.ErrInternalServerError
	}
	gacha.Rarity = rarity

	return s.grepo.Update(gacha)
}

func (s *AdminService) DeleteGacha(gachaId string) error {
	gid, err := utils.ParseUUID(gachaId)
	if err != nil {
		return models.ErrInternalServerError
	}

	gacha, error := s.grepo.FindByID(gid)
	if error != nil {
		return err
	}
	return s.grepo.Delete(gacha)
}

func (s *AdminService) GetAllGachas() ([]models.Gacha, error) {
	return s.grepo.GetAll()
}

func (s *AdminService) FindGachaByID(gachaId string) (*models.Gacha, error) {
	gid, err := utils.ParseUUID(gachaId)
	if err != nil {
		return nil, models.ErrInternalServerError
	}

	return s.grepo.FindByID(gid)
}

// Market service functions =================================================

func (s *AdminService) GetMarketHistory() ([]models.Transaction, error) {
	return s.mrepo.GetAllTransactions()
}

func (s *AdminService) GetAllAuctions() ([]models.Auction, error) {
	return s.mrepo.GetAll()
}

func (s *AdminService) FindAuctionByID(auctionId string) (*models.Auction, error) {
	aid, err := utils.ParseUUID(auctionId)
	if err != nil {
		return nil, models.ErrInternalServerError
	}
	return s.mrepo.FindByID(aid)
}

func (s *AdminService) UpdateAuction(auctionId string, gachaId string) error {
	aid, err := utils.ParseUUID(auctionId)
	if err != nil {
		return models.ErrInternalServerError
	}

	gid, err := utils.ParseUUID(gachaId)
	if err != nil {
		return models.ErrInternalServerError
	}

	auction, err := s.mrepo.FindByID(aid)
	if err != nil {
		return err
	}

	gacha, err := s.grepo.FindByID(gid)
	if err != nil {
		return err
	}

	if auction.EndTime.Before(time.Now()) {
		return models.ErrAuctionEnded
	}

	auction.GachaID = gacha.GachaID
	return s.mrepo.Update(auction)
}
