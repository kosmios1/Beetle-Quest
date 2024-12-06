package models

import (
  "time"
)

// Auth ======================================

type RegisterRequest struct {
	Username string `json:"username" binding:"required,ascii,min=4,max=50"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,ascii,min=8"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,ascii,min=4,max=50"`
	Password string `json:"password" binding:"required,ascii,min=8"`
}

// User ======================================
type TransactionView struct {
	TransactionID   string
	TransactionType string
	UserID          string
	Amount          int64
	DateTime        time.Time
	EventType       string
	EventID         string
}

type GetUserAccountDetailsTemplatesData struct {
	UserID          string            `json:"user_id"`
	Username        string            `json:"username"`
	Email           string            `json:"email"`
	Currency        int64             `json:"currency"`
	GachaList       []Gacha           `json:"gacha_list"`
	TransactionList []TransactionView `json:"transaction_list"`
}

type UpdateUserAccountDetailsRequest struct {
	Username    string `json:"username"      binding:"required,ascii,min=4,max=50"`
	Email       string `json:"email"         binding:"required,email"`
	OldPassword string `json:"old_password"  binding:"required,ascii,min=8"`
	NewPassword string `json:"new_password"  binding:"required,ascii,min=8"`
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
	Amount string `json:"amount" binding:"required,number,min=0,max=1000000"`
}

type CreateAuctionRequest struct {
	GachaID string `json:"gacha_id"  binding:"required,uuid4"`
	EndTime string `json:"end_time"`
}

type BidRequest struct {
	BidAmount string `json:"bid_amount"  binding:"required,number"`
}

// Admin ====================================
type AdminLoginRequest struct {
	AdminID  string `json:"admin_id"  binding:"required,uuid4"`
	Password string `json:"password"  binding:"required,ascii,min=4"`
	OtpCode  string `json:"otp_code"  binding:"required,number,len=6"`
}

type AdminUpdateUserAccount struct {
	Username string `json:"username"`
	Email    string `json:"email" binding:"required,email"`
	Currency string `json:"currency"`
}

type AdminAddGachaRequest struct {
	Name      string `json:"name"        binding:"required,ascii"`
	Rarity    string `json:"rarity"      binding:"required,alpha"`
	Price     string `json:"price"       binding:"required,number"`
	ImagePath string `json:"image_path"  binding:"required,filepath"`
}

type AdminUpdateGachaRequest struct {
  Name      string `json:"name"        binding:"required,ascii"`
	Rarity    string `json:"rarity"      binding:"required,alpha"`
	Price     string `json:"price"       binding:"required,number"`
	ImagePath string `json:"image_path"  binding:"required,filepath"`
}

type GetAllAuctionDataResponse struct {
	AuctionList []Auction `json:"AuctionList"`
}

type AdminUpdateAuctionRequest struct {
	GachaID string `json:"gacha_id"  binding:"required,uuid4"`
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
	UserID UUID `json:"user_id"`
}

type FindUserByUsernameData struct {
	Username string `json:"username"`
}

type FindUserByEmailData struct {
	Email string `json:"email"`
}

type GetAllUsersDataResponse struct {
	UserList []User `json:"UserList"`
}

// Gacha
type GetUserGachasData struct {
	UserID UUID `json:"user_id"`
}

type RemoveUserGachasData struct {
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

type RemoveGachaFromUserData struct {
	UserID  UUID `json:"user_id"`
	GachaID UUID `json:"gacha_id"`
}

type FindGachaByIDData struct {
	GachaID string `json:"gacha_id"`
}

// Market
type GetUserTransactionHistoryData struct {
	UserID UUID `json:"user_id"`
}

type DeleteUserTransactionHistoryData struct {
	UserID UUID `json:"user_id"`
}

type GetUserTransactionHistoryDataResponse struct {
	TransactionHistory []Transaction `json:"TransactionHistory"`
}

type GetAllTransactionDataResponse struct {
	TransactionHistory []Transaction `json:"TransactionHistory"`
}

type FindAuctionByIDData struct {
	AuctionID UUID `json:"auction_id"`
}

type FindAuctionByIDDataResponse struct {
	Auction *Auction `json:"Auction"`
}

type GetUserAuctionsData struct {
	UserID UUID `json:"user_id"`
}

type GetUserAuctionsDataResponse struct {
	AuctionList []Auction `json:"AuctionList"`
}

// Admin

type FindAdminByIDData struct {
	AdminID UUID `json:"admin_id"`
}
