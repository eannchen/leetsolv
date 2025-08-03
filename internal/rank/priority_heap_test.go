package rank

import (
	"testing"
	"time"

	"leetsolv/core"
)

func TestTopKSortedQuestions(t *testing.T) {
	tests := []struct {
		name       string
		questions  []core.Question
		k          int
		scoreFunc  func(*core.Question) float64
		wantLen    int
		wantSorted bool
	}{
		{
			name:      "empty questions",
			questions: []core.Question{},
			k:         5,
			scoreFunc: func(q *core.Question) float64 {
				return float64(q.ID)
			},
			wantLen: 0,
		},
		{
			name: "k <= 0",
			questions: []core.Question{
				{ID: 1}, {ID: 2}, {ID: 3},
			},
			k: 0,
			scoreFunc: func(q *core.Question) float64 {
				return float64(q.ID)
			},
			wantLen: 0,
		},
		{
			name: "k > len(questions)",
			questions: []core.Question{
				{ID: 1}, {ID: 2}, {ID: 3},
			},
			k: 10,
			scoreFunc: func(q *core.Question) float64 {
				return float64(q.ID)
			},
			wantLen:    3,
			wantSorted: true,
		},
		{
			name: "normal case - top 3",
			questions: []core.Question{
				{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5},
			},
			k: 3,
			scoreFunc: func(q *core.Question) float64 {
				return float64(q.ID)
			},
			wantLen:    3,
			wantSorted: true,
		},
		{
			name: "reverse order scores",
			questions: []core.Question{
				{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5},
			},
			k: 3,
			scoreFunc: func(q *core.Question) float64 {
				return float64(10 - q.ID) // Higher scores for lower IDs
			},
			wantLen:    3,
			wantSorted: true,
		},
		{
			name: "complex scoring function",
			questions: []core.Question{
				{ID: 1, Familiarity: core.Easy, Importance: core.HighImportance},
				{ID: 2, Familiarity: core.Hard, Importance: core.MediumImportance},
				{ID: 3, Familiarity: core.Medium, Importance: core.CriticalImportance},
				{ID: 4, Familiarity: core.VeryEasy, Importance: core.LowImportance},
				{ID: 5, Familiarity: core.VeryHard, Importance: core.HighImportance},
			},
			k: 3,
			scoreFunc: func(q *core.Question) float64 {
				// Complex scoring: importance * (6 - familiarity) / 5
				return float64(q.Importance) * float64(6-int(q.Familiarity)) / 5.0
			},
			wantLen:    3,
			wantSorted: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TopKSortedQuestions(tt.questions, tt.k, tt.scoreFunc)

			if len(got) != tt.wantLen {
				t.Errorf("TopKSortedQuestions() length = %v, want %v", len(got), tt.wantLen)
			}

			if tt.wantSorted && len(got) > 1 {
				// Check that results are sorted in descending order (highest scores first)
				for i := 1; i < len(got); i++ {
					prevScore := tt.scoreFunc(&got[i-1])
					currScore := tt.scoreFunc(&got[i])
					if prevScore < currScore {
						t.Errorf("Results not sorted correctly: score[%d] = %f < score[%d] = %f",
							i-1, prevScore, i, currScore)
					}
				}
			}
		})
	}
}

func TestTopKMinHeap_BasicOperations(t *testing.T) {
	heap := NewTopKMinHeap(3)

	// Test empty heap
	if heap.Len() != 0 {
		t.Errorf("Expected empty heap length 0, got %d", heap.Len())
	}

	// Test push operations
	items := []HeapItem{
		{Item: "A", Score: 5.0},
		{Item: "B", Score: 3.0},
		{Item: "C", Score: 7.0},
		{Item: "D", Score: 1.0},
		{Item: "E", Score: 9.0},
	}

	for _, item := range items {
		heap.Push(item)
	}

	// Should only keep top 3 (highest scores)
	if heap.Len() != 3 {
		t.Errorf("Expected heap length 3, got %d", heap.Len())
	}

	// Test pop operations - should return items in ascending order (min heap)
	// The heap keeps the K largest scores, so when we pop, we get smallest of the largest first
	expectedScores := []float64{5.0, 7.0, 9.0} // Smallest to largest of the top 3 scores
	for i := 0; i < 3; i++ {
		item, ok := heap.Pop()
		if !ok {
			t.Errorf("Expected to pop item %d, but got false", i)
		}
		if item.Score != expectedScores[i] {
			t.Errorf("Expected score %f, got %f", expectedScores[i], item.Score)
		}
	}

	// Test empty pop
	_, ok := heap.Pop()
	if ok {
		t.Error("Expected pop on empty heap to return false")
	}
}

func TestTopKMinHeap_EdgeCases(t *testing.T) {
	// Test with k = 1
	heap := NewTopKMinHeap(1)
	heap.Push(HeapItem{Item: "A", Score: 5.0})
	heap.Push(HeapItem{Item: "B", Score: 3.0})
	heap.Push(HeapItem{Item: "C", Score: 7.0})

	if heap.Len() != 1 {
		t.Errorf("Expected heap length 1, got %d", heap.Len())
	}

	item, ok := heap.Pop()
	if !ok {
		t.Error("Expected to pop item, but got false")
	}
	if item.Score != 7.0 {
		t.Errorf("Expected score 7.0, got %f", item.Score)
	}

	// Test with k = 0
	heap = NewTopKMinHeap(0)
	heap.Push(HeapItem{Item: "A", Score: 5.0})
	if heap.Len() != 0 {
		t.Errorf("Expected heap length 0, got %d", heap.Len())
	}

	// Test pop on empty heap
	_, ok = heap.Pop()
	if ok {
		t.Error("Expected pop on empty heap to return false")
	}
}

func TestTopKMinHeap_HeapProperty(t *testing.T) {
	heap := NewTopKMinHeap(10)

	// Add items in random order
	items := []HeapItem{
		{Item: "A", Score: 10.0},
		{Item: "B", Score: 5.0},
		{Item: "C", Score: 15.0},
		{Item: "D", Score: 3.0},
		{Item: "E", Score: 8.0},
		{Item: "F", Score: 12.0},
		{Item: "G", Score: 1.0},
		{Item: "H", Score: 20.0},
	}

	for _, item := range items {
		heap.Push(item)
	}

	// Verify heap property: parent should be smaller than children
	for i := 0; i < heap.Len(); i++ {
		left := 2*i + 1
		right := 2*i + 2

		if left < heap.Len() && heap.Items[i].Score > heap.Items[left].Score {
			t.Errorf("Heap property violated: parent[%d] = %f > left[%d] = %f",
				i, heap.Items[i].Score, left, heap.Items[left].Score)
		}

		if right < heap.Len() && heap.Items[i].Score > heap.Items[right].Score {
			t.Errorf("Heap property violated: parent[%d] = %f > right[%d] = %f",
				i, heap.Items[i].Score, right, heap.Items[right].Score)
		}
	}
}

func TestTopKMinHeap_PercolateUp(t *testing.T) {
	heap := NewTopKMinHeap(5)

	// Manually construct a heap that violates the heap property
	heap.Items = []HeapItem{
		{Item: "A", Score: 10.0}, // Root
		{Item: "B", Score: 5.0},  // Left child
		{Item: "C", Score: 15.0}, // Right child
		{Item: "D", Score: 3.0},  // This should percolate up
	}

	// Test percolateUp
	heap.percolateUp(3) // Percolate the last item (score 3.0)

	// Verify the smallest item is at the root
	if heap.Items[0].Score != 3.0 {
		t.Errorf("Expected root score 3.0, got %f", heap.Items[0].Score)
	}
}

func TestTopKMinHeap_PercolateDown(t *testing.T) {
	heap := NewTopKMinHeap(5)

	// Manually construct a heap that violates the heap property
	heap.Items = []HeapItem{
		{Item: "A", Score: 15.0}, // Root (should percolate down)
		{Item: "B", Score: 5.0},  // Left child
		{Item: "C", Score: 10.0}, // Right child
		{Item: "D", Score: 8.0},  // Left-left child
		{Item: "E", Score: 12.0}, // Left-right child
	}

	// Test percolateDown
	heap.percolateDown(0) // Percolate the root

	// Verify the smallest item is at the root
	if heap.Items[0].Score != 5.0 {
		t.Errorf("Expected root score 5.0, got %f", heap.Items[0].Score)
	}
}

func TestTopKSortedQuestions_WithRealQuestions(t *testing.T) {
	now := time.Now()
	questions := []core.Question{
		{
			ID: 1, Familiarity: core.Easy, Importance: core.HighImportance,
			LastReviewed: now.Add(-24 * time.Hour), ReviewCount: 5, EaseFactor: 2.5,
		},
		{
			ID: 2, Familiarity: core.Hard, Importance: core.CriticalImportance,
			LastReviewed: now.Add(-48 * time.Hour), ReviewCount: 2, EaseFactor: 1.8,
		},
		{
			ID: 3, Familiarity: core.Medium, Importance: core.MediumImportance,
			LastReviewed: now.Add(-12 * time.Hour), ReviewCount: 3, EaseFactor: 2.0,
		},
		{
			ID: 4, Familiarity: core.VeryHard, Importance: core.HighImportance,
			LastReviewed: now.Add(-72 * time.Hour), ReviewCount: 1, EaseFactor: 1.5,
		},
		{
			ID: 5, Familiarity: core.VeryEasy, Importance: core.LowImportance,
			LastReviewed: now.Add(-6 * time.Hour), ReviewCount: 8, EaseFactor: 3.0,
		},
	}

	// Test with different scoring functions
	scoreFuncs := []struct {
		name string
		fn   func(*core.Question) float64
	}{
		{
			name: "importance only",
			fn: func(q *core.Question) float64 {
				return float64(q.Importance)
			},
		},
		{
			name: "familiarity only",
			fn: func(q *core.Question) float64 {
				return float64(q.Familiarity)
			},
		},
		{
			name: "combined score",
			fn: func(q *core.Question) float64 {
				return float64(q.Importance)*2 + float64(q.Familiarity) + q.EaseFactor
			},
		},
	}

	for _, sf := range scoreFuncs {
		t.Run(sf.name, func(t *testing.T) {
			result := TopKSortedQuestions(questions, 3, sf.fn)

			if len(result) != 3 {
				t.Errorf("Expected 3 results, got %d", len(result))
			}

			// Verify sorting
			for i := 1; i < len(result); i++ {
				prevScore := sf.fn(&result[i-1])
				currScore := sf.fn(&result[i])
				if prevScore < currScore {
					t.Errorf("Results not sorted correctly: score[%d] = %f < score[%d] = %f",
						i-1, prevScore, i, currScore)
				}
			}
		})
	}
}

func BenchmarkTopKSortedQuestions(b *testing.B) {
	questions := make([]core.Question, 1000)
	for i := range questions {
		questions[i] = core.Question{
			ID:          i + 1,
			Familiarity: core.Familiarity(i % 5),
			Importance:  core.Importance(i % 4),
			EaseFactor:  1.5 + float64(i%10)/10.0,
		}
	}

	scoreFunc := func(q *core.Question) float64 {
		return float64(q.Importance)*2 + float64(q.Familiarity) + q.EaseFactor
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TopKSortedQuestions(questions, 100, scoreFunc)
	}
}

func BenchmarkTopKMinHeap_Push(b *testing.B) {
	heap := NewTopKMinHeap(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Push(HeapItem{
			Item:  i,
			Score: float64(i),
		})
	}
}
