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

func Test_ValidateAccept(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repository := mocks.NewMockStorage(t)
			packageService := mocks.NewMockPackageService(t)
			validationService := NewValidationService(repository, packageService)

			input := models.Dto{
				ID:           "1",
				UserID:       "1",
				StorageUntil: "2077-01-01",
				OrderPrice:   "999.99",
				Weight:       "10",
				PackageType:  "film",
			}
			_, err := validationService.ValidateAccept(tt.expectDto(input))

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

type expectMockStorage func(s *mocks.MockStorage)

func Test_ValidateIssue(t *testing.T) {
	tests := []struct {
		name              string
		expectedErr       error
		input             []string
		expectMockStorage expectMockStorage
	}{
		{
			"ErrUserIdNotProvided",
			util.ErrUserIdNotProvided,
			[]string{},
			func(s *mocks.MockStorage) {
			},
		},
		{
			"ErrOrderNotFound",
			util.ErrOrderNotFound,
			[]string{"1"},
			func(s *mocks.MockStorage) {
				s.EXPECT().Exists("1").Return(false)
			},
		},
		{
			"ErrOrderIssued",
			util.ErrOrderIssued,
			[]string{"1"},
			func(s *mocks.MockStorage) {
				s.EXPECT().Exists("1").Return(true)
				s.EXPECT().Get("1").Return(models.Order{Issued: true}, nil)
			},
		},
		{
			"ErrOrderReturned",
			util.ErrOrderReturned,
			[]string{"1"},
			func(s *mocks.MockStorage) {
				s.EXPECT().Exists("1").Return(true)
				s.EXPECT().Get("1").Return(models.Order{Returned: true}, nil)
			},
		},
		{
			"ErrOrderExpired",
			util.ErrOrderExpired,
			[]string{"1"},
			func(s *mocks.MockStorage) {
				order := models.Order{StorageUntil: time.Now()}
				order.StorageUntil = order.StorageUntil.AddDate(-1000, 0, 0)
				s.EXPECT().Exists("1").Return(true)
				s.EXPECT().Get("1").Return(order, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := mocks.NewMockStorage(t)
			tt.expectMockStorage(s)
			validationService := NewValidationService(s, mocks.NewMockPackageService(t))

			_, err := validationService.ValidateIssue(tt.input)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
