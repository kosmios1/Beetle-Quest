package models

type constError string

func (e constError) Error() string {
	return string(e)
}

const (
	ErrInvalidEndTime        constError = "invalid end time"
	ErrInvalidUserID         constError = "invalid user id"
	ErrInvalidGachaID        constError = "invalid gacha id"
	ErrCouldNotParseTime     constError = "time format not correct"
	ErrCouldNotDecodeUserID  constError = "could not decode user id"
	ErrCouldNotDecodeGachaID constError = "could not decode gacha id"
)
