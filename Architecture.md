### Архитектурный паттерн MVC
• Модель предоставляет данные и реагирует на команды контроллера, изменяя своё состояние.  
• Представление отвечает за отображение данных модели пользователю, реагируя на изменение модели.  
• Контроллер(Сервис) интерпретирует действия пользователя, оповещая модель о необходимости изменений.

### Polymorphism
```go
// /service/order.go
func ApplyPackaging(weightFloat float64, packageType string) (Package, error) {
	pkg, err := NewPackage(packageType, weightFloat)
	if err != nil {
		return nil, err
	}
	if err := pkg.Validate(weightFloat); err != nil {
		return nil, err
	}

	return pkg, nil
}
```

### Factory creational pattern https://github.com/AlexanderGrom/go-patterns/blob/master/Creational/FactoryMethod
```go
// /service/package.go
type Package interface {
	Validate(weight float64) error
	GetPrice() float64
	GetType() string
}

func newFilm() *film {
	return &film{packageType: filmType, packagePrice: filmPrice}
}
func (pkg *film) Validate(weight float64) error {
	if weight > 0 {
		return nil
	}
	return util.ErrWeightExceeds
}
func (pkg *film) GetPrice() float64 {
	return float64(pkg.packagePrice)
}
func (pkg *film) GetType() string {
	return string(pkg.packageType)
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
