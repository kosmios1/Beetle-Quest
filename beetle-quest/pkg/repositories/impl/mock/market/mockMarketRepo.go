package repository

import (
	"beetle-quest/pkg/models"
	"sync"
)

type MarketRepo struct {
	mux sync.RWMutex

	auctions    map[models.UUID]models.Auction
	auctionBids map[models.UUID][]models.Bid

	transactions map[models.UUID]models.Transaction
}

func NewMarketRepo() *MarketRepo {
	repo := &MarketRepo{
		// NOTE: we should use more than one mutex but for the sake of simplicity,
		// being a mock repo, we will use only one
		mux:          sync.RWMutex{},
		auctions:     make(map[models.UUID]models.Auction),
		auctionBids:  make(map[models.UUID][]models.Bid),
		transactions: make(map[models.UUID]models.Transaction),
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
	return nil
}

func (r *MarketRepo) Update(auction *models.Auction) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	if _, ok := r.auctions[auction.GachaID]; !ok {
		return models.ErrGachaNotFound
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
		outTransactions = append(outTransactions, t)
	}
	return outTransactions, nil
}

func (r *MarketRepo) DeleteUserTransactionHistory(uid models.UUID) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	for tid, t := range r.transactions {
		if t.UserID == uid {
			delete(r.transactions, tid)
		}
	}
	return nil
}

func (r *MarketRepo) GetUserTransactionHistory(uid models.UUID) ([]models.Transaction, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	var outTransactions []models.Transaction
	for _, t := range r.transactions {
		if t.UserID == uid {
			outTransactions = append(outTransactions, t)
		}
	}
	return outTransactions, nil
}

func (r *MarketRepo) AddTransaction(transaction *models.Transaction) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	if _, ok := r.transactions[transaction.TransactionID]; ok {
		return models.ErrCouldNotAddTransaction
	}
	r.transactions[transaction.TransactionID] = *transaction
	return nil
}

// Utils ================================================================================================================

func populateMockRepo(repo *MarketRepo) {
	// TODO: Polulate auctions, auctionBids, and transactions repo
}
