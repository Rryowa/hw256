package service

import (
	"homework/internal/models"
	"homework/internal/service/package"
	"homework/internal/util"
)

type packageContext struct {
	strategies map[models.PackageType]PackageStrategy
}

type PackageService interface {
	ValidatePackage(weight models.Weight, packageType models.PackageType) error
	ApplyPackage(order *models.Order, packageType models.PackageType)
}

type PackageStrategy interface {
	Validate(weight models.Weight) error
	Apply(order *models.Order)
}

func NewPackageService() PackageService {
	return &packageContext{
		strategies: map[models.PackageType]PackageStrategy{
			_package.FilmType:   _package.NewFilmPackage(),
			_package.PacketType: _package.NewPacketPackage(),
			_package.BoxType:    _package.NewBoxPackage(),
		},
	}
}

func (pc *packageContext) ValidatePackage(weight models.Weight, packageType models.PackageType) error {
	if strategy, ok := pc.strategies[packageType]; ok {
		return strategy.Validate(weight)
	}
	return util.ErrPackageTypeInvalid
}

func (pc *packageContext) ApplyPackage(order *models.Order, packageType models.PackageType) {
	if strategy, ok := pc.strategies[packageType]; ok {
		strategy.Apply(order)
		return
	}
}
