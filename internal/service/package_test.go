package service

import (
	"github.com/stretchr/testify/assert"
	"homework/internal/models"
	"homework/internal/util"
	"testing"
)

func Test_ValidatePackage(t *testing.T) {
	tests := []struct {
		name        string
		expectedErr error
		weight      models.Weight
		pkgType     models.PackageType
	}{
		{
			"ErrWeightExceedsPacket",
			util.ErrWeightExceeds,
			10,
			"packet",
		},
		{
			"ErrWeightExceedsBox",
			util.ErrWeightExceeds,
			30,
			"box",
		},
		{
			"ErrPackageTypeInvalid",
			util.ErrPackageTypeInvalid,
			10,
			"LOL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pkg := NewPackageService()
			err := pkg.ValidatePackage(tt.weight, tt.pkgType)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
