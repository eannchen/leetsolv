package clock

import (
	"testing"
	"time"
)

func TestNewClock(t *testing.T) {
	clock := NewClock()
	// ClockImpl is a struct, so it can't be nil
	// Just verify it's created successfully
	_ = clock
}

func TestClockImpl_Now(t *testing.T) {
	clock := NewClock()
	now := clock.Now()

	// Now() should return a time close to the current time
	// Allow for a small delay in test execution
	if time.Since(now) > 100*time.Millisecond {
		t.Errorf("Now() returned time too far in the past: %v", now)
	}
}

func TestClockImpl_Today(t *testing.T) {
	clock := NewClock()
	today := clock.Today()

	// Today() should return today's date at midnight
	now := time.Now()
	expected := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if !today.Equal(expected) {
		t.Errorf("Today() returned %v, expected %v", today, expected)
	}

	// Should be at midnight (00:00:00)
	if today.Hour() != 0 || today.Minute() != 0 || today.Second() != 0 {
		t.Errorf("Today() should return time at midnight, got %v", today)
	}
}

func TestClockImpl_ToDate(t *testing.T) {
	clock := NewClock()

	// Test with a specific time
	testTime := time.Date(2023, 12, 25, 15, 30, 45, 123456789, time.UTC)
	dateOnly := clock.ToDate(testTime)

	expected := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	if !dateOnly.Equal(expected) {
		t.Errorf("ToDate() returned %v, expected %v", dateOnly, expected)
	}

	// Should be at midnight
	if dateOnly.Hour() != 0 || dateOnly.Minute() != 0 || dateOnly.Second() != 0 {
		t.Errorf("ToDate() should return time at midnight, got %v", dateOnly)
	}
}

func TestClockImpl_AddDays(t *testing.T) {
	clock := NewClock()

	// Test adding positive days
	startTime := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
	result := clock.AddDays(startTime, 5)
	expected := time.Date(2023, 12, 30, 15, 30, 45, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("AddDays(5) returned %v, expected %v", result, expected)
	}

	// Test adding negative days
	result = clock.AddDays(startTime, -3)
	expected = time.Date(2023, 12, 22, 15, 30, 45, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("AddDays(-3) returned %v, expected %v", result, expected)
	}

	// Test adding zero days
	result = clock.AddDays(startTime, 0)
	if !result.Equal(startTime) {
		t.Errorf("AddDays(0) returned %v, expected %v", result, startTime)
	}

	// Test month/year boundary crossing
	startTime = time.Date(2023, 12, 31, 12, 0, 0, 0, time.UTC)
	result = clock.AddDays(startTime, 1)
	expected = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("AddDays(1) across year boundary returned %v, expected %v", result, expected)
	}

	// Test leap year
	startTime = time.Date(2024, 2, 28, 12, 0, 0, 0, time.UTC)
	result = clock.AddDays(startTime, 1)
	expected = time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("AddDays(1) in leap year returned %v, expected %v", result, expected)
	}
}

func TestClockImpl_InterfaceCompliance(t *testing.T) {
	// Test that ClockImpl implements the Clock interface
	var _ Clock = NewClock()
}

func TestClockImpl_Consistency(t *testing.T) {
	clock := NewClock()

	// Test that multiple calls to Now() return different times
	// (indicating real time progression)
	now1 := clock.Now()
	time.Sleep(1 * time.Millisecond) // Small delay
	now2 := clock.Now()

	if now1.Equal(now2) {
		t.Error("Multiple calls to Now() returned identical times")
	}

	// Test that Today() returns consistent results within the same day
	today1 := clock.Today()
	today2 := clock.Today()

	if !today1.Equal(today2) {
		t.Errorf("Multiple calls to Today() returned different results: %v vs %v", today1, today2)
	}
}
