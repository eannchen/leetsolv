package clock

import (
	"time"
)

type Clock interface {
	Now() time.Time
	Today() time.Time
	ToDate(t time.Time) time.Time
	AddDays(t time.Time, days int) time.Time
}

func NewClock() ClockImpl {
	return ClockImpl{}
}

type ClockImpl struct{}

func (ClockImpl) Now() time.Time {
	return time.Now()
}

func (ClockImpl) Today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func (ClockImpl) ToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func (ClockImpl) AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}
