package rank

import (
	"testing"
)

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
