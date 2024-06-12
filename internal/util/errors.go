package util

type CreateFileError struct{}

func (e CreateFileError) Error() string {
	return "create storage error"
}

type InvalidDateError struct{}

func (e InvalidDateError) Error() string {
	return "invalid date"
}

type ExistingOrderError struct{}

func (e ExistingOrderError) Error() string {
	return "order already exists"
}

type OrderNotFoundError struct{}

func (e OrderNotFoundError) Error() string {
	return "order not found"
}

type OrderIsExpiredError struct{}

func (e OrderIsExpiredError) Error() string {
	return "order is expired"
}

type OrderIsNotExpiredError struct{}

func (e OrderIsNotExpiredError) Error() string {
	return "order is not expired"
}

type OrderIssuedError struct{}

func (e OrderIssuedError) Error() string {
	return "order has been issued"
}

type OrderIdsNotProvidedError struct{}

func (e OrderIdsNotProvidedError) Error() string {
	return "order id is not provided"
}

type UserIdIsNotProvided struct{}

func (e UserIdIsNotProvided) Error() string {
	return "user id is not provided"
}

type OrdersRecipientDiffersError struct{}

func (e OrdersRecipientDiffersError) Error() string {
	return "all orders must belong to the same recipient"
}

type OrdersReturnedError struct{}

func (e OrdersReturnedError) Error() string {
	return "order has been returned"
}

type OrderDoesNotBelongError struct{}

func (e OrderDoesNotBelongError) Error() string {
	return "order does not belong to recipient"
}

type OrderHasNotBeenIssuedError struct{}

func (e OrderHasNotBeenIssuedError) Error() string {
	return "order has not been issued"
}

type OrderCantBeReturnedError struct{}

func (e OrderCantBeReturnedError) Error() string {
	return "return period has expired"
}
