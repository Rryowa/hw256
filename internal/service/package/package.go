package _package

import (
	"homework/internal/models"
	"homework/internal/util"
)

<<<<<<< HEAD
// TODO: add LRU cache using Strategy pattern
type packageContext struct {
	strategies map[models.PackageType]PackageStrategy
}

=======
>>>>>>> dc762e3 (refactor into command pattern)
type PackageService interface {
	ValidatePackage(weight models.Weight, packageType models.PackageType) error
	ApplyPackage(order *models.Order, packageType models.PackageType)
}

<<<<<<< HEAD
type PackageStrategy interface {
	Validate(weight models.Weight) error
	Apply(order *models.Order)
}

func NewPackageService() PackageService {
	return &packageContext{
		strategies: map[models.PackageType]PackageStrategy{
			FilmType:   NewFilmPackage(),
			PacketType: NewPacketPackage(),
			BoxType:    NewBoxPackage(),
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
	//Assume that film has no weight limit
	pc.strategies[FilmType].Apply(order)
=======
type packageService struct {
	filmPackage   *FilmPackage
	packetPackage *PacketPackage
	boxPackage    *BoxPackage
}

func NewPackageService() PackageService {
	return &packageService{
		filmPackage:   NewFilmPackage(),
		packetPackage: NewPacketPackage(),
		boxPackage:    NewBoxPackage(),
	}
}

func (ps *packageService) ValidatePackage(weight models.Weight, packageType models.PackageType) error {
	switch packageType {
	case FilmType:
		if err := ps.filmPackage.Validate(weight); err != nil {
			return err
		}
	case PacketType:
		if err := ps.packetPackage.Validate(weight); err != nil {
			return err
		}
	case BoxType:
		if err := ps.boxPackage.Validate(weight); err != nil {
			return err
		}
	case "":
		return nil
	default:
		return util.ErrPackageTypeInvalid
	}
	return nil
}

func (ps *packageService) ApplyPackage(order *models.Order, packageType models.PackageType) {
	switch packageType {
	case FilmType:
		ps.filmPackage.Apply(order)
	case PacketType:
		ps.packetPackage.Apply(order)
	case BoxType:
		ps.boxPackage.Apply(order)
	case "":
		ps.filmPackage.Apply(order)
	}
>>>>>>> dc762e3 (refactor into command pattern)
}
