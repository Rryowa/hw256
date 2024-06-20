package service

import (
	"homework/internal/models"
	"testing"
	"time"
)

func BenchmarkApplyPackaging(b *testing.B) {
	storageUntil, _ := time.Parse("2006-01-02", "2077-07-07")

	for i := 0; i < b.N; i++ {
		order := models.Order{
			ID:           "1",
			UserID:       "1",
			StorageUntil: storageUntil,
			Issued:       false,
			Returned:     false,
			OrderPrice:   99.99,
			Weight:       float64(i + 1),
		}
		err := applyPackaging(&order, "")
		if err != nil {
			b.Fatalf("Error applying packaging: %v", err)
		}
	}
}
