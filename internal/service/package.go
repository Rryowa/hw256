package service

import (
	"errors"
	"homework/internal/util"
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

func NewPackage(packageType string) (Package, error) {
	switch PackageType(packageType) {
	case filmType:
		return &film{packageType: filmType, packagePrice: filmPrice}, nil
	case packetType:
		return &packet{packageType: packetType, packagePrice: packetPrice}, nil
	case boxType:
		return &box{packageType: boxType, packagePrice: boxPrice}, nil
	default:
		return nil, errors.New("invalid package type")
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

//func NewFilm() *Film {
//	return &Film{packageType: filmType, packagePrice: filmPrice}
//}
//
//func NewPacket() *Packet {
//	return &Packet{packageType: packetType, packagePrice: packetPrice}
//}
//
//func NewBox() *Box {
//	return &Box{packageType: boxType, packagePrice: boxPrice}
//}
//
//func GetPackageType(p Package) string {
//	switch p.(type) {
//	case *Film:
//		return string(filmType)
//	case *Packet:
//		return string(packetType)
//	case *Box:
//		return string(boxType)
//	}
//	return ""
//}
//
//func GetPackagePrice(p Package) float64 {
//	switch p.(type) {
//	case *Film:
//		return float64(filmPrice)
//	case *Packet:
//		return float64(packetPrice)
//	case *Box:
//		return float64(boxPrice)
//	}
//	return 0
//}

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
