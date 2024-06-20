package service

import (
	"homework/internal/util"
)

type PackageType string
type PackagePrice float64
type PackageWeight float64

const (
	FilmType  PackageType  = "film"
	FilmPrice PackagePrice = 1
)
const (
	PacketType      PackageType   = "packet"
	PacketPrice     PackagePrice  = 5
	MaxPacketWeight PackageWeight = 10
)
const (
	BoxType      PackageType   = "box"
	BoxPrice     PackagePrice  = 20
	MaxBoxWeight PackageWeight = 30
)

// TODO: Создаем заказ и применяем к нему упаковку (темплейт метод)

// PackageInterface provides an interface to validate different packages
type PackageInterface interface {
	ValidatePackage(weight float64) error
	GetType() string
	GetPrice() float64
}

// Package implements a Template method
type Package struct {
	PackageInterface
}

// Validate is the Template Method.
func (p *Package) Validate(weight float64) error {
	return p.ValidatePackage(weight)
}

// NewPackage is the Package constructor.
func NewPackage(p PackageInterface) *Package {
	return &Package{p}
}

// FilmPackage implements ValidatePackage
type FilmPackage struct {
}

func NewFilmPackage() *FilmPackage {
	return &FilmPackage{}
}

// ValidatePackage provides validation
func (p *FilmPackage) ValidatePackage(weight float64) error {
	return nil
}
func (p *FilmPackage) GetType() string {
	return string(FilmType)
}
func (p *FilmPackage) GetPrice() float64 {
	return float64(FilmPrice)
}

type PacketPackage struct {
}

func NewPacketPackage() *PacketPackage {
	return &PacketPackage{}
}

func ChoosePackage(weight float64) *Package {
	if weight < float64(MaxPacketWeight) {
		return NewPackage(NewPacketPackage())
	} else if weight < float64(MaxBoxWeight) {
		return NewPackage(NewBoxPackage())
	} else {
		return NewPackage(NewFilmPackage())
	}
}

func (p *PacketPackage) ValidatePackage(weight float64) error {
	if weight < float64(MaxPacketWeight) {
		return nil
	}
	return util.ErrWeightExceeds
}
func (p *PacketPackage) GetType() string {
	return string(PacketType)
}
func (p *PacketPackage) GetPrice() float64 {
	return float64(PacketPrice)
}

type BoxPackage struct {
}

func NewBoxPackage() *BoxPackage {
	return &BoxPackage{}
}

func (p *BoxPackage) ValidatePackage(weight float64) error {
	if weight < float64(MaxBoxWeight) {
		return nil
	}
	return util.ErrWeightExceeds
}
func (p *BoxPackage) GetType() string {
	return string(BoxType)
}
func (p *BoxPackage) GetPrice() float64 {
	return float64(BoxPrice)
}
