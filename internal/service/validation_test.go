package service

import (
	"github.com/stretchr/testify/assert"
	"homework/internal/models"
	"homework/internal/util"
	"homework/mocks"
	"testing"
	"time"
)

func Test_Validate_IsEmpty(t *testing.T) {
	assert.Equal(t, true, isArgEmpty(""))
}

type expectDtoStorage func(s models.Dto) models.Dto

func (uts *UnitTestSuite) Test_ValidateAccept() {
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
		uts.T().Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := uts.validationService.ValidateAccept(tt.expectDto(uts.inputDto))

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

type expectMockOrderService func(s *mocks.MockOrderService)

func (uts *UnitTestSuite) Test_ValidateIssue() {
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
		uts.T().Run(tt.name, func(t *testing.T) {
			t.Parallel()
			orderSrvc := mocks.NewMockOrderService(t)
			tt.expectMockOrderService(orderSrvc)
			validation := NewValidationService(orderSrvc, mocks.NewMockPackageService(t))

			_, err := validation.ValidateIssue(tt.input)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func (uts *UnitTestSuite) Test_ValidateAcceptReturn() {
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
				order := uts.expectedOrder
				s.EXPECT().Exists("1").Return(order, true)
			},
		},
		{
			"ErrOrderNotIssued",
			util.ErrOrderNotIssued,
			[]string{"1", "1"},
			func(s *mocks.MockOrderService) {
				order := uts.expectedOrder
				order.Issued = false
				s.EXPECT().Exists("1").Return(order, true)
			},
		},
		{
			"ErrReturnPeriodExpired",
			util.ErrReturnPeriodExpired,
			[]string{"1", "1"},
			func(s *mocks.MockOrderService) {
				order := uts.expectedOrder
				order.Issued = true
				order.StorageUntil = order.StorageUntil.AddDate(-1000, 0, 0)
				s.EXPECT().Exists("1").Return(order, true)
			},
		},
	}
	for _, tt := range tests {
		uts.T().Run(tt.name, func(t *testing.T) {
			t.Parallel()
			orderSrvc := mocks.NewMockOrderService(t)
			tt.expectMockOrderService(orderSrvc)
			validation := NewValidationService(orderSrvc, mocks.NewMockPackageService(t))

			_, err := validation.ValidateAcceptReturn(tt.input[0], tt.input[1])

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func (uts *UnitTestSuite) Test_ValidateReturnToCourier() {
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
				order := uts.expectedOrder
				order.Issued = true
				s.EXPECT().Exists("1").Return(order, true)
			},
		},
	}
	for _, tt := range tests {
		uts.T().Run(tt.name, func(t *testing.T) {
			t.Parallel()
			orderSrvc := mocks.NewMockOrderService(t)
			tt.expectMockOrderService(orderSrvc)
			validation := NewValidationService(orderSrvc, mocks.NewMockPackageService(t))

			err := validation.ValidateReturnToCourier(tt.input)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func (uts *UnitTestSuite) Test_ValidateList_ErrNotProvided() {
	uts.T().Parallel()
	_, _, err := uts.validationService.ValidateList("", "1")
	uts.EqualError(err, util.ErrOffsetNotProvided.Error())
	_, _, err = uts.validationService.ValidateList("1", "")
	uts.EqualError(err, util.ErrLimitNotProvided.Error())
}
