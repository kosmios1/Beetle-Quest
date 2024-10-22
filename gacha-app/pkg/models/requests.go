package models

// Auction ===================================

type CreateAuctionReq struct {
	OwnerID string `json:"owner_id"`
	GachaID string `json:"gacha_id"`
	EndTime string `json:"end_time"`
}

type CreateAuctionRes struct {
	Auction *Auction `json:"auction"`
}

type GetAuctionRes struct {
	Auction *Auction `json:"auction"`
}
