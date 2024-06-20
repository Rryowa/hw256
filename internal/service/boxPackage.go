package service

import "homework/internal/util"

const (
	BoxType      PackageType   = "box"
	BoxPrice     PackagePrice  = 20
	MaxBoxWeight PackageWeight = 30
)

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
