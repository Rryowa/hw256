package models

import "time"

// Order Escape circular dependencies in file.go and storage.go
type Order struct {
	ID           string    `db:"id"`
	UserID       string    `db:"user_id"`
	StorageUntil time.Time `db:"storage_until"`
	Issued       bool      `db:"issued"`
	IssuedAt     time.Time `db:"issued_at"`
	Returned     bool      `db:"returned"`
	Hash         string    `db:"hash"`
}
