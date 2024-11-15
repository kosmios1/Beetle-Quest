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
	UserID          string        `json:"user_id"`
	Username        string        `json:"username"`
	Email           string        `json:"email"`
	Currency        int64         `json:"currency"`
	GachaList       []Gacha       `json:"gacha_list"`
	TransactionList []Transaction `json:"transaction_list"`
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
	Price     int64  `json:"price"`
	ImagePath string `json:"image_path"`
}

type GetGachaListResponse struct {
	Gachas []GetGachaDetailsResponse `json:"gachas"`
}

// Market ===================================
type BuyBugscoinRequest struct {
	Amount string `json:"amount"`
}

type CreateAuctionRequest struct {
	GachaID string `json:"gacha_id"`
	EndTime string `json:"end_time"`
}

type BidRequest struct {
	BidAmount string `json:"bid_amount"`
}

// ============================================
// Internal models
// ============================================

// User
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

type FindUserByEmailData struct {
	Email string `json:"email"`
}

// Gacha
type GetUserGachasData struct {
	UserID UUID `json:"user_id"`
}

type GetAllGachasDataResponse struct {
	GachaList []Gacha `json:"GachaList"`
}

type GetUserGachasDataResponse struct {
	GachaList []Gacha `json:"GachaList"`
}

type AddGachaToUserData struct {
	UserID  UUID `json:"user_id"`
	GachaID UUID `json:"gacha_id"`
}

type FindGachaByIDData struct {
	GachaID string `json:"gacha_id"`
}
