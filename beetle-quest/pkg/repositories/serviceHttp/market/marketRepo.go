package repository

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/sony/gobreaker/v2"
)

var (
	findAuctionByIDEndpoint = utils.FindEnv("FIND_AUCTION_BY_ID_ENDPOINT")
	getAllAuctionsEndpoint  = utils.FindEnv("GET_ALL_AUCTIONS_ENDPOINT")

	getAllTransactionsEndpoint           = utils.FindEnv("GET_ALL_TRANSACTIONS_ENDPOINT")
	getUserTransactionHistoryEndpoint    = utils.FindEnv("GET_USER_TRANSACTION_HISTORY_ENDPOINT")
	deleteUserTransactionHistoryEndpoint = utils.FindEnv("DELETE_USER_TRANSACTION_HISTORY_ENDPOINT")
)

type MarketRepo struct {
	client *http.Client
	cb     *gobreaker.CircuitBreaker[*http.Response]
}

func NewMarketRepo() *MarketRepo {
	return &MarketRepo{
		client: utils.SetupHTTPSClient(),
		cb:     gobreaker.NewCircuitBreaker[*http.Response](gobreaker.Settings{}),
	}
}

func (r *MarketRepo) GetUserTransactionHistory(userID models.UUID) ([]models.Transaction, bool) {
	requestData := models.GetUserTransactionHistoryData{
		UserID: userID,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return []models.Transaction{}, false
	}

	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Post(
			getUserTransactionHistoryEndpoint,
			"application/json",
			bytes.NewBuffer(jsonData),
		)

		if err != nil {
			return nil, err
		}

		return resp, nil
	})
	defer resp.Body.Close()

	if err != nil {
		return []models.Transaction{}, false
	}

	if resp.StatusCode != http.StatusOK {
		return []models.Transaction{}, false
	}

	var result models.GetUserTransactionHistoryDataResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return []models.Transaction{}, false
	}

	return result.TransactionHistory, true
}

func (r *MarketRepo) DeleteUserTransactionHistory(userID models.UUID) bool {
	requestData := models.DeleteUserTransactionHistoryData{
		UserID: userID,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return false
	}

	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Post(
			deleteUserTransactionHistoryEndpoint,
			"application/json",
			bytes.NewBuffer(jsonData),
		)

		if err != nil {
			return nil, err
		}

		return resp, nil
	})
	defer resp.Body.Close()

	if err != nil {
		return false
	}

	if resp.StatusCode != http.StatusOK {
		return false
	}

	return true
}

func (r *MarketRepo) GetAll() ([]models.Auction, bool) {
	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Get(getAllAuctionsEndpoint)
		if err != nil {
			return nil, err
		}

		return resp, nil
	})
	defer resp.Body.Close()

	if err != nil {
		return []models.Auction{}, false
	}

	if resp.StatusCode != http.StatusOK {
		return []models.Auction{}, false
	}

	var result models.GetAllAuctionDataResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return []models.Auction{}, false
	}

	return result.AuctionList, true
}

func (r *MarketRepo) GetAllTransactions() ([]models.Transaction, bool) {
	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Get(getAllTransactionsEndpoint)
		if err != nil {
			return nil, err
		}

		return resp, nil
	})
	defer resp.Body.Close()

	if err != nil {
		return []models.Transaction{}, false
	}

	if resp.StatusCode != http.StatusOK {
		return []models.Transaction{}, false
	}

	var result models.GetAllTransactionDataResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return []models.Transaction{}, false
	}

	return result.TransactionHistory, true
}

func (r *MarketRepo) FindByID(aid models.UUID) (*models.Auction, bool) {
	requestData := models.FindAuctionByIDData{
		AuctionID: aid,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, false
	}

	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Post(
			findAuctionByIDEndpoint,
			"application/json",
			bytes.NewBuffer(jsonData),
		)

		if err != nil {
			return nil, err
		}

		return resp, nil
	})
	defer resp.Body.Close()

	if err != nil {
		return nil, false
	}

	if resp.StatusCode != http.StatusOK {
		return nil, false
	}

	var result models.FindAuctionByIDDataResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, false
	}

	return result.Auction, true
}

// Not to be implemented, never used

func (r *MarketRepo) Create(*models.Auction) bool { return false }
func (r *MarketRepo) Update(*models.Auction) bool { return false }
func (r *MarketRepo) Delete(*models.Auction) bool { return false }
func (r *MarketRepo) GetUserAuctions(models.UUID) ([]models.Auction, bool) {
	return []models.Auction{}, false
}
func (r *MarketRepo) GetBidListOfAuction(models.UUID) ([]models.Bid, bool) {
	return []models.Bid{}, false
}
func (r *MarketRepo) BidToAuction(*models.Bid) bool           { return false }
func (r *MarketRepo) AddTransaction(*models.Transaction) bool { return false }
