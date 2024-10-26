package models

type constError string

func (e constError) Error() string {
	return string(e)
}

const (
	// Services errors
	ErrInternalServerError          constError = "internal server error"
	ErrInvalidEndTime               constError = "invalid end time"
	ErrInvalidUserID                constError = "invalid user id"
	ErrInvalidGachaID               constError = "invalid gacha id"
	ErrInvalidAuctionID             constError = "invalid auction id"
	ErrCouldNotGenerateAuction      constError = "could not generate auction"
	ErrUserNotFound                 constError = "user not found"
	ErrInvalidUsernameOrPass        constError = "invalid username or password"
	ErrInvalidUsernameOrPassOrEmail constError = "invalid username or password or email"
	ErrCheckingPassword             constError = "error when checking password"

	// Controllers errors
	ErrCouldNotParseTime          constError = "time format not correct"
	ErrCouldNotDecodeUserID       constError = "could not decode user id"
	ErrCouldNotDecodeGachaID      constError = "could not decode gacha id"
	ErrCouldNotDecodeAuctionID    constError = "could not decode auction id"
	ErrCouldNotFindResourceByUUID constError = "could not find resource by uuid"

	// Repositories errors
	ErrAuctionNotFound       constError = "auction not found"
	ErrAuctionAltreadyExists constError = "auction already exists"
)
