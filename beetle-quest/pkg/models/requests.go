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

// User ======================================

type GetUserAccountDetailsResponse struct {
	Username     string        `json:"username"`
	Email        string        `json:"email"`
	Currency     int64         `json:"currency"`
	Gachas       []Gacha       `json:"gachas"`
	Transactions []Transaction `json:"transactions"`
}

type UpdateUserAccountDetailsRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// Gacha =====================================

type GetGachaDetailsResponse struct {
	GachaID   string `json:"gacha_id"`
	Name      string `json:"name"`
	Rarity    string `json:"rarity"`
	Price     uint64 `json:"price"`
	ImagePath string `json:"image_path"`
}

type GetGachaListResponse struct {
	Gachas []GetGachaDetailsResponse `json:"gachas"`
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
