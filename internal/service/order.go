package service

import (
	"errors"
	"fmt"
	"homework/internal/models"
	"homework/internal/storage"
	"homework/internal/util"
	"homework/pkg/hash"
	"strconv"
	"strings"
	"time"
)

type OrderService interface {
	Accept(id, userId, dateStr, orderPrice, weight, packageType string) error
	ReturnToCourier(id string) error
	Issue(ids []string) error
	Return(id, userId string) error
	ListReturns(offset, limit string) ([]models.Order, error)
	ListOrders(userId, offset, limit string) ([]models.Order, error)
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
	if len(weight) == 0 {
		return util.ErrWeightNotProvided
	}
	if len(orderPrice) == 0 {
		return util.ErrPriceNotProvided
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

	//Check for existence
	_, err = os.repository.Get(id)
	if err == nil {
		return util.ErrOrderExists
	}

	order := createOrder(id, userId, storageUntil, orderPriceFloat, weightFloat)

	if err := applyPackaging(&order, packageType); err != nil {
		return err
	}

	calculateHash(&order)

	return os.repository.Insert(order)
}

func (os *orderService) Issue(ids []string) error {
	if len(ids) == 0 {
		return util.ErrUserIdNotProvided
	}

	var ordersToIssue []models.Order

	order, err := os.repository.Get(ids[0])
	if err != nil {
		return util.ErrOrderNotFound
	}
	recipientID := order.UserID

	for _, id := range ids {
		order, err = os.repository.Get(id)
		if err != nil {
			return util.ErrOrderNotFound
		}

		if order.Issued {
			return util.ErrOrderIssued
		}
		if order.Returned {
			return util.ErrOrderReturned
		}
		if time.Now().After(order.StorageUntil) {
			return util.ErrOrderExpired
		}

		//Check if users are equal
		if order.UserID != recipientID {
			return util.ErrOrdersUserDiffers
		}

		ordersToIssue = append(ordersToIssue, order)
	}

	modifiedOrders := issueOrders(ordersToIssue)

	return os.repository.IssueUpdate(modifiedOrders)
}

func (os *orderService) Return(id, userId string) error {
	if len(id) == 0 {
		return util.ErrOrderIdNotProvided
	}

	if len(userId) == 0 {
		return util.ErrUserIdNotProvided
	}

	order, err := os.repository.Get(id)
	if err != nil {
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

	order, err := os.repository.Get(id)
	if err != nil {
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

func (os *orderService) ListReturns(offset, limit string) ([]models.Order, error) {
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return nil, err
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}

	return os.repository.GetReturns(offsetInt, limitInt)
}

func (os *orderService) ListOrders(userId, offset, limit string) ([]models.Order, error) {
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return nil, err
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}

	return os.repository.GetOrders(userId, offsetInt, limitInt)
}

func (os *orderService) PrintList(orders []models.Order) {
	if len(orders) == 0 {
		defer fmt.Printf("\n\n")
	}
	fmt.Printf("%-5s%-10s%-15s%-15v%-10v%-13v%-10v%-13s%-13v\n", "id", "user_id", "storage_until", "issued_at", "returned", "order_price", "weight", "package_type", "package_price")
	fmt.Println(strings.Repeat("-", 100))
	for _, order := range orders {
		fmt.Printf("%-5s%-10s%-15s%-15v%-10v%-13v%-10v%-13s%-13v\n",
			order.ID,
			order.UserID,
			order.StorageUntil.Format("2006-01-02"),
			order.IssuedAt.Format("2006-01-02"),
			order.Returned,
			order.OrderPrice,
			order.Weight,
			order.PackageType,
			order.PackagePrice)
	}
	fmt.Printf("\n")
}

func createOrder(id, userId string, storageUntil time.Time, orderPrice, weight float64) models.Order {
	order := models.Order{
		ID:           id,
		UserID:       userId,
		StorageUntil: storageUntil,
		Issued:       false,
		Returned:     false,
		OrderPrice:   orderPrice,
		Weight:       weight,
	}

	return order
}

func applyPackaging(order *models.Order, packageType string) error {
	var pkg PackageInterface

	switch PackageType(packageType) {
	case FilmType:
		pkg = NewFilmPackage()
	case PacketType:
		pkg = NewPacketPackage()
	case BoxType:
		pkg = NewBoxPackage()
	case "":
		pkg = ChoosePackage(order.Weight)
	default:
		return util.ErrPackageTypeInvalid
	}

	p := NewPackage(pkg)

	if err := p.Validate(order.Weight); err != nil {
		return err
	}

	//Apply packaging and calculate order price
	order.PackageType = p.GetType()
	order.PackagePrice = p.GetPrice()
	order.OrderPrice += p.GetPrice()

	return nil
}

func calculateHash(order *models.Order) {
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
	}(order, ticker, done)

	<-done
}

func issueOrders(orders []models.Order) []models.Order {
	modifiedOrders := make([]models.Order, 0)
	for _, order := range orders {
		order.Issued = true
		order.IssuedAt = time.Now()
		modifiedOrders = append(modifiedOrders, order)
	}

	return modifiedOrders
}
