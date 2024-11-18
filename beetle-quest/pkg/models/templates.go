package models

type AuctionTemplate struct {
	Auction
	GachaName     string
	ImagePath     string
	OwnerUsername string
}

type UserDetailsTemplate struct {
	User         *User
	Gachas       []Gacha
	Auctions     []Auction
	Transactions []Transaction
}
