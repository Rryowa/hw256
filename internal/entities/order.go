package entities

import "time"

// Order Escape circular dependencies in file.go and orders.go
type Order struct {
	ID           string    `json:"id"`
	RecipientID  string    `json:"r_id"`
	StorageUntil time.Time `json:"storage_until"`
	Issued       bool      `json:"issued"`
	IssuedAt     time.Time `json:"issued_at"`
	Returned     bool      `json:"returned"`
	Hash         string    `json:"hash"`
}
