package _package

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"homework/internal/models"
	"homework/internal/util"
	"testing"
)

type UnitTestSuite struct {
	suite.Suite
	packageService PackageService

	weight  models.Weight
	pkgType models.PackageType
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (uts *UnitTestSuite) SetupTest() {
	uts.packageService = NewPackageService()
}

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

			err := uts.packageService.ValidatePackage(tt.weight, tt.pkgType)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
