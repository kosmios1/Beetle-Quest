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
	getUserTransactionHistoryEndpoint = utils.FindEnv("GET_USER_TRANSACTION_HISTORY_ENDPOINT")
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

// Not to be implemented, never used

func (r *MarketRepo) Create(*models.Auction) bool      { return false }
func (r *MarketRepo) Update(*models.Auction) bool      { return false }
func (r *MarketRepo) Delete(*models.Auction) bool      { return false }
func (r *MarketRepo) GetAll() ([]models.Auction, bool) { return []models.Auction{}, false }
func (r *MarketRepo) GetUserAuctions(models.UUID) ([]models.Auction, bool) {
	return []models.Auction{}, false
}
func (r *MarketRepo) FindByID(models.UUID) (*models.Auction, bool) { return nil, false }
func (r *MarketRepo) GetBidListOfAuction(models.UUID) ([]models.Bid, bool) {
	return []models.Bid{}, false
}
func (r *MarketRepo) BidToAuction(*models.Bid) bool           { return false }
func (r *MarketRepo) AddTransaction(*models.Transaction) bool { return false }
