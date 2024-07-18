package util

import "errors"

var (
	ErrPriceNotProvided    = errors.New("error - price not provided")
	ErrOrderPriceInvalid   = errors.New("error - invalid order price")
	ErrWeightNotProvided   = errors.New("error - weight not provided")
	ErrWeightExceeds       = errors.New("error - weight exceeds limit for this type of package")
	ErrWeightInvalid       = errors.New("error - invalid weight")
	ErrPackageTypeInvalid  = errors.New("error - invalid package type")
	ErrParsingDate         = errors.New("error - parsing date")
	ErrDateInvalid         = errors.New("error - invalid date")
	ErrOrderExists         = errors.New("error - order already exists")
	ErrOrderNotFound       = errors.New("error - order not found")
	ErrOrderIdInvalid      = errors.New("error - order id must be number")
	ErrUserIdInvalid       = errors.New("error - user id must be number")
	ErrOffsetInvalid       = errors.New("error - offset must be number")
	ErrLimitInvalid        = errors.New("error - limit must be number")
	ErrOrderExpired        = errors.New("error - order expired")
	ErrOrderNotIssued      = errors.New("error - order not issued")
	ErrOrderIssued         = errors.New("error - order issued")
	ErrOrderIdNotProvided  = errors.New("error - order id not provided")
	ErrUserIdNotProvided   = errors.New("error - user ids not provided")
	ErrOrdersUserDiffers   = errors.New("error - order's user differs")
	ErrOrderNotExpired     = errors.New("error - order cant be returned (not expired)")
	ErrOrderDoesNotBelong  = errors.New("error - order does not belong to user")
	ErrReturnPeriodExpired = errors.New("error - order cant be returned (period is expired)")
	ErrOffsetNotProvided   = errors.New("error - offset not provided")
	ErrLimitNotProvided    = errors.New("error - limit not provided")
	ErrCacheDelete         = errors.New("error - cant delete from cache")
)