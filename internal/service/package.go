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

func choosePackage(weightFloat float64) Package {
	var packageType string
	defer func() {
		log.Println("1! Based on weight, package type is:", packageType)
	}()
	defer func(packageType string) {
		log.Println("2! Based on weight, package type is:", packageType)
	}(packageType)
	defer func(packageType *string) {
		log.Println("3! Based on weight, package type is:", *packageType)
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

type Package interface {
	Validate(weight float64) error
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

func newFilm() *film {
	return &film{packageType: filmType, packagePrice: filmPrice}
}

func newPacket() *packet {
	return &packet{packageType: packetType, packagePrice: packetPrice}
}

func newBox() *box {
	return &box{packageType: boxType, packagePrice: boxPrice}
}

func GetPackageType(p Package) string {
	switch p.(type) {
	case *film:
		return string(filmType)
	case *packet:
		return string(packetType)
	case *box:
		return string(boxType)
	default:
		return ""
	}
}

func GetPackagePrice(p Package) float64 {
	switch p.(type) {
	case *film:
		return float64(filmPrice)
	case *packet:
		return float64(packetPrice)
	case *box:
		return float64(boxPrice)
	default:
		return 0
	}
}

func (p *film) Validate(weight float64) error {
	if weight > 0 {
		return nil
	}
	return util.ErrWeightExceeds
}

func (p *packet) Validate(weight float64) error {
	if weight < 10 {
		return nil
	}
	return util.ErrWeightExceeds
}

func (p *box) Validate(weight float64) error {
	if weight < 30 {
		return nil
	}
	return util.ErrWeightExceeds
}
