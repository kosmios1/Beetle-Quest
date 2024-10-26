package models

// Auth ======================================

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Auction ===================================

type CreateAuctionRequest struct {
	GachaUUID string `json:"gacha_uuid"`
	EndTime   string `json:"end_time"`
}

type CreateAuctionResponse struct {
	Auction *Auction `json:"auction"`
}

type GetAuctionResponse struct {
	Auction *Auction `json:"auction"`
}
