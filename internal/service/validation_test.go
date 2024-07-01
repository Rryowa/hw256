package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"homework/internal/models"
	"homework/internal/util"
	"homework/mocks"
	"strconv"
	"testing"
	"time"
)

type ValidationTestSuite struct {
	suite.Suite
	orderService      *mocks.MockOrderService
	packageService    *mocks.MockPackageService
	validationService ValidationService
	schemaName        string

	inputDto      models.Dto
	expectedOrder models.Order
}

func TestValidationTestSuite(t *testing.T) {
	suite.Run(t, new(ValidationTestSuite))
}

func (vts *ValidationTestSuite) SetupTest() {
	vts.orderService = mocks.NewMockOrderService(vts.T())
	vts.packageService = mocks.NewMockPackageService(vts.T())
	vts.validationService = NewValidationService(vts.orderService, vts.packageService)
	dto := models.Dto{
		ID:           "1",
		UserID:       "1",
		StorageUntil: "2077-01-01",
		OrderPrice:   "999.99",
		Weight:       "10",
		PackageType:  "box",
	}
	vts.inputDto = dto
	storageUntilDate, _ := time.Parse(time.DateOnly, dto.StorageUntil)
	orderPriceFloat, _ := strconv.ParseFloat(dto.OrderPrice, 64)
	weightFloat, _ := strconv.ParseFloat(dto.Weight, 64)
	vts.expectedOrder = models.Order{
		ID:           dto.ID,
		UserID:       dto.UserID,
		StorageUntil: storageUntilDate,
		OrderPrice:   models.Price(orderPriceFloat),
		Weight:       models.Weight(weightFloat),
	}
}

func (vts *ValidationTestSuite) Test_ValidateAccept_HappyPath() {
	vts.T().Run("ValidateAccept_HappyPath", func(t *testing.T) {
		t.Parallel()
		_, err := vts.validationService.ValidateAccept(vts.inputDto)

		assert.Equal(t, tt.expectedErr, err)
	})
}

type expectDtoStorage func(s models.Dto) models.Dto

func (vts *ValidationTestSuite) Test_ValidateAccept() {
	tests := []struct {
		name        string
		expectedErr error
		expectDto   expectDtoStorage
	}{
		{
			"ErrOrderIdNotProvided",
			util.ErrOrderIdNotProvided,
			func(s models.Dto) models.Dto {
				s.ID = ""
				return s
			},
		},
		{
			"ErrUserIdNotProvided",
			util.ErrUserIdNotProvided,
			func(s models.Dto) models.Dto {
				s.UserID = ""
				return s
			},
		},
		{
			"ErrWeightNotProvided",
			util.ErrWeightNotProvided,
			func(s models.Dto) models.Dto {
				s.Weight = ""
				return s
			},
		},
		{
			"ErrPriceNotProvided",
			util.ErrPriceNotProvided,
			func(s models.Dto) models.Dto {
				s.OrderPrice = ""
				return s
			},
		},
		{
			"ErrParsingDate",
			util.ErrParsingDate,
			func(s models.Dto) models.Dto {
				s.StorageUntil = "1234-56-78"
				return s
			},
		},
		{
			"ErrDateInvalid",
			util.ErrDateInvalid,
			func(s models.Dto) models.Dto {
				s.StorageUntil = "2007-07-07"
				return s
			},
		},
		{
			"ErrOrderPriceInvalid",
			util.ErrOrderPriceInvalid,
			func(s models.Dto) models.Dto {
				s.OrderPrice = "LOL"
				return s
			},
		},
		{
			"ErrOrderPriceInvalidNegative",
			util.ErrOrderPriceInvalid,
			func(s models.Dto) models.Dto {
				s.OrderPrice = "-1"
				return s
			},
		},
		{
			"ErrWeightInvalid",
			util.ErrWeightInvalid,
			func(s models.Dto) models.Dto {
				s.Weight = "LOL"
				return s
			},
		},
		{
			"ErrWeightInvalidNegative",
			util.ErrWeightInvalid,
			func(s models.Dto) models.Dto {
				s.Weight = "-1"
				return s
			},
		},
	}
	for _, tt := range tests {
		vts.T().Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := vts.validationService.ValidateAccept(tt.expectDto(vts.inputDto))

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

type expectMockOrderService func(s *mocks.MockOrderService)

func (vts *ValidationTestSuite) Test_ValidateIssue() {
	tests := []struct {
		name                   string
		expectedErr            error
		input                  string
		expectMockOrderService expectMockOrderService
	}{
		{
			"ErrUserIdNotProvided",
			util.ErrUserIdNotProvided,
			"",
			func(s *mocks.MockOrderService) {
			},
		},
		{
			"ErrOrderNotFound",
			util.ErrOrderNotFound,
			"1",
			func(s *mocks.MockOrderService) {
				s.EXPECT().Exists("1").Return(models.Order{}, false)
			},
		},
		{
			"ErrOrderIssued",
			util.ErrOrderIssued,
			"1",
			func(s *mocks.MockOrderService) {
				s.EXPECT().Exists("1").Return(models.Order{Issued: true}, true)
			},
		},
		{
			"ErrOrderReturned",
			util.ErrOrderReturned,
			"1",
			func(s *mocks.MockOrderService) {
				s.EXPECT().Exists("1").Return(models.Order{Returned: true}, true)
			},
		},
		{
			"ErrOrderExpired",
			util.ErrOrderExpired,
			"1",
			func(s *mocks.MockOrderService) {
				order := models.Order{StorageUntil: time.Now()}
				order.StorageUntil = order.StorageUntil.AddDate(-1000, 0, 0)
				s.EXPECT().Exists("1").Return(order, true)
			},
		},
	}
	for _, tt := range tests {
		vts.T().Run(tt.name, func(t *testing.T) {
			t.Parallel()
			orderSrvc := mocks.NewMockOrderService(t)
			tt.expectMockOrderService(orderSrvc)
			validation := NewValidationService(orderSrvc, mocks.NewMockPackageService(t))

			_, err := validation.ValidateIssue(tt.input)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func (vts *ValidationTestSuite) Test_ValidateAcceptReturn() {
	tests := []struct {
		name                   string
		expectedErr            error
		input                  []string
		expectMockOrderService expectMockOrderService
	}{
		{
			"ErrOrderIdNotProvided",
			util.ErrOrderIdNotProvided,
			[]string{"", "1"},
			func(s *mocks.MockOrderService) {
			},
		},
		{
			"ErrUserIdNotProvided",
			util.ErrUserIdNotProvided,
			[]string{"1", ""},
			func(s *mocks.MockOrderService) {
			},
		},
		{
			"ErrOrderNotFound",
			util.ErrOrderNotFound,
			[]string{"1", "1"},
			func(s *mocks.MockOrderService) {
				s.EXPECT().Exists("1").Return(models.Order{}, false)
			},
		},
		{
			"ErrOrderDoesNotBelong",
			util.ErrOrderDoesNotBelong,
			[]string{"1", "2"},
			func(s *mocks.MockOrderService) {
				order := vts.expectedOrder
				s.EXPECT().Exists("1").Return(order, true)
			},
		},
		{
			"ErrOrderNotIssued",
			util.ErrOrderNotIssued,
			[]string{"1", "1"},
			func(s *mocks.MockOrderService) {
				order := vts.expectedOrder
				order.Issued = false
				s.EXPECT().Exists("1").Return(order, true)
			},
		},
		{
			"ErrReturnPeriodExpired",
			util.ErrReturnPeriodExpired,
			[]string{"1", "1"},
			func(s *mocks.MockOrderService) {
				order := vts.expectedOrder
				order.Issued = true
				order.StorageUntil = order.StorageUntil.AddDate(-1000, 0, 0)
				s.EXPECT().Exists("1").Return(order, true)
			},
		},
	}
	for _, tt := range tests {
		vts.T().Run(tt.name, func(t *testing.T) {
			t.Parallel()
			orderSrvc := mocks.NewMockOrderService(t)
			tt.expectMockOrderService(orderSrvc)
			validation := NewValidationService(orderSrvc, mocks.NewMockPackageService(t))

			_, err := validation.ValidateAcceptReturn(tt.input[0], tt.input[1])

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func (vts *ValidationTestSuite) Test_ValidateReturnToCourier() {
	tests := []struct {
		name                   string
		expectedErr            error
		input                  string
		expectMockOrderService expectMockOrderService
	}{
		{
			"ErrOrderIdNotProvided",
			util.ErrOrderIdNotProvided,
			"",
			func(s *mocks.MockOrderService) {
			},
		},
		{
			"ErrOrderNotFound",
			util.ErrOrderNotFound,
			"1",
			func(s *mocks.MockOrderService) {
				s.EXPECT().Exists("1").Return(models.Order{}, false)
			},
		},
		{
			"ErrOrderIssued",
			util.ErrOrderIssued,
			"1",
			func(s *mocks.MockOrderService) {
				order := vts.expectedOrder
				order.Issued = true
				s.EXPECT().Exists("1").Return(order, true)
			},
		},
	}
	for _, tt := range tests {
		vts.T().Run(tt.name, func(t *testing.T) {
			t.Parallel()
			orderSrvc := mocks.NewMockOrderService(t)
			tt.expectMockOrderService(orderSrvc)
			validation := NewValidationService(orderSrvc, mocks.NewMockPackageService(t))

			err := validation.ValidateReturnToCourier(tt.input)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func (vts *ValidationTestSuite) Test_ValidateList_ErrNotProvided() {
	vts.T().Parallel()
	_, _, err := vts.validationService.ValidateList("", "1")
	vts.EqualError(err, util.ErrOffsetNotProvided.Error())
	_, _, err = vts.validationService.ValidateList("1", "")
	vts.EqualError(err, util.ErrLimitNotProvided.Error())
}

func (vts *ValidationTestSuite) Test_Validate_IsEmpty() {
	assert.Equal(vts.T(), true, isArgEmpty(""))
}
