package repository

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/sony/gobreaker/v2"
)

var (
	findAuctionByIDEndpoint = os.Getenv("FIND_AUCTION_BY_ID_ENDPOINT")
	getAllAuctionsEndpoint  = os.Getenv("GET_ALL_AUCTIONS_ENDPOINT")
	getUserAuctionsEndpoint = os.Getenv("GET_USER_AUCTIONS_ENDPOINT")

	getAllTransactionsEndpoint           = os.Getenv("GET_ALL_TRANSACTIONS_ENDPOINT")
	getUserTransactionHistoryEndpoint    = os.Getenv("GET_USER_TRANSACTION_HISTORY_ENDPOINT")
	deleteUserTransactionHistoryEndpoint = os.Getenv("DELETE_USER_TRANSACTION_HISTORY_ENDPOINT")
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

func (r *MarketRepo) GetUserTransactionHistory(userID models.UUID) ([]models.Transaction, error) {
	requestData := models.GetUserTransactionHistoryData{
		UserID: userID,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, models.ErrInternalServerError
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
	if err != nil {
		return nil, models.ErrInternalServerError
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		return nil, models.ErrInternalServerError
	case http.StatusNotFound:
		return nil, models.ErrTransactionNotFound
	}

	if resp.StatusCode == http.StatusOK {
		var result models.GetUserTransactionHistoryDataResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return nil, models.ErrInternalServerError
		}
		return result.TransactionHistory, nil
	}

	panic("unreachable code")
}

func (r *MarketRepo) DeleteUserTransactionHistory(userID models.UUID) error {
	requestData := models.DeleteUserTransactionHistoryData{
		UserID: userID,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return models.ErrInternalServerError
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
	if err != nil {
		return models.ErrInternalServerError
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		return models.ErrInternalServerError
	case http.StatusNotFound:
		return models.ErrTransactionNotFound
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	panic("unreachable code")
}

func (r *MarketRepo) GetAll() ([]models.Auction, error) {
	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Get(getAllAuctionsEndpoint)
		if err != nil {
			return nil, err
		}

		return resp, nil
	})
	if err != nil {
		return nil, models.ErrInternalServerError
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		return nil, models.ErrInternalServerError
	case http.StatusNotFound:
		return nil, models.ErrAuctionNotFound
	}

	if resp.StatusCode == http.StatusOK {
		var result models.GetAllAuctionDataResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return nil, models.ErrInternalServerError
		}

		return result.AuctionList, nil
	}
	panic("unreachable code")
}

func (r *MarketRepo) GetAllTransactions() ([]models.Transaction, error) {
	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Get(getAllTransactionsEndpoint)
		if err != nil {
			return nil, err
		}

		return resp, nil
	})
	if err != nil {
		return nil, models.ErrInternalServerError
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		return nil, models.ErrInternalServerError
	case http.StatusNotFound:
		return nil, models.ErrAuctionNotFound
	}

	if resp.StatusCode == http.StatusOK {
		var result models.GetAllTransactionDataResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return nil, models.ErrInternalServerError
		}
		return result.TransactionHistory, nil
	}
	panic("unreachable code")
}

func (r *MarketRepo) FindByID(aid models.UUID) (*models.Auction, error) {
	requestData := models.FindAuctionByIDData{
		AuctionID: aid,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, models.ErrInternalServerError
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
	if err != nil {
		return nil, models.ErrInternalServerError
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		return nil, models.ErrInternalServerError
	case http.StatusNotFound:
		return nil, models.ErrAuctionNotFound
	}

	if resp.StatusCode == http.StatusOK {

		var result models.FindAuctionByIDDataResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return nil, models.ErrInternalServerError
		}

		return result.Auction, nil
	}

	panic("unreachable code")
}

func (r *MarketRepo) GetUserAuctions(uid models.UUID) ([]models.Auction, error) {
	requestData := models.GetUserAuctionsData{
		UserID: uid,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, models.ErrInternalServerError
	}

	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Post(
			getUserAuctionsEndpoint,
			"application/json",
			bytes.NewBuffer(jsonData),
		)

		if err != nil {
			return nil, err
		}

		return resp, nil
	})
	if err != nil {
		return nil, models.ErrInternalServerError
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		return nil, models.ErrInternalServerError
	case http.StatusNotFound:
		return nil, models.ErrAuctionNotFound
	}

	if resp.StatusCode == http.StatusOK {
		var result models.GetUserAuctionsDataResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return nil, models.ErrInternalServerError
		}
		return result.AuctionList, nil
	}

	panic("unreachable code")
}

// Not to be implemented, never used

func (r *MarketRepo) Create(*models.Auction) error { return nil }
func (r *MarketRepo) Update(*models.Auction) error { return nil }
func (r *MarketRepo) Delete(*models.Auction) error { return nil }
func (r *MarketRepo) GetBidListOfAuction(models.UUID) ([]models.Bid, error) {
	return []models.Bid{}, nil
}
func (r *MarketRepo) BidToAuction(*models.Bid) error           { return nil }
func (r *MarketRepo) AddTransaction(*models.Transaction) error { return nil }
