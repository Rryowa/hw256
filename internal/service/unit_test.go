package service

import (
	"github.com/stretchr/testify/suite"
	"homework/internal/models"
	"homework/mocks"
	"strconv"
	"testing"
	"time"
)

type UnitTestSuite struct {
	suite.Suite
	orderService      *mocks.MockOrderService
	packageService    *mocks.MockPackageService
	validationService ValidationService
	schemaName        string

	inputDto      models.Dto
	expectedOrder models.Order
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (uts *UnitTestSuite) SetupTest() {
	uts.orderService = mocks.NewMockOrderService(uts.T())
	uts.packageService = mocks.NewMockPackageService(uts.T())
	uts.validationService = NewValidationService(uts.orderService, uts.packageService)
	uts.schemaName = "public"
	dto := models.Dto{
		ID:           "1",
		UserID:       "1",
		StorageUntil: "2077-01-01",
		OrderPrice:   "999.99",
		Weight:       "10",
		PackageType:  "box",
	}
	uts.inputDto = dto
	storageUntilDate, _ := time.Parse(time.DateOnly, dto.StorageUntil)
	orderPriceFloat, _ := strconv.ParseFloat(dto.OrderPrice, 64)
	weightFloat, _ := strconv.ParseFloat(dto.Weight, 64)
	uts.expectedOrder = models.Order{
		ID:           dto.ID,
		UserID:       dto.UserID,
		StorageUntil: storageUntilDate,
		OrderPrice:   models.Price(orderPriceFloat),
		Weight:       models.Weight(weightFloat),
	}
}
