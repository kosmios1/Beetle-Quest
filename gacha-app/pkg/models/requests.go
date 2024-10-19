package models

type CreateAuctionReq struct {
	OwnerID string `json:"owner_id"`
	GachaID string `json:"gacha_id"`
	EndTime string `json:"end_time"`
}

type CreateAuctionRes struct {
	Auction Auction `json:"auction"`
}
