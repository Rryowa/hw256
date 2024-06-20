package service

type PackageType string
type PackagePrice float64
type PackageWeight float64

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

func ChoosePackage(weight float64) *Package {
	if weight < float64(MaxPacketWeight) {
		return NewPackage(NewPacketPackage())
	} else if weight < float64(MaxBoxWeight) {
		return NewPackage(NewBoxPackage())
	} else {
		return NewPackage(NewFilmPackage())
	}
}
