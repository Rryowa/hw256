package models

import "time"

type Event struct {
	ID          int         `db:"id"`
	MethodName  string      `db:"method_name"`
	Request     string      `db:"request"`
	Status      EventStatus `db:"status"`
	RequestedAt time.Time   `db:"requested_at"`
}

type EventStatus string

const (
	EventStatusRequested EventStatus = "requested"
)