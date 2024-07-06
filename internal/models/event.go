package models

import "time"

type Event struct {
	ID          int       `db:"id"`
	Request     string    `db:"request"`
	MethodName  string    `db:"method_name"`
	Acquired    bool      `db:"acquired"`
	Processed   bool      `db:"processed"`
	AcquiredAt  time.Time `db:"acquired_at"`
	ProcessedAt time.Time `db:"processed_at"`
}