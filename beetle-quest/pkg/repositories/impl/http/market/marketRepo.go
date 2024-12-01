package repository

import (
	"beetle-quest/pkg/httpserver"
	"beetle-quest/pkg/models"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sony/gobreaker/v2"
)

var (
	findAuctionByIDEndpoint = os.Getenv("FIND_AUCTION_BY_ID_ENDPOINT")
	getAllAuctionsEndpoint  = os.Getenv("GET_ALL_AUCTIONS_ENDPOINT")
	getUserAuctionsEndpoint = os.Getenv("GET_USER_AUCTIONS_ENDPOINT")

	getAllTransactionsEndpoint           = os.Getenv("GET_ALL_TRANSACTIONS_ENDPOINT")
	getUserTransactionHistoryEndpoint    = os.Getenv("GET_USER_TRANSACTION_HISTORY_ENDPOINT")
	deleteUserTransactionHistoryEndpoint = os.Getenv("DELETE_USER_TRANSACTION_HISTORY_ENDPOINT")

	updateAuctionEndpoint = os.Getenv("UPDATE_AUCTION_ENDPOINT")
)

type MarketRepo struct {
	client *http.Client
	cb     *gobreaker.CircuitBreaker[*http.Response]
}

func NewMarketRepo() *MarketRepo {
	return &MarketRepo{
		client: httpserver.SetupHTTPSClient(),
		// closed: ok, open: service does not respond
		cb: gobreaker.NewCircuitBreaker[*http.Response](gobreaker.Settings{
			MaxRequests: 5,
			Interval:    5 * time.Second,  // When to flush counters int the Closed state
			Timeout:     10 * time.Second, // Time to switch from open -> half-open
			ReadyToTrip: func(counts gobreaker.Counts) bool { // When to switch from closed -> open
				failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
				return (counts.Requests >= 10 && failureRatio >= 0.6) || counts.ConsecutiveFailures > 10
			},
			OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
				log.Printf("[INFO] Circuit breaker changed from %s to %s", from, to)
			},
		}),
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

	log.Panicf("Unreachable code, status code received: %d", resp.StatusCode)
	return nil, models.ErrInternalServerError
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

	log.Panicf("Unreachable code, status code received: %d", resp.StatusCode)
	return models.ErrInternalServerError
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
	log.Panicf("Unreachable code, status code received: %d", resp.StatusCode)
	return nil, models.ErrInternalServerError
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
	log.Panicf("Unreachable code, status code received: %d", resp.StatusCode)
	return nil, models.ErrInternalServerError
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

	log.Panicf("Unreachable code, status code received: %d", resp.StatusCode)
	return nil, models.ErrInternalServerError
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

	log.Panicf("Unreachable code, status code received: %d", resp.StatusCode)
	return nil, models.ErrInternalServerError
}

func (r *MarketRepo) Update(auction *models.Auction) error {
	jsonData, err := json.Marshal(auction)
	if err != nil {
		return models.ErrInternalServerError
	}

	resp, err := r.cb.Execute(func() (*http.Response, error) {
		resp, err := r.client.Post(
			updateAuctionEndpoint,
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
		return models.ErrAuctionNotFound
	case http.StatusBadRequest:
		return models.ErrInvalidData
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	log.Panicf("Unreachable code, status code received: %d", resp.StatusCode)
	return models.ErrInternalServerError
}

// Not to be implemented, never used

func (r *MarketRepo) Create(*models.Auction) error { return nil }
func (r *MarketRepo) Delete(*models.Auction) error { return nil }
func (r *MarketRepo) GetBidListOfAuction(models.UUID) ([]models.Bid, error) {
	return []models.Bid{}, nil
}
func (r *MarketRepo) BidToAuction(*models.Bid) error           { return nil }
func (r *MarketRepo) AddTransaction(*models.Transaction) error { return nil }
