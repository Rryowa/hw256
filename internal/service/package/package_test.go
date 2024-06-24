package _package

import (
	"github.com/stretchr/testify/assert"
	"homework/internal/models"
	"homework/internal/util"
	"testing"
)

func Test_ValidatePackage(t *testing.T) {
	packageService := NewPackageService()
	err := packageService.ValidatePackage(models.Weight(10), "packet")
	assert.EqualError(t, err, util.ErrWeightExceeds.Error())
	err = packageService.ValidatePackage(models.Weight(30), "box")
	assert.EqualError(t, err, util.ErrWeightExceeds.Error())
}

func Test_ValidatePackage_ErrPackageTypeInvalid(t *testing.T) {
	packageService := NewPackageService()
	err := packageService.ValidatePackage(models.Weight(10), "LOL")
	assert.EqualError(t, err, util.ErrPackageTypeInvalid.Error())
}

func Test_ApplyPackageFilm(t *testing.T) {
	order := models.Order{OrderPrice: 10}
	packageService := NewPackageService()
	expectedOrder := models.Order{
		PackageType:  FilmType,
		PackagePrice: FilmPrice,
		OrderPrice:   10 + FilmPrice,
	}

	packageService.ApplyPackage(&order, "film")

	assert.Equal(t, expectedOrder, order)
}
func Test_ApplyPackagePacket(t *testing.T) {
	order := models.Order{OrderPrice: 10}
	packageService := NewPackageService()
	expectedOrder := models.Order{
		PackageType:  PacketType,
		PackagePrice: PacketPrice,
		OrderPrice:   10 + PacketPrice,
	}

	packageService.ApplyPackage(&order, "packet")

	assert.Equal(t, expectedOrder, order)
}
func Test_ApplyPackageBox(t *testing.T) {
	order := models.Order{OrderPrice: 10}
	packageService := NewPackageService()
	expectedOrder := models.Order{
		PackageType:  BoxType,
		PackagePrice: BoxPrice,
		OrderPrice:   10 + BoxPrice,
	}

	packageService.ApplyPackage(&order, "box")

	assert.Equal(t, expectedOrder, order)
}
