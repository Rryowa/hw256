package service

const (
	FilmType  PackageType  = "film"
	FilmPrice PackagePrice = 1
)

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
