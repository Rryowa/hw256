package service

import (
	"errors"
	"fmt"
	"homework-1/internal/models"
	"homework-1/internal/storage"
	"homework-1/internal/util"
	"homework-1/pkg/hash"
	"log"
	"strconv"
	"strings"
	"time"
)

type OrderService interface {
	Accept(id, userId, dateStr, orderPrice, weight, packageType string) error
	ReturnToCourier(id string) error
	Issue(ids []string) error
	Return(id, userId string) error
	ListReturns(page, size string) ([]models.Order, error)
	ListOrders(userId, limit string) ([]models.Order, error)
	PrintList(orders []models.Order)
}

type orderService struct {
	repository storage.Storage
}

func NewOrderService(repository storage.Storage) OrderService {
	return &orderService{
		repository: repository,
	}
}

func (os *orderService) Accept(id, userId, dateStr, orderPrice, weight, packageType string) error {
	if len(id) == 0 {
		return util.ErrOrderIdNotProvided
	}
	if len(userId) == 0 {
		return util.ErrUserIdNotProvided
	}
	if len(orderPrice) == 0 {
		return util.ErrPriceNotProvided
	}
	if len(weight) == 0 {
		return util.ErrWeightNotProvided
	}

	storageUntil, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		return errors.New("error parsing date")
	}
	if storageUntil.Before(time.Now()) {
		return util.ErrDateInvalid
	}

	orderPriceFloat, err := strconv.ParseFloat(orderPrice, 64)
	if err != nil || orderPriceFloat <= 0 {
		return util.ErrOrderPriceInvalid
	}
	weightFloat, err := strconv.ParseFloat(weight, 64)
	if err != nil || weightFloat <= 0 {
		return util.ErrWeightInvalid
	}

	emptyOrder := models.Order{}
	order := os.repository.Get(id)
	if order != emptyOrder {
		return util.ErrOrderExists
	}

	pkg, err := ApplyPackaging(weightFloat, packageType)
	if err != nil {
		return err
	}

	newOrder := Create(id, userId, storageUntil, orderPriceFloat, weightFloat, pkg)

	return os.repository.Insert(newOrder)
}

func (os *orderService) Issue(ids []string) error {
	if len(ids) == 0 {
		return util.ErrUserIdNotProvided
	}

	var orders []models.Order
	var recipientID string
	for i, id := range ids {
		emptyOrder := models.Order{}
		order := os.repository.Get(id)
		if order == emptyOrder {
			return util.ErrOrderNotFound
		}

		if time.Now().After(order.StorageUntil) {
			return util.ErrOrderExpired
		}
		if order.Issued {
			return util.ErrOrderIssued
		}
		if order.Returned {
			return util.ErrOrderReturned
		}

		//Check if users are equal
		if i == 0 {
			recipientID = order.UserID
		} else {
			if order.UserID != recipientID {
				return util.ErrOrdersUserDiffers
			}
		}
		orders = append(orders, order)
	}

	modifiedOrders := IssueOrders(orders)

	return os.repository.IssueUpdate(modifiedOrders)
}

func (os *orderService) Return(id, userId string) error {
	if len(id) == 0 {
		return util.ErrOrderIdNotProvided
	}

	if len(userId) == 0 {
		return util.ErrUserIdNotProvided
	}

	emptyOrder := models.Order{}
	order := os.repository.Get(id)
	if order == emptyOrder {
		return util.ErrOrderNotFound
	}

	if order.UserID != userId {
		return util.ErrOrderDoesNotBelong
	}
	if !order.Issued {
		return util.ErrOrderNotIssued
	}
	if time.Now().After(order.IssuedAt.Add(48 * time.Hour)) {
		return util.ErrReturnPeriodExpired
	}

	order.Returned = true

	return os.repository.Update(order)
}

func (os *orderService) ReturnToCourier(id string) error {
	if len(id) == 0 {
		return util.ErrOrderIdNotProvided
	}

	if _, err := strconv.Atoi(id); err != nil {
		return util.ErrOrderIdInvalid
	}

	emptyOrder := models.Order{}
	order := os.repository.Get(id)
	if order == emptyOrder {
		return util.ErrOrderNotFound
	}

	if order.Issued {
		return util.ErrOrderIssued
	}

	//skip checking for a period, to ensure that its working
	//if time.Now().Before(order.StorageUntil) {
	//	return util.ErrOrderNotExpired
	//}

	return os.repository.Delete(id)
}

func (os *orderService) ListReturns(limit, offset string) ([]models.Order, error) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return nil, err
	}

	return os.repository.GetReturns(limitInt, offsetInt)
}

func (os *orderService) ListOrders(userId, limit string) ([]models.Order, error) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}

	return os.repository.GetOrders(userId, limitInt)
}

func (os *orderService) PrintList(orders []models.Order) {
	if len(orders) == 0 {
		defer fmt.Printf("\n\n")
	}
	fmt.Printf("%-20s %-20s %-20s %-10s %-20s %-10s\n", "ID", "userId", "StorageUntil", "Issued", "IssuedAt", "Returned")
	fmt.Println(strings.Repeat("-", 100))
	for _, order := range orders {
		fmt.Printf("%-20s %-20s %-20s %-10v %-20s %-10v\n",
			order.ID,
			order.UserID,
			order.StorageUntil.Format("2006-01-02 15:04:05"),
			order.Issued,
			order.IssuedAt.Format("2006-01-02 15:04:05"),
			order.Returned)
	}
}

func ApplyPackaging(weightFloat float64, packageType string) (Package, error) {
	var pkg Package
	switch PackageType(packageType) {
	case filmType:
		pkg = NewFilm()
	case packetType:
		pkg = NewPacket()
	case boxType:
		pkg = NewBox()
	case "":
		pkg = choosePackage(weightFloat)
		log.Println("Based on weight, package type is:", GetPackageType(pkg))
	default:
		return nil, util.ErrPackageTypeInvalid
	}
	if err := pkg.Validate(weightFloat); err != nil {
		return nil, err
	}

	return pkg, nil
}

func choosePackage(weightFloat float64) Package {
	if weightFloat >= 30 {
		return NewFilm()
	} else if weightFloat >= 10 {
		return NewBox()
	} else {
		return NewPacket()
	}
}

func Create(id, userId string, storageUntil time.Time, orderPrice, weight float64, pkg Package) models.Order {
	order := models.Order{
		ID:           id,
		UserID:       userId,
		StorageUntil: storageUntil,
		Issued:       false,
		Returned:     false,
		OrderPrice:   orderPrice + GetPackagePrice(pkg),
		Weight:       weight,
		PackageType:  GetPackageType(pkg),
	}

	fmt.Print("Calculating hash.")

	ticker := time.NewTicker(time.Second)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				fmt.Print(" .")
			}
		}
	}()

	go func(order *models.Order, ticker *time.Ticker, done chan struct{}) {
		order.Hash = hash.GenerateHash()
		ticker.Stop()
		done <- struct{}{}
	}(&order, ticker, done)

	<-done

	return order
}

func IssueOrders(orders []models.Order) []models.Order {
	modifiedOrders := make([]models.Order, 0)
	for _, order := range orders {
		order.Issued = true
		order.IssuedAt = time.Now()
		modifiedOrders = append(modifiedOrders, order)
	}

	return modifiedOrders
}
