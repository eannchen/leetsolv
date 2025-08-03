package rank

import (
	"leetsolv/core"
)

func TopKSortedQuestions(questions []core.Question, k int, scoreFunc func(*core.Question) float64) []core.Question {
	if k <= 0 || len(questions) == 0 {
		return []core.Question{}
	}

	heap := NewTopKMinHeap(k)

	for i := range questions {
		heap.Push(HeapItem{
			Item:  &questions[i],
			Score: scoreFunc(&questions[i]),
		})
	}

	// Pop items in reverse order to get the highest scores first
	res := make([]core.Question, heap.Len())
	for i := len(res) - 1; i > -1; i-- {
		item, _ := heap.Pop()
		res[i] = *(item.Item.(*core.Question))
	}

	return res
}

type HeapItem struct {
	Item  any
	Score float64
}

func NewTopKMinHeap(k int) *TopKMinHeap {
	return &TopKMinHeap{
		Items: make([]HeapItem, 0, k),
		K:     k,
	}
}

type TopKMinHeap struct {
	Items []HeapItem
	K     int
}

func (h *TopKMinHeap) Len() int {
	return len(h.Items)
}

func (h *TopKMinHeap) Push(item HeapItem) {
	if len(h.Items) < h.K {
		h.Items = append(h.Items, item)
		h.percolateUp(len(h.Items) - 1)
		return
	}

	// If the item is smaller than the smallest item in the heap, it cannot be in the top-K
	if item.Score < h.Items[0].Score {
		return
	}

	h.Items[0] = item
	h.percolateDown(0)
}

func (h *TopKMinHeap) Pop() (HeapItem, bool) {
	if len(h.Items) == 0 {
		return HeapItem{}, false
	}

	item := h.Items[0]
	h.Items[0] = h.Items[len(h.Items)-1]
	h.Items = h.Items[:len(h.Items)-1]
	h.percolateDown(0)
	return item, true
}

func (h *TopKMinHeap) percolateUp(i int) {
	if i >= len(h.Items) {
		return
	}

	item := h.Items[i]

	parentI := (i - 1) / 2
	for i > 0 && item.Score < h.Items[parentI].Score {
		h.Items[i] = h.Items[parentI]
		i = parentI
		parentI = (i - 1) / 2
	}
	h.Items[i] = item
}

func (h *TopKMinHeap) percolateDown(i int) {
	if i >= len(h.Items) {
		return
	}

	for {
		smallest := i
		l := 2*i + 1
		r := 2*i + 2

		if l < len(h.Items) && h.Items[l].Score < h.Items[smallest].Score {
			smallest = l
		}
		if r < len(h.Items) && h.Items[r].Score < h.Items[smallest].Score {
			smallest = r
		}
		if smallest == i {
			break
		}
		h.Items[i], h.Items[smallest] = h.Items[smallest], h.Items[i]
		i = smallest
	}
}
