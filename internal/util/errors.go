package util

import "errors"

var (
	ErrOrderPriceInvalid   = errors.New("error - invalid order price")
	ErrWeightExceeds       = errors.New("error - weight exceeds limit for this type of package")
	ErrWeightInvalid       = errors.New("error - invalid weight")
	ErrPackageTypeInvalid  = errors.New("error - invalid package type")
	ErrDateInvalid         = errors.New("error - invalid date")
	ErrOrderExists         = errors.New("error - order already exists")
	ErrOrderNotFound       = errors.New("error - order not found")
	ErrOrderIdInvalid      = errors.New("error - order id must be number")
	ErrOrderExpired        = errors.New("error - order is expired")
	ErrOrderNotIssued      = errors.New("error - order is not issued")
	ErrOrderIssued         = errors.New("error - order is issued")
	ErrOrderIdNotProvided  = errors.New("error - order id is not provided")
	ErrUserIdNotProvided   = errors.New("error - user ids is not provided")
	ErrOrdersUserDiffers   = errors.New("error - order's user differs")
	ErrOrderReturned       = errors.New("error - order has been returned")
	ErrOrderDoesNotBelong  = errors.New("error - order does not belong to user")
	ErrReturnPeriodExpired = errors.New("error - order cant be returned (period is expired)")
)
