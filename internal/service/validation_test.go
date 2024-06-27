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

type UnitTestSuite struct {
	suite.Suite
	repository        *mocks.MockStorage
	packageService    *mocks.MockPackageService
	validationService ValidationService

	inputDto      models.Dto
	expectedOrder models.Order
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (uts *UnitTestSuite) SetupTest() {
	uts.repository = mocks.NewMockStorage(uts.T())
	uts.packageService = mocks.NewMockPackageService(uts.T())
	uts.validationService = NewValidationService(uts.repository, uts.packageService)

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

type expectMockStorage func(s *mocks.MockStorage)

func (uts *UnitTestSuite) Test_ValidateIssue() {
	tests := []struct {
		name              string
		expectedErr       error
		input             string
		expectMockStorage expectMockStorage
	}{
		{
			"ErrUserIdNotProvided",
			util.ErrUserIdNotProvided,
			"",
			func(s *mocks.MockStorage) {
			},
		},
		{
			"ErrOrderNotFound",
			util.ErrOrderNotFound,
			"1",
			func(s *mocks.MockStorage) {
				s.EXPECT().Get("1").Return(models.Order{}, util.ErrOrderNotFound)
			},
		},
		{
			"ErrOrderIssued",
			util.ErrOrderIssued,
			"1",
			func(s *mocks.MockStorage) {
				s.EXPECT().Get("1").Return(models.Order{Issued: true}, nil)
			},
		},
		{
			"ErrOrderReturned",
			util.ErrOrderReturned,
			"1",
			func(s *mocks.MockStorage) {
				s.EXPECT().Get("1").Return(models.Order{Returned: true}, nil)
			},
		},
		{
			"ErrOrderExpired",
			util.ErrOrderExpired,
			"1",
			func(s *mocks.MockStorage) {
				order := models.Order{StorageUntil: time.Now()}
				order.StorageUntil = order.StorageUntil.AddDate(-1000, 0, 0)
				s.EXPECT().Get("1").Return(order, nil)
			},
		},
	}
	for _, tt := range tests {
		uts.T().Run(tt.name, func(t *testing.T) {
			t.Parallel()
			/*Cant use uts.repository because it is a pointer
			inside subtests will be conflict between each EXPECT()
			repository will be overwritten by other subtests
			uts.SetupTest() wont help(in parallel run)*/
			repository := mocks.NewMockStorage(t)
			tt.expectMockStorage(repository)
			validation := NewValidationService(repository, mocks.NewMockPackageService(t))

			_, err := validation.ValidateIssue(tt.input)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func (uts *UnitTestSuite) Test_ValidateAcceptReturn() {
	tests := []struct {
		name              string
		expectedErr       error
		input             []string
		expectMockStorage expectMockStorage
	}{
		{
			"ErrOrderIdNotProvided",
			util.ErrOrderIdNotProvided,
			[]string{"", "1"},
			func(s *mocks.MockStorage) {
			},
		},
		{
			"ErrUserIdNotProvided",
			util.ErrUserIdNotProvided,
			[]string{"1", ""},
			func(s *mocks.MockStorage) {
			},
		},
		{
			"ErrOrderNotFound",
			util.ErrOrderNotFound,
			[]string{"1", "1"},
			func(s *mocks.MockStorage) {
				s.EXPECT().Get("1").Return(models.Order{}, util.ErrOrderNotFound)
			},
		},
		{
			"ErrOrderDoesNotBelong",
			util.ErrOrderDoesNotBelong,
			[]string{"1", "2"},
			func(s *mocks.MockStorage) {
				order := uts.expectedOrder
				s.EXPECT().Get("1").Return(order, nil)
			},
		},
		{
			"ErrOrderNotIssued",
			util.ErrOrderNotIssued,
			[]string{"1", "1"},
			func(s *mocks.MockStorage) {
				order := uts.expectedOrder
				order.Issued = false
				s.EXPECT().Get("1").Return(order, nil)
			},
		},
		{
			"ErrReturnPeriodExpired",
			util.ErrReturnPeriodExpired,
			[]string{"1", "1"},
			func(s *mocks.MockStorage) {
				order := uts.expectedOrder
				order.Issued = true
				order.StorageUntil = order.StorageUntil.AddDate(-1000, 0, 0)
				s.EXPECT().Get("1").Return(order, nil)
			},
		},
	}
	for _, tt := range tests {
		uts.T().Run(tt.name, func(t *testing.T) {
			t.Parallel()
			/*Cant use uts.repository because it is a pointer
			inside subtests will be conflict between each EXPECT()
			repository will be overwritten by other subtests
			uts.SetupTest() wont help(in parallel run)*/
			repository := mocks.NewMockStorage(t)
			tt.expectMockStorage(repository)
			validation := NewValidationService(repository, mocks.NewMockPackageService(t))

			_, err := validation.ValidateAcceptReturn(tt.input[0], tt.input[1])

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func (uts *UnitTestSuite) Test_ValidateReturnToCourier() {
	tests := []struct {
		name              string
		expectedErr       error
		input             string
		expectMockStorage expectMockStorage
	}{
		{
			"ErrOrderIdNotProvided",
			util.ErrOrderIdNotProvided,
			"",
			func(s *mocks.MockStorage) {
			},
		},
		{
			"ErrOrderNotFound",
			util.ErrOrderNotFound,
			"1",
			func(s *mocks.MockStorage) {
				s.EXPECT().Get("1").Return(models.Order{}, util.ErrOrderNotFound)
			},
		},
		{
			"ErrOrderIssued",
			util.ErrOrderIssued,
			"1",
			func(s *mocks.MockStorage) {
				order := uts.expectedOrder
				order.Issued = true
				s.EXPECT().Get("1").Return(order, nil)
			},
		},
	}
	for _, tt := range tests {
		uts.T().Run(tt.name, func(t *testing.T) {
			t.Parallel()
			/*Cant use uts.repository because it is a pointer
			inside subtests will be conflict between each EXPECT()
			repository will be overwritten by other subtests
			uts.SetupTest() wont help(in parallel run)*/
			repository := mocks.NewMockStorage(t)
			tt.expectMockStorage(repository)
			validation := NewValidationService(repository, mocks.NewMockPackageService(t))

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
