package models

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
