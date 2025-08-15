package core

import (
	"testing"
	"time"

	"leetsolv/config"
)

// MockClock implements clock.Clock for testing
type MockClock struct {
	currentTime time.Time
}

func NewMockClock(t time.Time) *MockClock {
	return &MockClock{currentTime: t}
}

func (m *MockClock) Now() time.Time {
	return m.currentTime
}

func (m *MockClock) Today() time.Time {
	return time.Date(m.currentTime.Year(), m.currentTime.Month(), m.currentTime.Day(), 0, 0, 0, 0, m.currentTime.Location())
}

func (m *MockClock) ToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func (m *MockClock) AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

func TestNewSM2Scheduler(t *testing.T) {
	mockClock := NewMockClock(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
	_, cfg := config.MockEnv(t)
	scheduler := NewSM2Scheduler(cfg, mockClock)

	// Test that scheduler is properly initialized
	if scheduler.Clock != mockClock {
		t.Errorf("Expected Clock to be set to mockClock, got %v", scheduler.Clock)
	}

	// Test interval settings
	expectedMaxInterval := 90
	if scheduler.maxInterval != expectedMaxInterval {
		t.Errorf("Expected maxInterval to be %d, got %d", expectedMaxInterval, scheduler.maxInterval)
	}

	// Test base intervals
	expectedBaseIntervals := map[Importance]int{
		LowImportance:      8,
		MediumImportance:   6,
		HighImportance:     5,
		CriticalImportance: 4,
	}
	for importance, expected := range expectedBaseIntervals {
		if scheduler.baseIntervals[importance] != expected {
			t.Errorf("Expected baseInterval for %v to be %d, got %d", importance, expected, scheduler.baseIntervals[importance])
		}
	}

	// Test memory multipliers
	expectedMemoryMultipliers := map[MemoryUse]float64{
		MemoryReasoned: 1.00,
		MemoryPartial:  1.10,
		MemoryFull:     1.25,
	}
	for memory, expected := range expectedMemoryMultipliers {
		if scheduler.memoryMultipliers[memory] != expected {
			t.Errorf("Expected memoryMultiplier for %v to be %f, got %f", memory, expected, scheduler.memoryMultipliers[memory])
		}
	}

	// Test ease factor bounds
	if scheduler.minEaseFactor != 1.3 {
		t.Errorf("Expected minEaseFactor to be 1.3, got %f", scheduler.minEaseFactor)
	}
	if scheduler.maxEaseFactor != 2.6 {
		t.Errorf("Expected maxEaseFactor to be 2.6, got %f", scheduler.maxEaseFactor)
	}
}

func TestScheduleNewQuestion(t *testing.T) {
	mockClock := NewMockClock(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
	_, cfg := config.MockEnv(t)
	scheduler := NewSM2Scheduler(cfg, mockClock)

	tests := []struct {
		name     string
		question *Question
		memory   MemoryUse
		check    func(*Question) error
	}{
		{
			name: "Low importance question with MemoryReasoned",
			question: &Question{
				ID:          1,
				Importance:  LowImportance,
				Familiarity: Medium,
			},
			memory: MemoryReasoned,
			check: func(q *Question) error {
				if q.EaseFactor != 2.0 {
					t.Errorf("Expected EaseFactor to be 2.0, got %f", q.EaseFactor)
				}
				if q.ReviewCount != 1 {
					t.Errorf("Expected ReviewCount to be 1, got %d", q.ReviewCount)
				}
				if !q.LastReviewed.Equal(mockClock.Today()) {
					t.Errorf("Expected LastReviewed to be today, got %v", q.LastReviewed)
				}
				// Base interval for LowImportance is 8, Medium familiarity adds 2 days
				// Total: (8 + 2) × 1.00 = 10 days
				// With randomization (-1 to +2 days), expect 9-12 days
				expectedMin := mockClock.AddDays(mockClock.Today(), 9)
				expectedMax := mockClock.AddDays(mockClock.Today(), 12)
				if q.NextReview.Before(expectedMin) || q.NextReview.After(expectedMax) {
					t.Errorf("Expected NextReview to be between %v and %v, got %v", expectedMin, expectedMax, q.NextReview)
				}
				return nil
			},
		},
		{
			name: "Critical importance question with MemoryFull",
			question: &Question{
				ID:          2,
				Importance:  CriticalImportance,
				Familiarity: VeryEasy,
			},
			memory: MemoryFull,
			check: func(q *Question) error {
				if q.EaseFactor != 1.7 {
					t.Errorf("Expected EaseFactor to be 1.7, got %f", q.EaseFactor)
				}
				// Base interval for CriticalImportance is 4, VeryEasy adds 7 days, MemoryFull multiplies by 1.25
				// Expected: (4+7) * 1.25 = 13.75 ≈ 14 days, with randomization (-1 to +2): 13-16 days
				expectedMin := mockClock.AddDays(mockClock.Today(), 13)
				expectedMax := mockClock.AddDays(mockClock.Today(), 16)
				if q.NextReview.Before(expectedMin) || q.NextReview.After(expectedMax) {
					t.Errorf("Expected NextReview to be between %v and %v, got %v", expectedMin, expectedMax, q.NextReview)
				}
				return nil
			},
		},
		{
			name: "High importance question with MemoryPartial",
			question: &Question{
				ID:          3,
				Importance:  HighImportance,
				Familiarity: Easy,
			},
			memory: MemoryPartial,
			check: func(q *Question) error {
				if q.EaseFactor != 1.8 {
					t.Errorf("Expected EaseFactor to be 1.8, got %f", q.EaseFactor)
				}
				// Base interval for HighImportance is 5, Easy adds 5 days, MemoryPartial multiplies by 1.10
				// Expected: (5+5) * 1.10 = 11 days, with randomization (-1 to +2): 10-13 days
				expectedMin := mockClock.AddDays(mockClock.Today(), 10)
				expectedMax := mockClock.AddDays(mockClock.Today(), 13)
				if q.NextReview.Before(expectedMin) || q.NextReview.After(expectedMax) {
					t.Errorf("Expected NextReview to be between %v and %v, got %v", expectedMin, expectedMax, q.NextReview)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scheduler.ScheduleNewQuestion(tt.question, tt.memory)
			if result != tt.question {
				t.Errorf("Expected ScheduleNewQuestion to return the same question pointer")
			}
			if err := tt.check(tt.question); err != nil {
				t.Errorf("Check failed: %v", err)
			}
		})
	}
}

func TestSchedule(t *testing.T) {
	mockClock := NewMockClock(time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC))
	_, cfg := config.MockEnv(t)
	scheduler := NewSM2Scheduler(cfg, mockClock)

	tests := []struct {
		name     string
		question *Question
		memory   MemoryUse
		check    func(*Question) error
	}{
		{
			name: "Question with VeryHard familiarity should reset",
			question: &Question{
				ID:           1,
				Importance:   MediumImportance,
				Familiarity:  VeryHard,
				ReviewCount:  5,
				EaseFactor:   2.0,
				LastReviewed: time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC),
				NextReview:   time.Date(2024, 1, 12, 12, 0, 0, 0, time.UTC),
			},
			memory: MemoryReasoned,
			check: func(q *Question) error {
				if q.ReviewCount != 6 {
					t.Errorf("Expected ReviewCount to be 6, got %d", q.ReviewCount)
				}
				if !q.LastReviewed.Equal(mockClock.Today()) {
					t.Errorf("Expected LastReviewed to be today, got %v", q.LastReviewed)
				}
				// Should reset to base interval for MediumImportance (6 days)
				// With randomization (-1 to +2 days), expect 5-8 days
				expectedMin := mockClock.AddDays(mockClock.Today(), 5)
				expectedMax := mockClock.AddDays(mockClock.Today(), 8)
				if q.NextReview.Before(expectedMin) || q.NextReview.After(expectedMax) {
					t.Errorf("Expected NextReview to be between %v and %v, got %v", expectedMin, expectedMax, q.NextReview)
				}
				return nil
			},
		},
		{
			name: "Normal scheduling with growth",
			question: &Question{
				ID:           2,
				Importance:   HighImportance,
				Familiarity:  Medium,
				ReviewCount:  3,
				EaseFactor:   1.8,
				LastReviewed: time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC),
				NextReview:   time.Date(2024, 1, 12, 12, 0, 0, 0, time.UTC),
			},
			memory: MemoryReasoned,
			check: func(q *Question) error {
				if q.ReviewCount != 4 {
					t.Errorf("Expected ReviewCount to be 4, got %d", q.ReviewCount)
				}
				if !q.LastReviewed.Equal(mockClock.Today()) {
					t.Errorf("Expected LastReviewed to be today, got %v", q.LastReviewed)
				}
				// Previous interval was 2 days, EaseFactor 1.8, MemoryReasoned 1.0
				// Expected: 2 * 1.8 * 1.0 = 3.6 ≈ 4 days
				// With randomization (-1 to +2 days), expect 3-6 days
				expectedMin := mockClock.AddDays(mockClock.Today(), 3)
				expectedMax := mockClock.AddDays(mockClock.Today(), 6)
				if q.NextReview.Before(expectedMin) || q.NextReview.After(expectedMax) {
					t.Errorf("Expected NextReview to be between %v and %v, got %v", expectedMin, expectedMax, q.NextReview)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheduler.Schedule(tt.question, tt.memory)
			if err := tt.check(tt.question); err != nil {
				t.Errorf("Check failed: %v", err)
			}
		})
	}
}

func TestSetEaseFactor(t *testing.T) {
	mockClock := NewMockClock(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
	_, cfg := config.MockEnv(t)
	scheduler := NewSM2Scheduler(cfg, mockClock)

	tests := []struct {
		name     string
		question *Question
		memory   MemoryUse
		expected float64
	}{
		{
			name: "Low importance with Easy familiarity and MemoryReasoned",
			question: &Question{
				Importance:  LowImportance,
				Familiarity: Easy,
				EaseFactor:  2.0,
				ReviewCount: 1,
			},
			memory:   MemoryReasoned,
			expected: 2.0 + 0.15 + 0.05 + 0.00, // base + importance bonus + familiarity penalty + memory penalty
		},
		{
			name: "Critical importance with VeryHard familiarity and MemoryPartial",
			question: &Question{
				Importance:  CriticalImportance,
				Familiarity: VeryHard,
				EaseFactor:  1.7,
				ReviewCount: 1,
			},
			memory:   MemoryPartial,
			expected: 1.7 + 0.03 + (-0.40) + (-0.02), // base + importance bonus + familiarity penalty + memory penalty
		},
		{
			name: "High importance with Medium familiarity and MemoryFull, high review count",
			question: &Question{
				Importance:  HighImportance,
				Familiarity: Medium,
				EaseFactor:  1.8,
				ReviewCount: 5, // >= 3, should get stability bonus
			},
			memory:   MemoryFull,
			expected: 1.8 + 0.05 + (-0.10) + (-0.05) + (0.05 * 0.5), // base + bonus + penalties + stability bonus
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalEaseFactor := tt.question.EaseFactor
			scheduler.setEaseFactor(tt.question, tt.memory)

			// Check bounds
			if tt.question.EaseFactor < scheduler.minEaseFactor {
				t.Errorf("EaseFactor %f is below minimum %f", tt.question.EaseFactor, scheduler.minEaseFactor)
			}
			if tt.question.EaseFactor > scheduler.maxEaseFactor {
				t.Errorf("EaseFactor %f is above maximum %f", tt.question.EaseFactor, scheduler.maxEaseFactor)
			}

			// Check that ease factor changed
			if tt.question.EaseFactor == originalEaseFactor {
				t.Errorf("Expected EaseFactor to change from %f, but it remained the same", originalEaseFactor)
			}
		})
	}
}

func TestSetEaseFactorOverduePenalty(t *testing.T) {
	mockClock := NewMockClock(time.Date(2024, 1, 20, 12, 0, 0, 0, time.UTC))
	_, cfg := config.MockEnv(t)
	// Enable overdue penalty for this test
	cfg.OverduePenalty = true
	scheduler := NewSM2Scheduler(cfg, mockClock)

	tests := []struct {
		name           string
		question       *Question
		overdueLimit   int
		expectedChange bool
	}{
		{
			name: "Question overdue but within limit",
			question: &Question{
				Importance:  MediumImportance,
				Familiarity: Medium,
				EaseFactor:  1.8,
				NextReview:  time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC), // 5 days overdue
			},
			overdueLimit:   7,
			expectedChange: false, // Within limit, no penalty
		},
		{
			name: "Question overdue beyond limit",
			question: &Question{
				Importance:  HighImportance,
				Familiarity: Hard,
				EaseFactor:  1.8,
				NextReview:  time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC), // 10 days overdue
			},
			overdueLimit:   7,
			expectedChange: true, // Beyond limit, should get penalty
		},
		{
			name: "Low importance question should not get penalty",
			question: &Question{
				Importance:  LowImportance,
				Familiarity: Medium,
				EaseFactor:  1.8,
				NextReview:  time.Date(2024, 1, 5, 12, 0, 0, 0, time.UTC), // 15 days overdue
			},
			overdueLimit:   7,
			expectedChange: false, // Low importance, no penalty
		},
		{
			name: "VeryEasy question should not get penalty",
			question: &Question{
				Importance:  MediumImportance,
				Familiarity: VeryEasy,
				EaseFactor:  1.8,
				NextReview:  time.Date(2024, 1, 5, 12, 0, 0, 0, time.UTC), // 15 days overdue
			},
			overdueLimit:   7,
			expectedChange: false, // VeryEasy, no penalty
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the overdue limit for this test
			scheduler.cfg.OverdueLimit = tt.overdueLimit
			
			originalEaseFactor := tt.question.EaseFactor
			scheduler.setEaseFactorOverduePenalty(tt.question)

			if tt.expectedChange {
				if tt.question.EaseFactor >= originalEaseFactor {
					t.Errorf("Expected EaseFactor to decrease due to overdue penalty, but it didn't change or increased")
				}
			} else {
				if tt.question.EaseFactor != originalEaseFactor {
					t.Errorf("Expected EaseFactor to remain unchanged, but it changed from %f to %f", originalEaseFactor, tt.question.EaseFactor)
				}
			}
		})
	}
}

func TestSetNextReview(t *testing.T) {
	mockClock := NewMockClock(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
	_, cfg := config.MockEnv(t)
	scheduler := NewSM2Scheduler(cfg, mockClock)

	tests := []struct {
		name         string
		intervalDays int
		checkBounds  bool
		expectedMin  int
		expectedMax  int
	}{
		{
			name:         "Normal interval",
			intervalDays: 10,
			checkBounds:  true,
			expectedMin:  9,  // 10 - 1 (randomization)
			expectedMax:  12, // 10 + 2 (randomization)
		},
		{
			name:         "Minimum interval",
			intervalDays: 1,
			checkBounds:  true,
			expectedMin:  1, // Should not go below 1
			expectedMax:  3, // 1 + 2 (randomization)
		},
		{
			name:         "Maximum interval",
			intervalDays: 100,
			checkBounds:  true,
			expectedMin:  89, // 90 - 1 (randomization, capped at maxInterval)
			expectedMax:  90, // Capped at maxInterval
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			question := &Question{}
			scheduler.setNextReview(question, mockClock.Today(), tt.intervalDays)

			if tt.checkBounds {
				daysDiff := int(question.NextReview.Sub(mockClock.Today()).Hours() / 24)
				if daysDiff < tt.expectedMin || daysDiff > tt.expectedMax {
					t.Errorf("Expected interval to be between %d and %d days, got %d", tt.expectedMin, tt.expectedMax, daysDiff)
				}
			}

			// Should always be in the future
			if question.NextReview.Before(mockClock.Today()) {
				t.Errorf("NextReview %v is before today %v", question.NextReview, mockClock.Today())
			}
		})
	}
}

func TestCalculatePriorityScore(t *testing.T) {
	mockClock := NewMockClock(time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC))
	_, cfg := config.MockEnv(t)
	scheduler := NewSM2Scheduler(cfg, mockClock)

	tests := []struct {
		name     string
		question *Question
		check    func(float64) error
	}{
		{
			name: "High priority question",
			question: &Question{
				Importance:  CriticalImportance,
				Familiarity: VeryHard,
				ReviewCount: 1,
				EaseFactor:  1.7,
				NextReview:  time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC), // 5 days overdue
			},
			check: func(score float64) error {
				// Should be high priority due to CriticalImportance, VeryHard familiarity, and overdue
				if score < 10 {
					t.Errorf("Expected high priority score, got %f", score)
				}
				return nil
			},
		},
		{
			name: "Low priority question",
			question: &Question{
				Importance:  LowImportance,
				Familiarity: VeryEasy,
				ReviewCount: 10,
				EaseFactor:  2.5,
				NextReview:  time.Date(2024, 1, 20, 12, 0, 0, 0, time.UTC), // 5 days in future
			},
			check: func(score float64) error {
				// Should be low priority due to LowImportance, VeryEasy familiarity, and high review count
				if score > -5 {
					t.Errorf("Expected low priority score, got %f", score)
				}
				return nil
			},
		},
		{
			name: "Medium priority question",
			question: &Question{
				Importance:  MediumImportance,
				Familiarity: Medium,
				ReviewCount: 5,
				EaseFactor:  1.9,
				NextReview:  time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC), // Due today
			},
			check: func(score float64) error {
				// Should be medium priority
				if score < -10 || score > 10 {
					t.Errorf("Expected medium priority score, got %f", score)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := scheduler.CalculatePriorityScore(tt.question)
			if err := tt.check(score); err != nil {
				t.Errorf("Check failed: %v", err)
			}
		})
	}
}

func TestSchedulerInterface(t *testing.T) {
	mockClock := NewMockClock(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
	_, cfg := config.MockEnv(t)
	var scheduler Scheduler = NewSM2Scheduler(cfg, mockClock)

	// Test that we can call interface methods
	question := &Question{
		ID:          1,
		Importance:  MediumImportance,
		Familiarity: Medium,
	}

	// Test ScheduleNewQuestion
	result := scheduler.ScheduleNewQuestion(question, MemoryReasoned)
	if result == nil {
		t.Error("Expected ScheduleNewQuestion to return a question")
	}

	// Test Schedule
	scheduler.Schedule(question, MemoryPartial)

	// Test CalculatePriorityScore
	score := scheduler.CalculatePriorityScore(question)
	if score < 0 {
		t.Error("Expected priority score to be non-negative")
	}
}

func TestEdgeCases(t *testing.T) {
	mockClock := NewMockClock(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
	_, cfg := config.MockEnv(t)
	scheduler := NewSM2Scheduler(cfg, mockClock)

	t.Run("Question with zero interval", func(t *testing.T) {
		question := &Question{
			Importance:  MediumImportance,
			Familiarity: Medium,
		}
		scheduler.setNextReview(question, mockClock.Today(), 0)

		// Should default to 1 day minimum, with randomization could be 1-3 days
		expectedMin := mockClock.AddDays(mockClock.Today(), 1)
		expectedMax := mockClock.AddDays(mockClock.Today(), 3)
		if question.NextReview.Before(expectedMin) || question.NextReview.After(expectedMax) {
			t.Errorf("Expected NextReview to be between %v and %v for zero interval, got %v", expectedMin, expectedMax, question.NextReview)
		}
	})

	t.Run("Question with negative interval", func(t *testing.T) {
		question := &Question{
			Importance:  MediumImportance,
			Familiarity: Medium,
		}
		scheduler.setNextReview(question, mockClock.Today(), -5)

		// Should default to 1 day minimum
		expectedNextReview := mockClock.AddDays(mockClock.Today(), 1)
		if !question.NextReview.Equal(expectedNextReview) {
			t.Errorf("Expected NextReview to be %v for negative interval, got %v", expectedNextReview, question.NextReview)
		}
	})

	t.Run("Question with very large interval", func(t *testing.T) {
		question := &Question{
			Importance:  MediumImportance,
			Familiarity: Medium,
		}
		scheduler.setNextReview(question, mockClock.Today(), 1000)

		// Should be capped at maxInterval (90 days)
		expectedNextReview := mockClock.AddDays(mockClock.Today(), 90)
		if !question.NextReview.Equal(expectedNextReview) {
			t.Errorf("Expected NextReview to be %v for large interval, got %v", expectedNextReview, question.NextReview)
		}
	})
}

func TestMemoryMultipliers(t *testing.T) {
	mockClock := NewMockClock(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
	_, cfg := config.MockEnv(t)
	scheduler := NewSM2Scheduler(cfg, mockClock)

	question := &Question{
		ID:          1,
		Importance:  MediumImportance,
		Familiarity: Medium,
	}

	// Test that different memory types result in different intervals
	// Base interval for MediumImportance is 6 days, Medium familiarity adds 2 days = 8 days
	// MemoryReasoned: 8 * 1.00 = 8 days
	// MemoryPartial: 8 * 1.10 = 8.8 ≈ 9 days
	// MemoryFull: 8 * 1.25 = 10 days

	// With randomization (-1 to +2 days), the ranges are:
	// MemoryReasoned: 7-10 days
	// MemoryPartial: 8-11 days
	// MemoryFull: 9-12 days

	// Test that the expected ranges are reasonable
	expectedRanges := map[MemoryUse]struct{ min, max int }{
		MemoryReasoned: {7, 10},
		MemoryPartial:  {8, 11},
		MemoryFull:     {9, 12},
	}

	for memory, expected := range expectedRanges {
		question.NextReview = time.Time{} // Reset
		scheduler.ScheduleNewQuestion(question, memory)
		interval := int(question.NextReview.Sub(mockClock.Today()).Hours() / 24)

		if interval < expected.min || interval > expected.max {
			t.Errorf("Expected %v interval to be between %d and %d days, got %d", memory, expected.min, expected.max, interval)
		}
	}
}
