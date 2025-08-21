# Motivation & Design Notes

*Author: Ian Chen, Last Update: 2025-08-21*

This document explains why I built LeetSolv, how I adapted the SM-2 algorithm, and what design and efficiency choices went into the project. In the AI "vibe coding" era, I want to show that this project is intentional and serious, not just toy code.

- [Motivation \& Design Notes](#motivation--design-notes)
  - [Motivation](#motivation)
    - [Zero Dependencies Philosophy](#zero-dependencies-philosophy)
  - [Why SM-2 Algorithm](#why-sm-2-algorithm)
    - [Custom Adaptations to SM-2](#custom-adaptations-to-sm-2)
  - [Data Structures](#data-structures)
    - [Indexing for Lookup](#indexing-for-lookup)
    - [Trie for Text Search](#trie-for-text-search)
    - [Heap for Top-K Problems](#heap-for-top-k-problems)
    - [Stack for Undo](#stack-for-undo)
  - [Storage](#storage)
    - [Atomic File Write](#atomic-file-write)
    - [Caching](#caching)
  - [Design Patterns](#design-patterns)
  - [Closing Note](#closing-note)


*<span style="color: #888">[🏗️ AN EXPLANATION VIDEO FOR THIS DOC WILL BE EMBEDED AT HERE LATER]</span>*

## Motivation

After solving 190+ LeetCode problems, I noticed something missing: my understanding didn’t always stick. I was moving forward but not deepening.

My old method was starring ⭐️ hard problems, but it wasn’t reliable: some stars became trivial with progression, while other tough problems slipped through. Moreover, "when" to review those problems is also a challenge.

I thought back to my English-learning experience: flashcards and spaced repetition worked well for vocabulary. But DSA isn’t like vocabulary. Memorization isn't the correct way to learn DSA, it requires reasoning, practice, and reviewing concepts in different contexts. I can't just use a software like Anki to review DSA.

So I made LeetSolv to solve my own learning problem: a revision tool that schedules problem reviews like flashcards, but adapts the method for algorithm practice.


### Zero Dependencies Philosophy

On LeetCode, it encourages developers to adapt the fundamentals to solve problems, cultivating developers to be tool makers instead of tool users.

Since I am practicing problem-solving, why don't I challenge myself to apply what I’ve learned here? So I gave myself the challenge to build LeetSolv without any dependencies. The tool is not only a revision tool, but also a chance to apply what I’ve learned. To achieve this, even some Go built-in libraries for data structures are not used. Given that, I can even control the subtle details for better time efficiency, e.g., [Heap for Top-K Problems](#heap-for-top-k-problems).

It’s a rare chance, because in real jobs developers usually rely on libraries for speed and reliability.


## Why SM-2 Algorithm

Newer versions of the SuperMemo algorithm exist and are backed by science, but I chose SM-2 for two reasons:

- It’s simple and easy to understand.
- It’s proven in practice (used by Anki, which works for millions of learners).

That balance made it a good foundation for customization. For my use case, adaptability matters more than chasing algorithmic “perfection.”

### Custom Adaptations to SM-2

I adjusted SM-2 so it works better for DSA instead of vocabulary.
The key customizations are:

- **Importance factor:** NeetCode 150 problems are essential building blocks, while NeetCode 250 includes duplicates for extra practice. They shouldn’t be scheduled equally.
- **Reasoning factor:** If a problem is solved from memory instead of reasoning, the algorithm treats it as a weaker recall signal. These should be spaced out more aggressively to prevent rote memorization and lower their retention score.
- **Randomization:** Many SM-2 apps bunch reviews into heavy days — a common pitfall I wanted to avoid.

Together, these changes shift SM-2 from “memorize facts” toward “practice reasoning.”

See my implementation in [scheduler.go](../core/scheduler.go).


## Data Structures

LeetSolv is small. Users might only add hundreds of problems, far below where efficiency actually matters. But I treated this project as a chance to apply what I’ve learned — to write efficient structures and think about time complexity in a real project, not just in interview problems.

Interestingly, AI models (in 2025) often suggested workable but inefficient solutions unless I explicitly point to improve efficiency. (e.g. linear scans with `O(n)` or `O(n * m)` time complexity instead of indexing or using more complex data structures). **I confirmed that learning data structures and algorithms is still essential.**

### Indexing for Lookup

Users should be able to search problems by both ID and URL. It provides `average O(1)` time operations for both lookups instead of linear scans.

```json
{
    "max_id": 1,
    "questions": {
        "1": {
            "id": "1",
            "url": "https://leetcode.com/problems/two-sum/"
        }
    },
    "url_index": {
        "https://leetcode.com/problems/two-sum/": "1"
    }
}
```

### Trie for Text Search

In linear scans, it takes `O(n * m)` time to search, where `n` is the number of strings in the dataset and `m` is the average length of a string.

```go
// search linearly through a dataset.
// Complexity: O(n * m)
// - n: number of strings in the dataset
// - m: average length of a string
func linearSearch(dataset []string, pattern string) []string {
	var results []string

	for _, item := range dataset {
		if strings.Contains(item, pattern) {
			results = append(results, item)
		}
	}

	return results
}
```

By implementing a trie, search time drops to `O(L)`, where `L` is just the length of the search query itself. The key win here is that search performance is now **completely independent of the dataset size**. Whether I have 100 problems or 10,000, the search takes the same amount of time.

```json
"url_trie": {
    "Root": {
      "Children": {
        "115": {
          "Children": {
            "117": {
              "Children": {
                "109": {
                  "Children": {},
                  "IDs": {
                    "1": {}
                  },
                  "IsWord": true
                }
              },
              // ...
              // ...
    "MinPrefixLength": 3
  }
```

```go
func (t *Trie) SearchPrefix(prefix string) map[int]struct{} {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if prefix == "" {
		idsCopy := make(map[int]struct{}, len(t.Root.IDs))
		for id := range t.Root.IDs {
			idsCopy[id] = struct{}{}
		}
		return idsCopy
	}

	if len([]rune(prefix)) < t.MinPrefixLength {
		return make(map[int]struct{})
	}

	node := t.Root
	for _, ch := range prefix {
		if _, ok := node.Children[ch]; !ok {
			return make(map[int]struct{})
		}
		node = node.Children[ch]
	}

	idsCopy := make(map[int]struct{}, len(node.IDs))
	for id := range node.IDs {
		idsCopy[id] = struct{}{}
	}
	return idsCopy
}
```
> *The app is single-threaded, so a race condition is impossible. But I kept the mutex around for ever features like a background process are added.*

See my implementation in [trie.go](../internal/search/trie.go).

There’s a trade-off: a trie only supports prefix search, not fuzzy search. But I chose it because I fully understand it and wanted to avoid external libraries.

### Heap for Top-K Problems

SM-2 alone can create a backlog problem: too many reviews pile up, and it’s hard to decide what to tackle first.

To address this, I added a Due Priority Scoring system. It surfaces the most urgent and important problems first, making the review queue manageable and meaningful.

In brute force, we can sort the problems by due priority and return the top-k problems. But this takes `O(n log n)` time, where `n` is the number of problems.

```go
if k > len(questions) {
    k = len(questions)
}
sort.Slice(questions, func(i, j int) bool {
	return questions[i].DuePriority < questions[j].DuePriority
})
return questions[:k]
```

By using a min-heap of size `k`, I can find the top k problems in `O(n log k)` time, which is a major improvement. We iterate through all n problems, maintaining a heap of the k problems with the highest priority scores seen so far.

```go
func (h *TopKMinHeap) Push(item HeapItem) {
	if h.K <= 0 {
		return
	}

	if h.h.Len() < h.K {
		heap.Push(h.h, item)
	} else {
		if item.Score > (*h.h)[0].Score {
			heap.Pop(h.h)         // O(log k)
			heap.Push(h.h, item)  // O(log k)
		}
	}
}
```

I also implemented a custom heap to optimize the process further. Instead of using the standard library's `heap.Pop` and `heap.Push` (**two** O`(log k)` operations), my implementation directly replaces the heap's root if a new item has a higher score and performs a **single** `percolateDown` operation.

```go
func (h *TopKMinHeap) Push(item HeapItem) {
	if h.K <= 0 {
		return
	}

	if len(h.Items) < h.K {
		h.Items = append(h.Items, item)
		h.percolateUp(len(h.Items) - 1)
		return
	}

	if item.Score <= h.Items[0].Score {
		return
	}

	h.Items[0] = item
	h.percolateDown(0)  // O(log k)
}
```

See my implementation in [priority_heap.go](../internal/rank/priority_heap.go).


### Stack for Undo

Stack is a natural choice for tracking a history owning to its LIFO behavior. Every change whether it's an `add`, `update`, or `delete` is captured in a `Delta` object. It is clean, efficient, and requires no external libraries.

```go
type Delta struct {
    Action     ActionType `json:"action"`
    QuestionID int        `json:"question_id"`
    OldState   *Question  `json:"old_state"` // The question's state before the change
    NewState   *Question  `json:"new_state"` // The question's state after the change
}
```

## Storage

### Atomic File Write

[README](../README.md) has depicted the atomic file write process with a diagram.

See my implementation in [fileutil.go](../internal/fileutil/fileutil.go).

### Caching

In interactive mode, the app caches data in memory instead of reloading the storage file on every operation.
I applied **Cache-aside** and **Write-through** strategies. This way, performance is improved.

> *The app is single-threaded, so a race condition is impossible. But I kept the mutex around for ever features like a background process are added.*

```go
func (fs *FileStorage) LoadDeltas() ([]core.Delta, error) {
	// Try to load from cache first
	fs.deltasCacheMutex.RLock()
	if fs.deltasCache != nil {
		fs.deltasCacheMutex.RUnlock()
		return fs.deltasCache, nil
	}
	fs.deltasCacheMutex.RUnlock()

	fs.deltasCacheMutex.Lock()
	defer fs.deltasCacheMutex.Unlock()

	// Load deltas from file
	var deltas []core.Delta
	err := fs.file.Load(&deltas, fs.deltasFileName)
	if err != nil {
		return nil, err
	}

	// Update cache
	fs.deltasCache = deltas

	return deltas, nil
}

func (fs *FileStorage) SaveDeltas(deltas []core.Delta) error {
	fs.deltasCacheMutex.Lock()
	defer fs.deltasCacheMutex.Unlock()

	err := fs.file.Save(deltas, fs.deltasFileName)
	if err != nil {
		return err
	}

	// Update cache after successful save
	fs.deltasCache = deltas

	return nil
}
```

See my implementation in [storage.go](../storage/file.go).

## Design Patterns

*<span style="color: #888">[🏗️ Work in progress]</span>*


## Closing Note

LeetSolv is both a tool I use and a learning project I grow with.
It’s built not just to work, but to show my care for efficiency, design, and fundamentals — even when that depth isn’t strictly necessary for a small tool.
