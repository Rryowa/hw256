package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"homework-1/internal/util"
	"homework-1/pkg/hash"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

type Order struct {
	ID           string    `json:"id"`
	RecipientID  string    `json:"r_id"`
	StorageUntil time.Time `json:"storage_until"`
	Issued       bool      `json:"issued"`
	IssuedAt     time.Time `json:"issued_at"`
	Returned     bool      `json:"returned"`
	Hash         string    `json:"hash"`
}

type OrderStorage struct {
	fileName string
}

func NewOrderStorage(fileName string) OrderStorage {
	return OrderStorage{fileName: fileName}
}

func (ost OrderStorage) AcceptOrder(order Order) error {
	if _, err := os.Stat(ost.fileName); errors.Is(err, os.ErrNotExist) {
		if errCreateFile := ost.createFile(); errCreateFile != nil {
			return errCreateFile
		}
	}
	if !order.StorageUntil.After(time.Now()) {
		return util.InvalidDateError{}
	}

	orders, err := ost.readOrders()
	if err != nil {
		return err
	}

	for _, existingOrder := range orders {
		if existingOrder.ID == order.ID {
			return util.ExistingOrderError{}
		}
	}

	fmt.Println("Calculating hash...")

	//to skip generating in tests
	if len(order.Hash) == 0 {
		order.Hash = hash.GenerateHash()
	}

	orders = append(orders, order)

	fmt.Println("Order accepted.")
	return ost.writeOrders(orders)
}

func (ost OrderStorage) ReturnOrderToCourier(orderID string) error {
	orders, err := ost.readOrders()
	if err != nil {
		return err
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].ID < orders[j].ID
	})
	index := binarySearch(orders, orderID)
	if index == -1 {
		return util.OrderNotFoundError{}
	}

	order := orders[index]
	if time.Now().After(order.StorageUntil) && !order.Issued {
		orders = append(orders[:index], orders[index+1:]...)
		log.Println("Order returned.")
		return ost.writeOrders(orders)
	}

	if order.Issued {
		return util.OrderIssuedError{}
	}

	if !time.Now().After(order.StorageUntil) {
		return util.OrderIsNotExpiredError{}
	}

	return util.OrderNotFoundError{}
}

func (ost OrderStorage) IssueOrders(orderIDs []string) error {
	if len(orderIDs) == 0 {
		return util.OrderIdsNotProvidedError{}
	}
	orders, err := ost.readOrders()
	if err != nil {
		return err
	}
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].ID < orders[j].ID
	})

	var recipientID string
	var index int
	for i, orderID := range orderIDs {
		index = binarySearch(orders, orderID)
		if index == -1 {
			return util.OrderNotFoundError{}
		}

		order := orders[index]
		if time.Now().After(order.StorageUntil) {
			return util.OrderIsExpiredError{}
		}
		if order.Issued {
			return util.OrderIssuedError{}
		}
		if order.Returned {
			return util.OrdersReturnedError{}
		}

		//check if recipients equal
		if i == 0 {
			recipientID = order.RecipientID
		} else {
			if order.RecipientID != recipientID {
				return util.OrdersRecipientDiffersError{}
			}
		}
		orders[index].Issued = true
		orders[index].IssuedAt = time.Now()
		log.Println("Order issued.")
	}
	return ost.writeOrders(orders)
}

func (ost OrderStorage) AcceptReturn(orderID, userID string) error {
	orders, err := ost.readOrders()
	if err != nil {
		return err
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].ID < orders[j].ID
	})
	index := binarySearch(orders, orderID)
	if index == -1 {
		return util.OrderNotFoundError{}
	}

	order := orders[index]
	if order.RecipientID != userID {
		return util.OrderDoesNotBelongError{}
	}
	if !order.Issued {
		return util.OrderHasNotBeenIssuedError{}
	}
	if time.Now().After(order.IssuedAt.Add(48 * time.Hour)) {
		return util.OrderCantBeReturnedError{}
	}

	order.Returned = true
	orders[index] = order

	fmt.Println("Return accepted.")
	return ost.writeOrders(orders)
}

func (ost OrderStorage) ListReturns(page, pageSize int) error {
	orders, err := ost.readOrders()
	if err != nil {
		return err
	}

	var returns []Order
	for _, order := range orders {
		if order.Returned {
			returns = append(returns, order)
		}
	}

	start := (page - 1) * pageSize
	if start >= len(returns) {
		return nil
	}

	end := start + pageSize
	if end > len(returns) {
		end = len(returns)
	}

	ost.PrintList(returns[start:end])
	return nil
}

func (ost OrderStorage) ListOrders(userID string, limit int) error {
	orders, err := ost.readOrders()
	if err != nil {
		return err
	}

	var userOrders []Order
	//If Order issued it cant be displayed
	for _, order := range orders {
		if order.RecipientID == userID && !order.Issued {
			userOrders = append(userOrders, order)
			if len(userOrders) == limit {
				break
			}
		}
	}

	if len(userOrders) == 0 {
		return util.EmptyOrderListError{}
	}
	// Sort orders by StorageUntil date in descending order
	sort.Slice(userOrders, func(i, j int) bool {
		return userOrders[i].StorageUntil.Before(userOrders[j].StorageUntil)
	})

	ost.PrintList(userOrders)
	return nil
}

func (ost OrderStorage) createFile() error {
	f, err := os.Create(ost.fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func (ost OrderStorage) readOrders() ([]Order, error) {
	b, err := os.ReadFile(ost.fileName)
	if err != nil {
		return nil, err
	}

	var orders []Order
	if err := json.Unmarshal(b, &orders); err != nil {
		orders = []Order{}
	}

	return orders, nil
}

func (ost OrderStorage) writeOrders(orders []Order) error {
	bWrite, err := json.MarshalIndent(orders, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ost.fileName, bWrite, 0666)
}

func binarySearch(orders []Order, orderID string) int {
	low, high := 0, len(orders)-1
	for low <= high {
		mid := (low + high) / 2
		if orders[mid].ID == orderID {
			return mid
		} else if orders[mid].ID < orderID {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return -1
}

func (ost OrderStorage) PrintList(orders []Order) {
	fmt.Printf("%-20s %-20s %-20s %-10s %-20s %-10s\n", "ID", "RecipientID", "StorageUntil", "Issued", "IssuedAt", "Returned")
	fmt.Println(strings.Repeat("-", 100))
	for _, order := range orders {
		fmt.Printf("%-20s %-20s %-20s %-10v %-20s %-10v\n",
			order.ID,
			order.RecipientID,
			order.StorageUntil.Format("2006-01-02 15:04:05"),
			order.Issued,
			order.IssuedAt.Format("2006-01-02 15:04:05"),
			order.Returned)
	}
}
