package models

type constError string

func (e constError) Error() string {
	return string(e)
}

const (
	ErrInternalServerError          constError = "internal server error"
	ErrInvalidEndTime               constError = "invalid end time"
	ErrInvalidUserID                constError = "invalid user id"
	ErrInvalidAuctionID             constError = "invalid auction id"
	ErrCouldNotGenerateAuction      constError = "could not generate auction"
	ErrUserNotFound                 constError = "user not found"
	ErrInvalidUsernameOrPass        constError = "invalid username or password"
	ErrInvalidUsernameOrPassOrEmail constError = "invalid username or password or email"
	ErrCheckingPassword             constError = "error when checking password"
	ErrUserParametersNotValid       constError = "inserted username or mail are already in the system"
	ErrCouldNotDelete               constError = "could not delete user"
	ErrInvalidPassword              constError = "invalid password"
	ErrCouldNotUseNewPassword       constError = "could not use new password"
	ErrCouldNotUpdate               constError = "could not update user"
	ErrUsernameAlreadyExists        constError = "username already exists"
	ErrEmailAlreadyExists           constError = "email already exists"
	ErrAmountNotValid               constError = "amount not valid"
	ErrNotEnoughMoneyToBuyGacha     constError = "not enough money to buy gacha"
	ErrCouldNotAddGachaToUser       constError = "could not add gacha to user"

	ErrCouldNotParseTime          constError = "time format not correct"
	ErrCouldNotDecodeUserID       constError = "could not decode user id"
	ErrCouldNotDecodeGachaID      constError = "could not decode gacha id"
	ErrCouldNotDecodeAuctionID    constError = "could not decode auction id"
	ErrCouldNotFindResourceByUUID constError = "could not find resource by uuid"

	ErrInvalidGachaID constError = "invalid gacha id"
	ErrGachaNotFound  constError = "gacha not found"

	ErrAuctionNotFound       constError = "auction not found"
	ErrAuctionAltreadyExists constError = "auction already exists"

	// Internal api errors
	ErrCouldNotBuyGacha constError = "internal: could not buy gacha (db)"
)
