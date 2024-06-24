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

	input         models.Dto
	expectedOrder models.Order
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (uts *UnitTestSuite) SetupTest() {
	uts.repository = mocks.NewMockStorage(uts.T())
	uts.packageService = mocks.NewMockPackageService(uts.T())
	uts.validationService = NewValidationService(uts.repository, uts.packageService)

	id, userId, dateStr, orderPriceStr, weightStr, pkgTypeStr :=
		"1", "1", "2077-01-01", "999.99", "10", "film"
	uts.input = models.Dto{
		ID:           id,
		UserID:       userId,
		StorageUntil: dateStr,
		OrderPrice:   orderPriceStr,
		Weight:       weightStr,
		PackageType:  pkgTypeStr,
	}

	storageUntilDate, _ := time.Parse(time.DateOnly, dateStr)
	orderPriceFloat, _ := strconv.ParseFloat(orderPriceStr, 64)
	weightFloat, _ := strconv.ParseFloat(weightStr, 64)
	uts.expectedOrder = models.Order{
		ID:           id,
		UserID:       userId,
		StorageUntil: storageUntilDate,
		OrderPrice:   models.Price(orderPriceFloat),
		Weight:       models.Weight(weightFloat),
	}
}

// TODO: Пересмотреть воркшоп и узнать за что снимут баллы!
func (uts *UnitTestSuite) Test_ValidateAccept_HappyPath() {
	uts.repository.EXPECT().Exists(uts.input.ID).Return(false)
	uts.packageService.EXPECT().ValidatePackage(uts.expectedOrder.Weight, models.PackageType(uts.input.PackageType)).Return(nil)

	order, err := uts.validationService.ValidateAccept(uts.input)

	uts.Equal(uts.expectedOrder, order)
	uts.NoError(err)
}

func Test_Validate_IsEmpty(t *testing.T) {
	assert.Equal(t, true, isArgEmpty(""))
}

func (uts *UnitTestSuite) Test_ValidateAccept() {
	tests := []struct {
		name        string
		expectedErr error
		input       models.Dto
	}{
		{
			"ErrOrderIdNotProvided",
			util.ErrOrderIdNotProvided,
			func(uts models.Dto) models.Dto {
				uts.ID = ""
				return uts
			}(uts.input),
		},
		{
			"ErrUserIdNotProvided",
			util.ErrUserIdNotProvided,
			func(uts models.Dto) models.Dto {
				uts.UserID = ""
				return uts
			}(uts.input),
		},
		{
			"ErrWeightNotProvided",
			util.ErrWeightNotProvided,
			func(uts models.Dto) models.Dto {
				uts.Weight = ""
				return uts
			}(uts.input),
		},
		{
			"ErrPriceNotProvided",
			util.ErrPriceNotProvided,
			func(uts models.Dto) models.Dto {
				uts.OrderPrice = ""
				return uts
			}(uts.input),
		},
		{
			"ErrParsingDate",
			util.ErrParsingDate,
			func(uts models.Dto) models.Dto {
				uts.StorageUntil = "1234-56-78"
				return uts
			}(uts.input),
		},
		{
			"ErrDateInvalid",
			util.ErrDateInvalid,
			func(uts models.Dto) models.Dto {
				uts.StorageUntil = "2007-07-07"
				return uts
			}(uts.input),
		},
		{
			"ErrOrderPriceInvalid",
			util.ErrOrderPriceInvalid,
			func(uts models.Dto) models.Dto {
				uts.OrderPrice = "LOL"
				return uts
			}(uts.input),
		},
		{
			"ErrOrderPriceInvalidNegative",
			util.ErrOrderPriceInvalid,
			func(uts models.Dto) models.Dto {
				uts.OrderPrice = "-1"
				return uts
			}(uts.input),
		},
		{
			"ErrWeightInvalid",
			util.ErrWeightInvalid,
			func(uts models.Dto) models.Dto {
				uts.Weight = "LOL"
				return uts
			}(uts.input),
		},
		{
			"ErrWeightInvalidNegative",
			util.ErrWeightInvalid,
			func(uts models.Dto) models.Dto {
				uts.Weight = "-1"
				return uts
			}(uts.input),
		},
	}
	for _, tt := range tests {
		uts.T().Run(tt.name, func(t *testing.T) {
			_, err := uts.validationService.ValidateAccept(tt.input)
			uts.EqualError(err, tt.expectedErr.Error())
		})
	}
}

func (uts *UnitTestSuite) Test_ValidateIssue_ErrUserIdNotProvided() {
	ids := []string{}

	_, err := uts.validationService.ValidateIssue(ids)

	uts.EqualError(err, util.ErrUserIdNotProvided.Error())
}

func (uts *UnitTestSuite) Test_ValidateIssue_ErrOrderNotFound() {
	ids := []string{"1"}
	uts.repository.EXPECT().Exists("1").Return(false)

	_, err := uts.validationService.ValidateIssue(ids)

	uts.EqualError(err, util.ErrOrderNotFound.Error())
}

func (uts *UnitTestSuite) Test_ValidateIssue_ErrOrderIssued() {
	ids := []string{"1"}
	uts.repository.EXPECT().Exists("1").Return(true)
	uts.repository.EXPECT().Get("1").Return(models.Order{Issued: true}, nil)

	_, err := uts.validationService.ValidateIssue(ids)

	uts.EqualError(err, util.ErrOrderIssued.Error())
}

func (uts *UnitTestSuite) Test_ValidateIssue_ErrOrderReturned() {
	ids := []string{"1"}
	uts.repository.EXPECT().Exists("1").Return(true)
	uts.repository.EXPECT().Get("1").Return(models.Order{Returned: true}, nil)

	_, err := uts.validationService.ValidateIssue(ids)

	uts.EqualError(err, util.ErrOrderReturned.Error())
}

func (uts *UnitTestSuite) Test_ValidateIssue_ErrOrderExpired() {
	ids := []string{"1"}
	order := uts.expectedOrder
	order.StorageUntil = order.StorageUntil.AddDate(-1000, 0, 0)
	uts.repository.EXPECT().Exists("1").Return(true)
	uts.repository.EXPECT().Get("1").Return(order, nil)

	_, err := uts.validationService.ValidateIssue(ids)

	uts.EqualError(err, util.ErrOrderExpired.Error())
}

func (uts *UnitTestSuite) Test_ValidateIssue_ErrOrdersUserDiffers() {
	ids := []string{"1", "2"}
	order1 := uts.expectedOrder
	order2 := uts.expectedOrder
	order2.UserID = "2"
	uts.repository.EXPECT().Exists("1").Return(true)
	uts.repository.EXPECT().Exists("2").Return(true)
	uts.repository.EXPECT().Get("1").Return(order1, nil)
	uts.repository.EXPECT().Get("2").Return(order2, nil)

	_, err := uts.validationService.ValidateIssue(ids)

	uts.EqualError(err, util.ErrOrdersUserDiffers.Error())
}

func (uts *UnitTestSuite) Test_ValidateAcceptReturn_ErrNotProvided() {
	_, err := uts.validationService.ValidateAcceptReturn("", "1")
	uts.EqualError(err, util.ErrOrderIdNotProvided.Error())
	_, err = uts.validationService.ValidateAcceptReturn("1", "")
	uts.EqualError(err, util.ErrUserIdNotProvided.Error())
}

func (uts *UnitTestSuite) Test_ValidateAcceptReturn_ErrOrderNotFound() {
	uts.repository.EXPECT().Exists("1").Return(false)

	_, err := uts.validationService.ValidateAcceptReturn("1", "1")

	uts.EqualError(err, util.ErrOrderNotFound.Error())
}

func (uts *UnitTestSuite) Test_ValidateAcceptReturn_ErrOrderDoesNotBelong() {
	order := uts.expectedOrder
	uts.repository.EXPECT().Exists("1").Return(true)
	uts.repository.EXPECT().Get("1").Return(order, nil)

	_, err := uts.validationService.ValidateAcceptReturn("1", "2")

	uts.EqualError(err, util.ErrOrderDoesNotBelong.Error())
}

func (uts *UnitTestSuite) Test_ValidateAcceptReturn_ErrOrderNotIssued() {
	order := uts.expectedOrder
	order.Issued = false
	uts.repository.EXPECT().Exists("1").Return(true)
	uts.repository.EXPECT().Get("1").Return(order, nil)

	_, err := uts.validationService.ValidateAcceptReturn("1", "1")

	uts.EqualError(err, util.ErrOrderNotIssued.Error())
}

func (uts *UnitTestSuite) Test_ValidateAcceptReturn_ErrReturnPeriodExpired() {
	order := uts.expectedOrder
	order.Issued = true
	order.StorageUntil = order.StorageUntil.AddDate(-1000, 0, 0)
	uts.repository.EXPECT().Exists("1").Return(true)
	uts.repository.EXPECT().Get("1").Return(order, nil)

	_, err := uts.validationService.ValidateAcceptReturn("1", "1")

	uts.EqualError(err, util.ErrReturnPeriodExpired.Error())
}

func (uts *UnitTestSuite) Test_ValidateReturnToCourier_ErrOrderIdNotProvided() {
	err := uts.validationService.ValidateReturnToCourier("")
	uts.EqualError(err, util.ErrOrderIdNotProvided.Error())
}

func (uts *UnitTestSuite) Test_ValidateReturnToCourier_ErrOrderNotFound() {
	uts.repository.EXPECT().Exists("1").Return(false)

	err := uts.validationService.ValidateReturnToCourier("1")

	uts.EqualError(err, util.ErrOrderNotFound.Error())
}

func (uts *UnitTestSuite) Test_ValidateReturnToCourier_ErrOrderIssued() {
	order := uts.expectedOrder
	order.Issued = true
	uts.repository.EXPECT().Exists("1").Return(true)
	uts.repository.EXPECT().Get("1").Return(order, nil)

	err := uts.validationService.ValidateReturnToCourier("1")

	uts.EqualError(err, util.ErrOrderIssued.Error())
}

func (uts *UnitTestSuite) Test_ValidateList_ErrNotProvided() {
	_, _, err := uts.validationService.ValidateList("", "1")
	uts.EqualError(err, util.ErrOffsetNotProvided.Error())
	_, _, err = uts.validationService.ValidateList("1", "")
	uts.EqualError(err, util.ErrLimitNotProvided.Error())
}
