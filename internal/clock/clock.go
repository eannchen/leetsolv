// Package clock implements the clock for the leetsolv application.
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
	return time.Now().UTC()
}

func (ClockImpl) Today() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

func (ClockImpl) ToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func (ClockImpl) AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// MockClock implements Clock for testing with a fixed time.
type MockClock struct {
	FixedTime time.Time
}

// NewMockClock creates a MockClock with the given fixed time.
func NewMockClock(t time.Time) *MockClock {
	return &MockClock{FixedTime: t}
}

func (m *MockClock) Now() time.Time {
	return m.FixedTime
}

func (m *MockClock) Today() time.Time {
	return time.Date(m.FixedTime.Year(), m.FixedTime.Month(), m.FixedTime.Day(), 0, 0, 0, 0, time.UTC)
}

func (m *MockClock) ToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func (m *MockClock) AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}
