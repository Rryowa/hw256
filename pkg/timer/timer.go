package timer

import (
	"time"
)

type Timer interface {
	TimeNow() time.Time
}

type TimeGenerator struct{}

func (tg *TimeGenerator) TimeNow() time.Time {
	return time.Now()
}
