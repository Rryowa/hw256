package service

import (
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
		}, {
			"ErrPriceNotProvided",
			util.ErrPriceNotProvided,
			func(uts models.Dto) models.Dto {
				uts.OrderPrice = ""
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

//func (uts *UnitTestSuite) Test_ValidateIssue_ErrOrderNotFound() {
//	ids := []string{"2"}
//	uts.repository.On("Exists", mock.Anything).Return(false)
//
//	_, err := uts.validationService.ValidateIssue(ids)
//
//	uts.EqualError(err, util.ErrOrderNotFound.Error())
//}
