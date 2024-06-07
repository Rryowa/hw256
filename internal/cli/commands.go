package cli

const (
	help                 = "help"
	acceptOrder          = "accept"
	returnOrderToCourier = "return_courier"
	issueOrders          = "issue"
	acceptReturn         = "accept_return"
	listReturns          = "list_returns"
	listOrders           = "list_orders"
	setWorkers           = "set_workers"
	exit                 = "exit"
)

type command struct {
	name        string
	description string
}
