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

type GetUserAccountDetailsTemplatesData struct {
	UserID       string        `json:"user_id"`
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

type CreateUserData struct {
	Email          string `json:"email"`
	Username       string `json:"username"`
	HashedPassword []byte `json:"password"`
	Currency       int64  `json:"currency"`
}

type FindUserByIDData struct {
	UserID string `json:"user_id"`
}

type FindUserByUsernameData struct {
	Username string `json:"username"`
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

type BuyBugscoinRequest struct {
	Amount uint64 `json:"amount"`
}

type CreateAuctionRequest struct {
	GachaUUID string `json:"gacha_uuid"`
	EndTime   string `json:"end_time"`
}
