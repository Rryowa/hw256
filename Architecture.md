### Архитектурный паттерн MVC
• Модель предоставляет данные и реагирует на команды контроллера, изменяя своё состояние.  
• Представление отвечает за отображение данных модели пользователю, реагируя на изменение модели.  
• Контроллер(Сервис) интерпретирует действия пользователя, оповещая модель о необходимости изменений.

### Polymorphism

```go
func applyPackaging(order *models.Order, packageType string) error {
<<<<<<< Architecture.md
    var pkg PackageInterface
    
    switch PackageType(packageType) {
    case FilmType:
        pkg = NewFilmPackage()
    case PacketType:
        pkg = NewPacketPackage()
    case BoxType:
        pkg = NewBoxPackage()
    case "":
        pkg = ChoosePackage(order.Weight)
    default:
        return util.ErrPackageTypeInvalid
    }

    p := NewPackage(pkg)

    if err := p.Validate(order.Weight); err != nil {
        return err
    }
    
    //Apply packaging and calculate order price
    order.PackageType = p.GetType()
    order.PackagePrice = p.GetPrice()
    order.OrderPrice += p.GetPrice()
    
    return nil
=======
	var pkg PackageInterface

	switch PackageType(packageType) {
	case FilmType:
		pkg = NewFilmPackage()
	case PacketType:
		pkg = NewPacketPackage()
	case BoxType:
		pkg = NewBoxPackage()
	case "":
		pkg = ChoosePackage(order.Weight)
	default:
		return util.ErrPackageTypeInvalid
	}

	p := NewPackage(pkg)

	if err := p.Validate(order.Weight); err != nil {
		return err
	}

	//Apply packaging and calculate order price
	order.PackageType = p.GetType()
	order.PackagePrice = p.GetPrice()
	order.OrderPrice += p.GetPrice()

	return nil
>>>>>>> Architecture.md
}
```

### Template behavioral pattern
https://github.com/AlexanderGrom/go-patterns/blob/master/Behavioral/TemplateMethod/
```go
// PackageInterface provides an interface to validate different packages
type PackageInterface interface {
<<<<<<< Architecture.md
    ValidatePackage(weight float64) error
    GetType() string
    GetPrice() float64
=======
	ValidatePackage(weight float64) error
	GetType() string
	GetPrice() float64
>>>>>>> Architecture.md
}

// Package implements a Template method
type Package struct {
<<<<<<< Architecture.md
    PackageInterface
=======
	PackageInterface
>>>>>>> Architecture.md
}

// Validate is the Template Method.
func (p *Package) Validate(weight float64) error {
<<<<<<< Architecture.md
    return p.ValidatePackage(weight)
=======
	return p.ValidatePackage(weight)
>>>>>>> Architecture.md
}

// NewPackage is the Package constructor.
func NewPackage(p PackageInterface) *Package {
<<<<<<< Architecture.md
    return &Package{p}
=======
	return &Package{p}
>>>>>>> Architecture.md
}
```
```go
package service
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
```

## Используемые стандарты описания архитектуры
Для описания архитектуры я использовал *диаграмму классов* и  
*диаграмму последовательности* стандарта UML,

### Почему UML
UML я выбрал, потому что C4 Model непригодна для настолько маленького проекта,  
2 и 3 слои выглядят идентично из-за отсутствия внешнего сервиса.


### Почему диаграмма последовательности
Я выбрал диаграмму последовательности, так как с её помощью можно понять, как пользователь  
и программа взаимодействуют от начал до конца.  
Некоторого рода общая картина взаимодействия.


### Почему диаграмма классов
Я выбрал диаграмму классов, потому что таким образом можно углубиться в то как работает программа,  
при этом не углубляясь в технические детали реализации методов интерфейса.
Нечто среднее между общей картиной и точной реализацией каждого интерфейса, если коллега захочет  
углубиться в детали реализации, то они будут доступны в самом коде проекта.

### Почему в диаграмме классов есть пакеты(namespace)?
Потому что пакеты в гошке это круто и я думаю, это не усложняет процесс понимания, а наоборот помогает.  
Но я могу и ошибаться...

### Таким образом, коллега сможет при анализе диаграммы классов держать в голове общую картину программы и при этом разбираться в нужном пакете или классе(структуре) если у него возникнут вопросы!