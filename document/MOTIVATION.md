# Motivation & Design Notes

*(Updated by: Ian Chen, Date: 2025-08-20)*

This document explains why I built LeetSolv, how I adapted the SM-2 algorithm, and what design and efficiency choices went into the project. In the AI era, I want to show that this project is intentional and serious, not just AI-generated code.

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
  - [Data Safety \& Caching](#data-safety--caching)
  - [Design Patterns](#design-patterns)
  - [Closing Note](#closing-note)


*<span style="color: #888">[🏗️ AN EXPLANATION VIDEO FOR THIS DOC WILL BE EMBEDED AT HERE LATER]</span>*

## Motivation

After solving 190+ LeetCode problems, I noticed something missing: my understanding didn’t always stick. I was moving forward but not deepening.

My old method was starring ⭐️ hard problems, but it wasn’t reliable: some stars became trivial with progression, while other tough problems slipped through.

I thought back to my English-learning experience: flashcards and spaced repetition worked well for vocabulary. But DSA isn’t like vocabulary. Memorization isn't the correct way to learn DSA, it requires reasoning, practice, and reviewing concepts in different contexts. I can't just use a software like Anki to review DSA.

So I made LeetSolv to solve my own learning problem: a revision tool that schedules problem reviews like flashcards, but adapts the method for algorithm practice.


### Zero Dependencies Philosophy

On LeetCode, it encourages developers to adapt the fundamentals to solve problems, cultivating developers to be tool makers instead of tool users.

Since I am practicing problem-solving, why don't I challenge myself to apply what I’ve learned here? So I gave myself the challenge to build LeetSolv without any dependencies. The tool is not only a revision tool, but also a chance to apply what I’ve learned.

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
- **Memory usage factor:** Some problems inevitably feel familiar once we’ve seen their patterns. These should be spaced out more aggressively to prevent rote memorization, by lowering their retention score.
- **Randomization:** Many SM-2 apps bunch reviews into heavy days — a common pitfall I wanted to avoid.

Together, these changes shift SM-2 from “memorize facts” toward “practice reasoning.”

See my implementation in [scheduler.go](../core/scheduler.go).


## Data Structures

LeetSolv is small. Users might only add hundreds of problems, far below where efficiency actually matters. But I treated this project as a chance to apply what I’ve learned — to write efficient structures and think about time complexity in a real project, not just in interview problems.

Interestingly, 2025’s AI often suggested workable but inefficient solutions (e.g. linear scans instead of binary search or hash maps). By pushing back and rethinking, I practiced writing efficiency consciously.

### Indexing for Lookup

Users should be able to search problems by both ID and URL. Originally, AI suggested arrays with linear scans, but I replaced them with a hash-based approach, it provides `amortized O(1)` lookup in both cases. Even binary search wasn’t needed.

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

In linear scans, it takes `O(n * m)` time to search, where `n` is the length of searching string and `m` is the average length of the strings in the dataset.


```go
func main() {
	text := "This is a sample text for a search."
	pattern := "sample"

	index := strings.Index(text, pattern)

	if index != -1 {
		fmt.Printf("Pattern found at index: %d\n", index)
	} else {
		fmt.Println("Pattern not found.")
	}
}
```

By implementing a trie, I can mitigate the time complexity to `O(n)`, with only `O(k)` extra space, where k is the number of unique characters across all entries. See my implementation in [trie.go](../internal/search/trie.go).

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

By using a heap, we can improve the efficiency to `O(log k)` time, where `k` is the number of problems to return, which is drastically better than `O(n log n)`.

But when we use the built-in Heap library, we need to spend `O(log k)` time to remove the smallest element and spend another `O(log k)` time to insert the new element, which takes `O(2 log k)` time in total.

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

With a custom heap implementation, I gained more control, avoided unnecessary percolate-up steps, and kept it purely `O(log k)`.

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

See my implementation in [priority_heap.go](../internal/heap/priority_heap.go).


### Stack for Undo



## Data Safety & Caching

In interactive mode, the app caches data in memory instead of reloading the storage file on every operation.

TODO

See my implementation in [fileutil.go](../internal/fileutil/fileutil.go).

I applied **write-through** and **cache-aside** strategies. This way, performance improves while data safety is still guaranteed.

TODO

See my implementation in [storage.go](../storage/file.go).



## Design Patterns

TODO


## Closing Note

LeetSolv is both a tool I use and a learning project I grow with.
It’s built not just to work, but to show my care for efficiency, design, and fundamentals — even when that depth isn’t strictly necessary for a small tool.
