package util

import "errors"

var ErrInvalidDate = errors.New("error - invalid date")

var ErrOrderExists = errors.New("error - order already exists")

var ErrOrderNotFound = errors.New("error - order not found")

var ErrOrderIdInvalid = errors.New("error - order id must be number")

var ErrOrderExpired = errors.New("error - order is expired")

var ErrOrderNotIssued = errors.New("error - order is not issued")

var ErrOrderIssued = errors.New("error - order is issued")

var ErrOrderIdNotProvided = errors.New("error - order id is not provided")

var ErrUserIdNotProvided = errors.New("error - user ids is not provided")

var ErrOrdersUserDiffers = errors.New("error - order's user differs")

var ErrOrderReturned = errors.New("error - order has been returned")

var ErrOrderDoesNotBelong = errors.New("error - order does not belong to user")

var ErrReturnPeriodExpired = errors.New("error - order cant be returned (period is expired)")

var ErrCreateFile = errors.New("error - create file error")
