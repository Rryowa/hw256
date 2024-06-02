package storage

import (
	"bufio"
	"fmt"
	"homework-1/internal/entities"
	"homework-1/internal/file"
	"homework-1/internal/util"
	"homework-1/pkg/hash"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

type OrderStorage struct {
	fileService file.FileService
	//To optimize search
	orders map[string]entities.Order

	//To keep order in ListReturns and ListOrders
	orderIDs []string
}

func NewOrderStorage(fs file.FileService) *OrderStorage {
	return &OrderStorage{
		fileService: fs,
		orders:      make(map[string]entities.Order),
		orderIDs:    []string{},
	}
}

func (ost *OrderStorage) AcceptOrder(order entities.Order) error {
	if err := ost.fileService.CheckFile(); err != nil {
		if err = ost.fileService.CreateFile(); err != nil {
			return err
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
	return ost.fileService.Write(ost.orders)
}

func (ost *OrderStorage) ReturnOrderToCourier(orderID string) error {
	var order entities.Order
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
	return ost.fileService.Write(ost.orders)
}

func (ost *OrderStorage) IssueOrders(orderIDs []string) error {
	if len(orderIDs) == 0 {
		return util.OrderIdsNotProvidedError{}
	}

	var recipientID string
	for i, orderID := range orderIDs {
		var order entities.Order
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
	return ost.fileService.Write(ost.orders)
}

func (ost *OrderStorage) AcceptReturn(orderID, userID string) error {
	var order entities.Order
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
	return ost.fileService.Write(ost.orders)
}

func (ost *OrderStorage) ListReturns(page, pageSize int) error {
	start := (page - 1) * pageSize
	if start >= len(ost.orderIDs) {
		return nil
	}

	end := start + pageSize
	if end > len(ost.orderIDs) {
		end = len(ost.orderIDs)
	}

	var returns []entities.Order
	for _, id := range ost.orderIDs[start:end] {
		order := ost.orders[id]
		if order.Returned {
			returns = append(returns, order)
		}
	}
	ost.PrintList(returns[start:end])
	return nil
}

func (ost *OrderStorage) ListOrders(userID string, limit int) error {
	var userOrders []entities.Order

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

func (ost *OrderStorage) UpdateCache() error {
	fmt.Println("Updating cache")
	writer := bufio.NewWriter(os.Stdout)

	var err error
	ost.orders, ost.orderIDs, err = ost.fileService.Read()
	if err != nil {
		return err
	}
	for i := 0; i < 5; i++ {
		fmt.Fprint(writer, ". ")
		writer.Flush()
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Fprintln(writer)
	writer.Flush()
	return nil
}

func (ost *OrderStorage) PrintList(orders []entities.Order) {
	if len(orders) == 0 {
		fmt.Println("There are no orders or they all issued!")
		return
	}
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
