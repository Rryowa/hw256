package tests_test

import (
	"errors"
	"homework-1/internal/entities"
	"homework-1/internal/file"
	"homework-1/internal/storage"
	"homework-1/internal/util"
	"log"
	"os"
	"testing"
	"time"
)

var testFile = "test.json"

func setupStorage(t *testing.T) *storage.OrderStorage {
	t.Helper()
	err := os.Remove(testFile)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}
	fileService := file.NewFileService(testFile)
	return storage.NewOrderStorage(fileService)
}

func newOrder(id, recipientID, hash string, date string, issued bool) entities.Order {
	storageUntil, _ := time.Parse(time.DateOnly, date)
	return entities.Order{
		ID:           id,
		RecipientID:  recipientID,
		StorageUntil: storageUntil,
		Issued:       issued,
		Hash:         hash,
	}
}

func TestAcceptOrderInvalidDate(t *testing.T) {
	ost := setupStorage(t)
	order := newOrder("1", "2", "abc", "2020-02-02", false)
	err := ost.AcceptOrder(order)
	if !errors.As(err, &util.InvalidDateError{}) {
		t.Fatalf("expected invalid date error, got %v", err)
	}
}

func TestAcceptExistingOrder(t *testing.T) {
	ost := setupStorage(t)
	order := newOrder("1", "2", "abc", "2077-02-02", false)
	_ = ost.AcceptOrder(order)
	order2 := newOrder("1", "3", "abcd", "2077-02-02", false)
	err := ost.AcceptOrder(order2)
	if !errors.As(err, &util.ExistingOrderError{}) {
		t.Fatalf("expected same order, got %v", err)
	}
}

func TestReturnOrderWrongId(t *testing.T) {
	ost := setupStorage(t)
	order := newOrder("1", "2", "abc", "2077-02-02", false)
	_ = ost.AcceptOrder(order)
	err := ost.ReturnOrderToCourier("100000")
	if !errors.As(err, &util.OrderNotFoundError{}) {
		t.Fatalf("expected OrderNotFoundError, got %v", err)
	}
}

func TestReturnOrderDateIsNotExpired(t *testing.T) {
	ost := setupStorage(t)
	order := newOrder("1", "2", "abc", "2077-02-02", false)
	_ = ost.AcceptOrder(order)
	err := ost.ReturnOrderToCourier("1")
	if !errors.As(err, &util.OrderIsNotExpiredError{}) {
		t.Fatalf("expected OrderIsNotExpiredError, got %v", err)
	}
}

func TestIssueOrdersEmptyIdSlice(t *testing.T) {
	ost := setupStorage(t)
	order := newOrder("1", "2", "abc", "2077-02-02", false)
	_ = ost.AcceptOrder(order)
	var s []string
	err := ost.IssueOrders(s)
	if !errors.As(err, &util.OrderIdsNotProvidedError{}) {
		t.Fatalf("expected OrderIdsNotProvidedError, got %v", err)
	}
}

func TestIssueOrderNotFound(t *testing.T) {
	ost := setupStorage(t)
	order := newOrder("1", "2", "abc", "2077-02-02", false)
	_ = ost.AcceptOrder(order)
	s := []string{"7", "8"}
	err := ost.IssueOrders(s)
	if !errors.As(err, &util.OrderNotFoundError{}) {
		t.Fatalf("expected OrderNotFoundError, got %v", err)
	}
}

func TestIssueOrderIssued(t *testing.T) {
	ost := setupStorage(t)
	order := newOrder("1", "2", "abc", "2077-02-02", false)
	_ = ost.AcceptOrder(order)
	s := []string{"1"}
	_ = ost.IssueOrders(s)
	err := ost.IssueOrders(s)
	if !errors.As(err, &util.OrderIssuedError{}) {
		t.Fatalf("expected OrderIssuedError, got %v", err)
	}
}

func TestIssueOrderRecipientDiffers(t *testing.T) {
	ost := setupStorage(t)
	order1 := newOrder("1", "2", "abc", "2077-02-02", false)
	order2 := newOrder("2", "3", "abcd", "2077-02-02", false)
	_ = ost.AcceptOrder(order1)
	_ = ost.AcceptOrder(order2)
	s := []string{"1", "2"}
	err := ost.IssueOrders(s)
	if !errors.As(err, &util.OrdersRecipientDiffersError{}) {
		t.Fatalf("expected OrdersRecipientDiffersError, got %v", err)
	}
}

func TestAcceptReturnOrderNotFound(t *testing.T) {
	ost := setupStorage(t)
	order := newOrder("1", "2", "abc", "2077-02-02", false)
	_ = ost.AcceptOrder(order)
	err := ost.AcceptReturn("666", "2")
	if !errors.As(err, &util.OrderNotFoundError{}) {
		t.Fatalf("expected OrderNotFoundError, got %v", err)
	}
}

func TestAcceptReturnOrderDoesNotBelong(t *testing.T) {
	ost := setupStorage(t)
	order := newOrder("1", "2", "abc", "2077-02-02", true)
	_ = ost.AcceptOrder(order)
	err := ost.AcceptReturn("1", "666")
	if !errors.As(err, &util.OrderDoesNotBelongError{}) {
		t.Fatalf("expected OrderDoesNotBelongError, got %v", err)
	}
}

func TestAcceptReturnOrderHasNotBeenIssued(t *testing.T) {
	ost := setupStorage(t)
	storageUntil, _ := time.Parse(time.DateOnly, "2024-06-29")
	order := entities.Order{
		ID:           "1",
		RecipientID:  "2",
		StorageUntil: storageUntil,
		Issued:       false,
		Hash:         "abc",
	}
	_ = ost.AcceptOrder(order)

	err := ost.AcceptReturn("1", "2")
	if !errors.As(err, &util.OrderHasNotBeenIssuedError{}) {
		t.Fatalf("expected OrderHasNotBeenIssuedError, got %v", err)
	}
}

func TestAcceptReturnOrderCantBeReturnedAfter48Hours(t *testing.T) {
	ost := setupStorage(t)
	storageUntil, _ := time.Parse(time.DateOnly, "2024-06-29")
	order := entities.Order{
		ID:           "1",
		RecipientID:  "2",
		StorageUntil: storageUntil,
		Issued:       true,
		IssuedAt:     time.Now().Add(-49 * time.Hour),
		Hash:         "abc",
	}
	_ = ost.AcceptOrder(order)

	err := ost.AcceptReturn("1", "2")
	if !errors.As(err, &util.OrderCantBeReturnedError{}) {
		t.Fatalf("expected OrderCantBeReturnedError, got %v", err)
	}
}

func TestAcceptReturnSuccessful(t *testing.T) {
	ost := setupStorage(t)
	storageUntil, _ := time.Parse(time.DateOnly, "2024-06-29")
	order := entities.Order{
		ID:           "1",
		RecipientID:  "2",
		StorageUntil: storageUntil,
		Issued:       true,
		IssuedAt:     time.Now().Add(-24 * time.Hour),
		Hash:         "abc",
	}
	_ = ost.AcceptOrder(order)

	err := ost.AcceptReturn("1", "2")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
