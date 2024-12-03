package repository

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"log"
	"sync"
	"time"
)

type MarketRepo struct {
	mux sync.RWMutex

	auctions    map[models.UUID]models.Auction
	auctionBids map[models.UUID][]models.Bid

	transactions map[models.UUID][]models.Transaction
}

func NewMarketRepo() *MarketRepo {
	repo := &MarketRepo{
		// NOTE: we should use more than one mutex but for the sake of simplicity,
		// being a mock repo, we will use only one
		mux: sync.RWMutex{},

		// map[auctionId] -> auction
		auctions: make(map[models.UUID]models.Auction),

		// map[auctionId] -> []bids
		auctionBids: make(map[models.UUID][]models.Bid),

		// map[userId] -> []transactions
		transactions: make(map[models.UUID][]models.Transaction),
	}

	populateMockRepo(repo)

	return repo
}

func (r *MarketRepo) Create(auction *models.Auction) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	if _, ok := r.auctions[auction.AuctionID]; ok {
		return models.ErrAuctionAltreadyExists
	}

	r.auctions[auction.AuctionID] = *auction
	r.auctionBids[auction.AuctionID] = make([]models.Bid, 0)
	return nil
}

func (r *MarketRepo) Update(auction *models.Auction) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	if _, ok := r.auctions[auction.AuctionID]; !ok {
		return models.ErrAuctionNotFound
	}
	r.auctions[auction.AuctionID] = *auction
	return nil
}

func (r *MarketRepo) Delete(auction *models.Auction) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	if _, ok := r.auctions[auction.AuctionID]; !ok {
		return models.ErrUserNotFound
	}
	delete(r.auctions, auction.AuctionID)
	if _, ok := r.auctionBids[auction.AuctionID]; ok {
		delete(r.auctionBids, auction.AuctionID)
	}
	return nil
}

func (r *MarketRepo) GetAll() ([]models.Auction, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()
	var values []models.Auction
	for _, v := range r.auctions {
		values = append(values, v)
	}
	return values, nil
}

func (r *MarketRepo) GetUserAuctions(uid models.UUID) ([]models.Auction, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()
	var values []models.Auction
	for _, v := range r.auctions {
		if v.OwnerID == uid {
			values = append(values, v)
		}
	}
	return values, nil
}

func (r *MarketRepo) FindByID(aid models.UUID) (*models.Auction, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()
	if auction, ok := r.auctions[aid]; ok {
		return &auction, nil
	}

	return nil, models.ErrAuctionNotFound
}

func (r *MarketRepo) GetBidListOfAuction(aid models.UUID) ([]models.Bid, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()
	if _, ok := r.auctionBids[aid]; !ok {
		return nil, models.ErrBidsNotFound
	}

	var bids []models.Bid
	copy(bids, r.auctionBids[aid])
	return bids, nil
}

func (r *MarketRepo) BidToAuction(bid *models.Bid) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	if _, ok := r.auctionBids[bid.AuctionID]; !ok {
		r.auctionBids[bid.AuctionID] = make([]models.Bid, 0)
	}
	r.auctionBids[bid.AuctionID] = append(r.auctionBids[bid.AuctionID], *bid)
	return nil
}

func (r *MarketRepo) GetAllTransactions() ([]models.Transaction, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()
	var outTransactions []models.Transaction
	for _, t := range r.transactions {
		outTransactions = append(outTransactions, t...)
	}
	return outTransactions, nil
}

func (r *MarketRepo) DeleteUserTransactionHistory(uid models.UUID) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	if _, ok := r.transactions[uid]; ok {
		delete(r.transactions, uid)
	}
	return nil
}

func (r *MarketRepo) GetUserTransactionHistory(uid models.UUID) ([]models.Transaction, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()
	if transactions, ok := r.transactions[uid]; ok {
		return transactions, nil
	}
	return []models.Transaction{}, nil
}

func (r *MarketRepo) AddTransaction(transaction *models.Transaction) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	if _, ok := r.transactions[transaction.UserID]; ok {
		for _, t := range r.transactions[transaction.UserID] {
			if t.TransactionID == transaction.TransactionID {
				return models.ErrCouldNotAddTransaction
			}
		}
	} else {
		r.transactions[transaction.UserID] = make([]models.Transaction, 0)
	}
	r.transactions[transaction.UserID] = append(r.transactions[transaction.UserID], *transaction)
	return nil
}

// Utils ================================================================================================================

func populateMockRepo(repo *MarketRepo) {
	mockAuctions := []models.Auction{
		{
			AuctionID: utils.PanicIfError[models.UUID](utils.ParseUUID("77934f96-38eb-4252-a426-7302ac26d58a")),
			OwnerID:   utils.PanicIfError[models.UUID](utils.ParseUUID("744a2f4d-a693-4352-916e-64f4ef94b709")),
			GachaID:   utils.PanicIfError[models.UUID](utils.ParseUUID("e455113c-655c-478d-bd24-b2a59c11e1f3")),
			StartTime: time.Now(),
			EndTime:   time.Now().Add(1 * time.Hour),
			WinnerID:  utils.PanicIfError[models.UUID](utils.ParseUUID("00000000-0000-0000-0000-000000000000")),
		},
	}

	for _, auction := range mockAuctions {
		if err := repo.Create(&auction); err != nil {
			log.Fatal("[FATAL] Could not create auction in mock repo!")
		}
	}

	mockTransactions := []models.Transaction{
	  {
			TransactionID:   utils.PanicIfError[models.UUID](utils.ParseUUID("a1c6f276-a8e0-421f-8b35-4f50e145922f")),
			TransactionType: models.Deposit,
			UserID:          utils.PanicIfError[models.UUID](utils.ParseUUID("744a2f4d-a693-4352-916e-64f4ef94b709")),
			Amount:          10,
			DateTime:        time.Now(),
			EventType:       models.GameEv,
			EventID:         utils.PanicIfError[models.UUID](utils.ParseUUID("15e631fb-b214-4c57-aca8-075420debacb")),
		},
	}

	for _, transaction := range mockTransactions {
		if err := repo.AddTransaction(&transaction); err != nil {
			log.Fatal("[FATAL] Could not add transaction in mock repo!")
		}
	}

}
