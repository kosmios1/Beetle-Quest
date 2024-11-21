package service

import (
	"beetle-quest/internal/market/repository"
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/repositories"
	"beetle-quest/pkg/utils"
	"log"
	"math/rand/v2"
	"time"
)

type MarketService struct {
	evrepo *repository.EventRepo
	urepo  repositories.UserRepo
	grepo  repositories.GachaRepo
	mrepo  repositories.MarketRepo
}

func NewMarketService(urepo repositories.UserRepo, grepo repositories.GachaRepo, mrepo repositories.MarketRepo) *MarketService {
	evrepo := repository.NewEventRepo()
	srv := &MarketService{
		evrepo: evrepo,
		urepo:  urepo,
		grepo:  grepo,
		mrepo:  mrepo,
	}

	go evrepo.StartSubscriber(srv.closeAuctionCallback, srv.closeAuctionErrorCallback)

	return srv
}

func (s *MarketService) AddBugsCoin(userId string, amount int64) error {
	id, err := utils.ParseUUID(userId)
	if err != nil {
		return models.ErrInternalServerError
	}

	if amount <= 0 {
		return models.ErrInternalServerError
	}

	user, err := s.urepo.FindByID(id)
	if err != nil {
		return err
	}

	t := &models.Transaction{
		TransactionID:   utils.GenerateUUID(),
		TransactionType: models.Deposit,
		UserID:          user.UserID,
		Amount:          amount,
		DateTime:        time.Now(),
		EventType:       models.MarketEv,
		EventID:         models.UUID{},
	}

	if err := s.mrepo.AddTransaction(t); err != nil {
		// NOTE: The client doesn't need to know why the request
		// failed, so we can just return a generic error.
		return models.ErrInternalServerError
	}

	if user.Currency+amount < 0 {
		return models.ErrMaxMoneyExceeded
	}

	user.Currency += amount
	if err := s.urepo.Update(user); err != nil {
		if err != models.ErrUserNotFound {
			// NOTE: Because the client should not know how we are updating the user in the backend
			//  and an error like models.ErrUsernameOrEmailAlreadyExists should not be reported
			// TODO: problem
			return models.ErrInternalServerError
		}
		return err
	}
	return nil
}

func (s *MarketService) RollGacha(userId string) (string, error) {
	uid, err := utils.ParseUUID(userId)
	if err != nil {
		return "", models.ErrInternalServerError
	}

	user, err := s.urepo.FindByID(uid)
	if err != nil {
		return "", err
	}

	if user.Currency < 1000 {
		return "", models.ErrNotEnoughMoneyToRollGacha
	}

	gachas, err := s.grepo.GetAll()
	if err != nil {
		return "", err
	}
	gid := gachas[rand.IntN(len(gachas))].GachaID

	t := &models.Transaction{
		TransactionID:   utils.GenerateUUID(),
		TransactionType: models.Withdraw,
		UserID:          user.UserID,
		Amount:          1000,
		DateTime:        time.Now(),
		EventType:       models.MarketEv,
		EventID:         models.UUID{},
	}

	if err := s.mrepo.AddTransaction(t); err != nil {
		// The client doesn't need to know why the request
		// failed, so we can just return a generic error.
		return "", models.ErrInternalServerError
	}

	user.Currency -= 1000
	if err := s.urepo.Update(user); err != nil {
		if err != models.ErrUserNotFound {
			// Because the client should not know how we are updating the user in the backend
			//  and an error like models.ErrUsernameOrEmailAlreadyExists should not be reported
			return "", models.ErrInternalServerError
		}
		return "", err
	}

	gachas, err = s.grepo.GetUserGachas(uid)
	if err != nil {
		user.Currency += 1000
		err = s.urepo.Update(user)
		if err != models.ErrUserNotFound {
			// Because the client should not know how we are updating the user in the backend
			//  and an error like models.ErrUsernameOrEmailAlreadyExists should not be reported
			// TODO: Problem
			return "", models.ErrInternalServerError
		}
		return "", err
	}

	for _, gacha := range gachas {
		if gacha.GachaID == gid {
			return "Opps you already have this gacha!", nil
		}
	}

	if err := s.grepo.AddGachaToUser(uid, gid); err != nil {
		user.Currency += 1000
		err = s.urepo.Update(user)
		if err != models.ErrUserNotFound {
			// Because the client should not know how we are updating the user in the backend
			//  and an error like models.ErrUsernameOrEmailAlreadyExists should not be reported
			// TODO: Problem
			return "", models.ErrInternalServerError
		}
		return "", err
	}

	return "Gacha successfully obtained, check your inventory!", nil
}

func (s *MarketService) BuyGacha(userId string, gachaId string) error {
	uid, err := utils.ParseUUID(userId)
	if err != nil {
		return models.ErrInternalServerError
	}

	gid, err := utils.ParseUUID(gachaId)
	if err != nil {
		return models.ErrInternalServerError
	}

	userGacha, err := s.grepo.GetUserGachas(uid)
	if err != nil {
		return err
	}

	for _, gacha := range userGacha {
		if gid == gacha.GachaID {
			return models.ErrUserAlreadyHasGacha
		}
	}

	user, err := s.urepo.FindByID(uid)
	if err != nil {
		return err
	}

	gacha, err := s.grepo.FindByID(gid)
	if err != nil {
		return err
	}

	if user.Currency < gacha.Price {
		return models.ErrNotEnoughMoneyToBuyGacha
	}

	t := &models.Transaction{
		TransactionID:   utils.GenerateUUID(),
		TransactionType: models.Withdraw,
		UserID:          user.UserID,
		Amount:          gacha.Price,
		DateTime:        time.Now(),
		EventType:       models.MarketEv,
		EventID:         models.UUID{},
	}

	if err := s.mrepo.AddTransaction(t); err != nil {
		// NOTE: The client doesn't need to know why the request
		// failed, so we can just return a generic error.
		return models.ErrInternalServerError
	}

	user.Currency -= gacha.Price
	if err := s.urepo.Update(user); err != nil {
		if err != models.ErrUserNotFound {
			// Because the client should not know how we are updating the user in the backend
			//  and an error like models.ErrUsernameOrEmailAlreadyExists should not be reported
			return models.ErrInternalServerError
		}
		return err
	}

	if err := s.grepo.AddGachaToUser(uid, gid); err != nil {
		// Compensating transaction
		if err := s.urepo.Update(user); err != nil {
			if err != models.ErrUserNotFound {
				// Because the client should not know how we are updating the user in the backend
				//  and an error like models.ErrUsernameOrEmailAlreadyExists should not be reported
				return models.ErrInternalServerError
			}
			// If the user is not found, we can't do anything
		}
		// No error should be returned to the client
	}
	return nil
}

func (s *MarketService) CreateAuction(userId, gachaId string, endTime time.Time) error {
	uid, err := utils.ParseUUID(userId)
	if err != nil {
		return models.ErrInternalServerError
	}

	gid, err := utils.ParseUUID(gachaId)
	if err != nil {
		return models.ErrInternalServerError
	}

	user, err := s.urepo.FindByID(uid)
	if err != nil {
		return err
	}

	gacha, err := s.grepo.FindByID(gid)
	if err != nil {
		return err
	}

	gachas, err := s.grepo.GetUserGachas(uid)
	if err != nil {
		return err
	}

	var found bool
	for _, g := range gachas {
		if g.GachaID == gid {
			found = true
			break
		}
	}

	if !found {
		return models.ErrUserDoesNotOwnGacha
	}

	auctions, err := s.mrepo.GetUserAuctions(uid)
	if err != nil {
		return err
	}

	for _, a := range auctions {
		if a.GachaID == gid {
			return models.ErrGachaAlreadyAuctioned
		}
	}

	startTime := time.Now()
	if endTime.Before(startTime) || endTime.After(startTime.Add(time.Hour*24)) {
		return models.ErrInvalidEndTime
	}

	auction := &models.Auction{
		AuctionID: utils.GenerateUUID(),
		OwnerID:   user.UserID,
		GachaID:   gacha.GachaID,
		StartTime: time.Now(),
		EndTime:   endTime,
		WinnerID:  models.UUID{},
	}

	if err = s.evrepo.AddEndAuctionEvent(auction); err != nil {
		return models.ErrInternalServerError
	}

	if err := s.mrepo.Create(auction); err != nil {
		return err
	}

	return nil
}

func (s *MarketService) DeleteAuction(userId, auctionId, password string) error {
	uid, err := utils.ParseUUID(userId)
	if err != nil {
		return models.ErrInternalServerError
	}

	aid, err := utils.ParseUUID(auctionId)
	if err != nil {
		return models.ErrInternalServerError
	}

	user, err := s.urepo.FindByID(uid)
	if err != nil {
		return err
	}

	auction, err := s.mrepo.FindByID(aid)
	if err != nil {
		return err
	}

	// Check if the user is the owner of the auction
	if auction.OwnerID != uid {
		return models.ErrUserNotOwnerOfAuction
	}

	// Check if the inserted password is correct
	if err = utils.CompareHashPassword([]byte(password), user.PasswordHash); err != nil {
		return models.ErrInvalidPassword
	}

	if auction.EndTime.Before(time.Now()) {
		return models.ErrAuctionEnded
	}

	// If the auction started more than 1/3 of the total time, it cannot be deleted
	timeNow := time.Now()
	totalAuctionTime := auction.EndTime.Sub(auction.StartTime)
	if timeNow.Sub(auction.StartTime) > totalAuctionTime/3 {
		return models.ErrAuctionIsTooCloseToEnd
	}

	bids, err := s.mrepo.GetBidListOfAuction(aid)
	if err != nil {
		return err
	}

	// If there are bids the auction cannot be deleted
	if len(bids) > 0 {
		return models.ErrAuctionHasBids
	}

	if err := s.mrepo.Delete(auction); err != nil {
		return err
	}

	return nil
}

func (s *MarketService) RetrieveAuctionTemplateList() ([]models.AuctionTemplate, error) {
	auctions, err := s.mrepo.GetAll()
	if err != nil {
		return nil, err
	}

	var data []models.AuctionTemplate = []models.AuctionTemplate{}
	for _, auction := range auctions {
		gacha, err := s.grepo.FindByID(auction.GachaID)
		if err != nil {
			if err != models.ErrGachaNotFound {
				return nil, err
			}

			gacha = &models.Gacha{
				Name:      "Unknown",
				ImagePath: "unknown.png",
			}
		}

		owner, err := s.urepo.FindByID(auction.OwnerID)
		if err != nil {
			if err != models.ErrUserNotFound {
				return nil, err
			}
			owner = &models.User{
				Username: "Unknown",
			}
		}

		data = append(data, models.AuctionTemplate{
			Auction:       auction,
			GachaName:     gacha.Name,
			ImagePath:     gacha.ImagePath,
			OwnerUsername: owner.Username,
		})
	}

	return data, nil
}

func (s *MarketService) GetAuctionDetails(auctionId string) (*models.Auction, []models.Bid, error) {
	auction, err := s.FindByID(auctionId)
	if err != nil {
		return nil, nil, err
	}

	bids, err := s.GetBidListOfAuctionID(auctionId)
	if err != nil {
		return nil, nil, err
	}

	return auction, bids, nil
}

func (s *MarketService) FindByID(auctionId string) (*models.Auction, error) {
	aid, err := utils.ParseUUID(auctionId)
	if err != nil {
		return nil, models.ErrInternalServerError
	}

	auction, err := s.mrepo.FindByID(aid)
	if err != nil {
		return nil, err
	}
	return auction, nil
}

func (s *MarketService) GetBidListOfAuctionID(auctionId string) ([]models.Bid, error) {
	aid, err := utils.ParseUUID(auctionId)
	if err != nil {
		return nil, models.ErrInternalServerError
	}

	bids, err := s.mrepo.GetBidListOfAuction(aid)
	if err != nil {
		return nil, err
	}
	return bids, nil
}

func (s *MarketService) MakeBid(userId, auctionId string, bidAmount int64) error {
	uid, err := utils.ParseUUID(userId)
	if err != nil {
		return models.ErrInternalServerError
	}

	aid, err := utils.ParseUUID(auctionId)
	if err != nil {
		return models.ErrInternalServerError
	}

	auction, err := s.mrepo.FindByID(aid)
	if err != nil {
		return err
	}

	if auction.OwnerID == uid {
		return models.ErrOwnerCannotBid
	}

	user, err := s.urepo.FindByID(uid)
	if err != nil {
		return err
	}

	if user.Currency < bidAmount {
		return models.ErrNotEnoughMoneyToBid
	}

	bids, err := s.mrepo.GetBidListOfAuction(aid)
	if err != nil {
		return err
	}

	var maxBid int64 = 0
	for _, bid := range bids {
		if bid.AmountSpend > maxBid {
			maxBid = bid.AmountSpend
		}
	}

	if bidAmount <= maxBid {
		return models.ErrBidAmountNotEnough
	}

	if auction.EndTime.Before(time.Now()) {
		return models.ErrAuctionEnded
	}

	bid := &models.Bid{
		BidID:       utils.GenerateUUID(),
		UserID:      uid,
		AuctionID:   aid,
		AmountSpend: int64(bidAmount),
		TimeStamp:   time.Now(),
	}

	user.Currency -= bidAmount
	if err := s.urepo.Update(user); err != nil {
		if err != models.ErrUserNotFound {
			return models.ErrInternalServerError
		}
		return err
	}

	if err := s.mrepo.BidToAuction(bid); err != nil {
		user.Currency += bidAmount
		if err := s.urepo.Update(user); err != nil {
			if err != models.ErrUserNotFound {
				return models.ErrInternalServerError
			}
			return err
		}
		if err == models.ErrCouldNotBidToAuction {
			return models.ErrInternalServerError
		}
		return err
	}
	return nil
}

// Timed events callbacks ================================================

func (s *MarketService) closeAuctionCallback(aid models.UUID) {
	auction, err := s.mrepo.FindByID(aid)
	if err != nil {
		s.closeAuctionErrorCallback(models.ErrAuctionNotFound)
		return
	}

	bids, err := s.mrepo.GetBidListOfAuction(aid)
	if err != nil {
		s.closeAuctionErrorCallback(models.ErrCouldNotRetrieveAuctionBids)
		return
	}

	if len(bids) == 0 {
		if err = s.mrepo.Delete(auction); err != nil {
			s.closeAuctionErrorCallback(models.ErrCouldNotDeleteAuction)
		}
		return
	}

	var maxBid int64 = 0
	var maxBidder models.UUID
	var totalUserBiddings map[models.UUID]int64 = make(map[models.UUID]int64)
	{ // Give back money to the losers
		for _, bid := range bids {
			if bid.AmountSpend > maxBid {
				maxBid = bid.AmountSpend
				maxBidder = bid.UserID
			}
			totalUserBiddings[bid.UserID] += bid.AmountSpend
		}

		for uid, totBidAmount := range totalUserBiddings {
			if uid == maxBidder {
				continue
			}

			user, err := s.urepo.FindByID(uid)
			if err != nil {
				s.closeAuctionErrorCallback(models.ErrUserNotFound)
				continue
			}

			user.Currency += totBidAmount // NOTE: If we go over math.MaxInt64
			if err := s.urepo.Update(user); err != nil {
				s.closeAuctionErrorCallback(models.ErrCouldNotUpdate)
			}
		}
	}

	gacha, err := s.grepo.FindByID(auction.GachaID)
	if err != nil {
		// NOTE: If we do not find the gacha, we still have informations to give back the money
		// to the winner
		user, err := s.urepo.FindByID(maxBidder)
		if err != nil {
			s.closeAuctionErrorCallback(models.ErrUserNotFound)
			return
		}

		user.Currency += totalUserBiddings[maxBidder]
		if err := s.urepo.Update(user); err != nil {
			s.closeAuctionErrorCallback(models.ErrCouldNotUpdate)
		}
	}

	{ // Winner actions
		auction.WinnerID = maxBidder
		if err := s.mrepo.Update(auction); err != nil {
			// NOTE: If we do not find the auction, we still have informations to give gacha to the winner
			s.closeAuctionErrorCallback(models.ErrCouldNotUpdateAuction)
		}

		user, err := s.urepo.FindByID(maxBidder)
		if err != nil {
			// NOTE: If the winner does not exist, we still have informations to give back the money to the owner
			s.closeAuctionErrorCallback(models.ErrUserNotFound)
		} else {
			t := &models.Transaction{
				TransactionID:   utils.GenerateUUID(),
				TransactionType: models.Withdraw,
				UserID:          user.UserID,
				Amount:          maxBid,
				DateTime:        time.Now(),
				EventType:       models.AuctionEv,
				EventID:         auction.AuctionID,
			}

			if err := s.mrepo.AddTransaction(t); err != nil {
				s.closeAuctionErrorCallback(models.ErrCouldNotAddTransaction)
				return
			}

			if err := s.grepo.AddGachaToUser(user.UserID, gacha.GachaID); err != nil {
				s.closeAuctionErrorCallback(models.ErrCouldNotAddGachaToUser)
				return
			}
		}
	}

	{ // Auction owner actions
		user, err := s.urepo.FindByID(auction.OwnerID)
		if err != nil {
			s.closeAuctionErrorCallback(models.ErrUserNotFound)
			return
		}

		if err := s.grepo.RemoveGachaFromUser(user.UserID, gacha.GachaID); err != nil {
			s.closeAuctionErrorCallback(models.ErrCouldNotAddGachaToUser)
			return
		}

		t := &models.Transaction{
			TransactionID:   utils.GenerateUUID(),
			TransactionType: models.Deposit,
			UserID:          user.UserID,
			Amount:          maxBid,
			DateTime:        time.Now(),
			EventType:       models.AuctionEv,
			EventID:         auction.AuctionID,
		}

		if err := s.mrepo.AddTransaction(t); err != nil {
			s.closeAuctionErrorCallback(models.ErrCouldNotAddTransaction)
		}

		user.Currency += maxBid
		if err := s.urepo.Update(user); err != nil {
			s.closeAuctionErrorCallback(models.ErrCouldNotUpdate)
			// TODO: What should we do here ?
			return
		}
	}
}

func (s *MarketService) closeAuctionErrorCallback(err error) {
	log.Printf("Error closing auction: %v", err)
}

// Internal functions ====================================================

func (s *MarketService) GetAuctionList() ([]models.Auction, error) {
	return s.mrepo.GetAll()
}

func (s *MarketService) GetAuctionListOfUser(uid models.UUID) ([]models.Auction, error) {
	return s.mrepo.GetUserAuctions(uid)
}

func (s *MarketService) GetAllTransactions() ([]models.Transaction, error) {
	return s.mrepo.GetAllTransactions()
}

func (s *MarketService) GetUserTransactionHistory(uid models.UUID) ([]models.Transaction, error) {
	return s.mrepo.GetUserTransactionHistory(uid)
}

func (s *MarketService) DeleteUserTransactionHistory(uid models.UUID) error {
	return s.mrepo.DeleteUserTransactionHistory(uid)
}
