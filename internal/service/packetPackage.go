package service

import "homework/internal/util"

const (
	PacketType      PackageType   = "packet"
	PacketPrice     PackagePrice  = 5
	MaxPacketWeight PackageWeight = 10
)

type PacketPackage struct {
}

func NewPacketPackage() *PacketPackage {
	return &PacketPackage{}
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
