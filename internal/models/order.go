package models

import "time"

type Order struct {
	ID           string    `db:"id"`
	UserID       string    `db:"user_id"`
	StorageUntil time.Time `db:"storage_until"`
	Issued       bool      `db:"issued"`
	IssuedAt     time.Time `db:"issued_at"`
	Returned     bool      `db:"returned"`
	OrderPrice   float64   `db:"order_price"`
	Weight       float64   `db:"weight"`
	PackageType  string    `db:"package_type"`
	Hash         string    `db:"hash"`
}
