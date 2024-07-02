package view

import (
	"github.com/stretchr/testify/require"
	"homework/internal/models"
	"homework/internal/util"
	"testing"
)

type expectDtoStorage func(dto models.Dto) models.Dto

func Test_ValidateAcceptArgs(t *testing.T) {
	tests := []struct {
		name        string
		expectedErr error
		expectDto   expectDtoStorage
	}{
		{
			"ErrOrderIdNotProvided",
			util.ErrOrderIdNotProvided,
			func(dto models.Dto) models.Dto {
				dto.ID = ""
				return dto
			},
		},
		{
			"ErrOrderIdInvalid",
			util.ErrOrderIdInvalid,
			func(dto models.Dto) models.Dto {
				dto.ID = "some id"
				return dto
			},
		},
		{
			"ErrUserIdNotProvided",
			util.ErrUserIdNotProvided,
			func(dto models.Dto) models.Dto {
				dto.UserID = ""
				return dto
			},
		},
		{
			"ErrUserIdInvalid",
			util.ErrUserIdInvalid,
			func(dto models.Dto) models.Dto {
				dto.UserID = "some user id"
				return dto
			},
		},
		{
			"ErrWeightNotProvided",
			util.ErrWeightNotProvided,
			func(dto models.Dto) models.Dto {
				dto.Weight = ""
				return dto
			},
		},
		{
			"ErrPriceNotProvided",
			util.ErrPriceNotProvided,
			func(dto models.Dto) models.Dto {
				dto.OrderPrice = ""
				return dto
			},
		},
	}
	dto := models.Dto{
		ID:         "1",
		UserID:     "1",
		OrderPrice: "999.99",
		Weight:     "10",
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateAcceptArgs(tt.expectDto(dto))

			require.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestValidateIssueArgs_ErrUserIdNotProvided(t *testing.T) {
	t.Parallel()

	err := ValidateIssueArgs("")

	require.Equal(t, util.ErrUserIdNotProvided, err)
}

func Test_ValidateAcceptReturnArgs(t *testing.T) {
	tests := []struct {
		name        string
		expectedErr error
		id          string
		userId      string
	}{
		{
			"ErrOrderIdNotProvided",
			util.ErrOrderIdNotProvided,
			"",
			"1",
		},
		{
			"ErrOrderIdInvalid",
			util.ErrOrderIdInvalid,
			"some id",
			"1",
		},
		{
			"ErrUserIdNotProvided",
			util.ErrUserIdNotProvided,
			"1",
			"",
		},
		{
			"ErrUserIdInvalid",
			util.ErrUserIdInvalid,
			"1",
			"some user id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateAcceptReturnArgs(tt.id, tt.userId)

			require.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestValidateReturnToCourierArgs_ErrOrderIdNotProvided(t *testing.T) {
	t.Parallel()

	err := ValidateReturnToCourierArgs("")

	require.Equal(t, util.ErrOrderIdNotProvided, err)
}

func TestValidateReturnToCourierArgs_ErrOrderIdInvalid(t *testing.T) {
	t.Parallel()

	err := ValidateReturnToCourierArgs("some id")

	require.Equal(t, util.ErrOrderIdInvalid, err)
}

func TestValidateListArgs_ErrOffsetNotProvided(t *testing.T) {
	t.Parallel()

	err := ValidateListArgs("", "10")

	require.Equal(t, util.ErrOffsetNotProvided, err)
}

func TestValidateListArgs_ErrLimitNotProvided(t *testing.T) {
	t.Parallel()

	err := ValidateListArgs("0", "")

	require.Equal(t, util.ErrLimitNotProvided, err)
}
