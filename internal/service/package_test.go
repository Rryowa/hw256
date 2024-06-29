package service

import (
	"github.com/stretchr/testify/assert"
	"homework/internal/models"
	_package "homework/internal/service/package"
	"homework/internal/util"
	"testing"
)

func (uts *UnitTestSuite) Test_ValidatePackage() {
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
		{
			"ErrPackageTypeInvalid",
			util.ErrPackageTypeInvalid,
			10,
			"LOL",
		},
	}
	for _, tt := range tests {
		uts.T().Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pkg := _package.NewPackageService()
			err := pkg.ValidatePackage(tt.weight, tt.pkgType)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
