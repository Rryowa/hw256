package hash

import (
	"time"

	"github.com/google/uuid"
)

type Hasher interface {
	GenerateHash() string
}

type HashGenerator struct{}

// GenerateHash возвращает случайный "хэш"
func (hg *HashGenerator) GenerateHash() string {
	time.Sleep(time.Second * 5) // имитируем долгую работу
	id := uuid.New()

	return id.String()
}
