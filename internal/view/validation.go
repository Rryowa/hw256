package view

import (
	"homework/internal/models"
	"homework/internal/util"
	"strconv"
)

func ValidateAcceptArgs(dto models.Dto) error {
	if len(dto.ID) == 0 {
		return util.ErrOrderIdNotProvided
	}
	if _, err := strconv.Atoi(dto.ID); err != nil {
		return util.ErrOrderIdInvalid
	}
	if len(dto.UserID) == 0 {
		return util.ErrUserIdNotProvided
	}
	if _, err := strconv.Atoi(dto.UserID); err != nil {
		return util.ErrUserIdInvalid
	}
	if len(dto.Weight) == 0 {
		return util.ErrWeightNotProvided
	}
	if len(dto.OrderPrice) == 0 {
		return util.ErrPriceNotProvided
	}

	return nil
}

func ValidateIssueArgs(idsStr string) error {
	if len(idsStr) == 0 {
		return util.ErrUserIdNotProvided
	}
	return nil
}

func ValidateAcceptReturnArgs(id, userId string) error {
	if len(id) == 0 {
		return util.ErrOrderIdNotProvided
	}
	if _, err := strconv.Atoi(id); err != nil {
		return util.ErrOrderIdInvalid
	}
	if len(userId) == 0 {
		return util.ErrUserIdNotProvided
	}
	if _, err := strconv.Atoi(userId); err != nil {
		return util.ErrUserIdInvalid
	}
	return nil
}

func ValidateReturnToCourierArgs(id string) error {
	if len(id) == 0 {
		return util.ErrOrderIdNotProvided
	}
	if _, err := strconv.Atoi(id); err != nil {
		return util.ErrOrderIdInvalid
	}

	return nil
}

func ValidateListArgs(offset, limit string) error {
	if len(offset) == 0 {
		return util.ErrOffsetNotProvided
	}
	if len(limit) == 0 {
		return util.ErrLimitNotProvided
	}
	return nil
}
