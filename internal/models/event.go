package models

import "time"

type Event struct {
	ID          int         `db:"id"`
	Request     string      `db:"request"`
	MethodName  string      `db:"method_name"`
	Status      EventStatus `db:"status"`
	AcquiredAt  time.Time   `db:"acquired_at"`
	ProcessedAt time.Time   `db:"processed_at"`
}

type EventStatus string

const (
	EventStatusAcquired  EventStatus = "acquired"
	EventStatusProcessed EventStatus = "processed"
)