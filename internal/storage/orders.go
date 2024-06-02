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
	//To optimize search
	orders map[string]Order

	//To keep order in ListReturns and ListOrders
	orderIDs []string
}

func NewOrderStorage(fileName string) OrderStorage {
	return OrderStorage{
		fileName: fileName,
		orders:   make(map[string]Order),
		orderIDs: []string{},
	}
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

	if _, exists := ost.orders[order.ID]; !exists {
		ost.orderIDs = append(ost.orderIDs, order.ID)
	} else {
		return util.ExistingOrderError{}
	}

	fmt.Println("Calculating hash...")

	//to skip generating in tests
	if len(order.Hash) == 0 {
		order.Hash = hash.GenerateHash()
	}

	ost.orders[order.ID] = order

	fmt.Println("Order accepted.")
	return ost.writeOrders(ost.orders)
}

func (ost OrderStorage) ReturnOrderToCourier(orderID string) error {
	var order Order
	if v, exists := ost.orders[orderID]; exists {
		order = v
	} else {
		return util.OrderNotFoundError{}
	}

	if order.Issued {
		return util.OrderIssuedError{}
	}

	if !time.Now().After(order.StorageUntil) {
		return util.OrderIsNotExpiredError{}
	}

	delete(ost.orders, orderID)
	log.Println("Order returned.")
	return ost.writeOrders(ost.orders)
}

func (ost OrderStorage) IssueOrders(orderIDs []string) error {
	if len(orderIDs) == 0 {
		return util.OrderIdsNotProvidedError{}
	}

	var recipientID string
	for i, orderID := range orderIDs {
		var order Order
		if v, exists := ost.orders[orderID]; exists {
			order = v
		} else {
			return util.OrderNotFoundError{}
		}

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
		order.Issued = true
		order.IssuedAt = time.Now()
		ost.orders[order.ID] = order
		log.Println("Order issued.")
	}
	return ost.writeOrders(ost.orders)
}

func (ost OrderStorage) AcceptReturn(orderID, userID string) error {
	var order Order
	if v, exists := ost.orders[orderID]; exists {
		order = v
	} else {
		return util.OrderNotFoundError{}
	}

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
	ost.orders[order.ID] = order

	fmt.Println("Return accepted.")
	return ost.writeOrders(ost.orders)
}

func (ost OrderStorage) ListReturns(page, pageSize int) error {
	start := (page - 1) * pageSize
	if start >= len(ost.orderIDs) {
		return nil
	}

	end := start + pageSize
	if end > len(ost.orderIDs) {
		end = len(ost.orderIDs)
	}

	var returns []Order
	for _, id := range ost.orderIDs[start:end] {
		returns = append(returns, ost.orders[id])
	}
	ost.PrintList(returns[start:end])
	return nil
}

func (ost OrderStorage) ListOrders(userID string, limit int) error {
	var userOrders []Order
	//If Order issued it cant be displayed
	for _, id := range ost.orderIDs {
		order := ost.orders[id]
		if order.RecipientID == userID && !order.Issued {
			userOrders = append(userOrders, order)
			if len(userOrders) == limit {
				break
			}
		}
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

func (ost OrderStorage) writeOrders(orders map[string]Order) error {
	bWrite, err := json.MarshalIndent(orders, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ost.fileName, bWrite, 0666)
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
