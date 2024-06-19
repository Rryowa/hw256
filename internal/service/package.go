package service

import (
	"homework/internal/util"
	"log"
)

type PackageType string
type PackagePrice float64

const (
	filmType   PackageType = "film"
	packetType PackageType = "packet"
	boxType    PackageType = "box"

	filmPrice   PackagePrice = 1
	packetPrice PackagePrice = 5
	boxPrice    PackagePrice = 20
)

func NewPackage(packageType string, weightFloat float64) (Package, error) {
	switch PackageType(packageType) {
	case filmType:
		return newFilm(), nil
	case packetType:
		return newPacket(), nil
	case boxType:
		return newBox(), nil
	case "":
		return choosePackage(weightFloat), nil
	default:
		return nil, util.ErrPackageTypeInvalid
	}
}

type Package interface {
	Validate(weight float64) error
	GetPrice() float64
	GetType() string
}

func newFilm() *film {
	return &film{packageType: filmType, packagePrice: filmPrice}
}
func (pkg *film) Validate(weight float64) error {
	if weight > 0 {
		return nil
	}
	return util.ErrWeightExceeds
}
func (pkg *film) GetPrice() float64 {
	return float64(pkg.packagePrice)
}
func (pkg *film) GetType() string {
	return string(pkg.packageType)
}

func newPacket() *packet {
	return &packet{packageType: packetType, packagePrice: packetPrice}
}
func (pkg *packet) Validate(weight float64) error {
	if weight < 10 {
		return nil
	}
	return util.ErrWeightExceeds
}
func (pkg *packet) GetPrice() float64 {
	return float64(pkg.packagePrice)
}
func (pkg *packet) GetType() string {
	return string(pkg.packageType)
}

func newBox() *box {
	return &box{packageType: boxType, packagePrice: boxPrice}
}
func (pkg *box) GetPrice() float64 {
	return float64(pkg.packagePrice)
}
func (pkg *box) GetType() string {
	return string(pkg.packageType)
}
func (pkg *box) Validate(weight float64) error {
	if weight < 30 {
		return nil
	}
	return util.ErrWeightExceeds
}

func choosePackage(weightFloat float64) Package {
	var packageType string
	defer func(packageType *string) {
		log.Println("Based on weight, package type is:", *packageType)
	}(&packageType)

	if weightFloat >= 30 {
		packageType = string(filmType)
		return newFilm()
	} else if weightFloat >= 10 {
		packageType = string(boxType)
		return newBox()
	} else {
		packageType = string(packetType)
		return newPacket()
	}
}

type film struct {
	packageType  PackageType
	packagePrice PackagePrice
}

type packet struct {
	packageType  PackageType
	packagePrice PackagePrice
}

type box struct {
	packageType  PackageType
	packagePrice PackagePrice
}
