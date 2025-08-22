package search

import (
	"reflect"
	"sort"
	"testing"
)

// newHelper is a helper function to compare the actual IDs with expected IDs.
// It improves readability and reduces boilerplate in tests.
func assertIDsMatch(t *testing.T, actual map[int]struct{}, expected ...int) {
	t.Helper() // Marks this function as a test helper.

	if len(actual) != len(expected) {
		t.Errorf("expected %d IDs, but got %d. Expected: %v, Got: %v", len(expected), len(actual), expected, actual)
		return
	}

	expectedMap := make(map[int]struct{})
	for _, id := range expected {
		expectedMap[id] = struct{}{}
	}

	if !reflect.DeepEqual(actual, expectedMap) {
		// Sorting slices for consistent error message output
		actualSlice := []int{}
		for id := range actual {
			actualSlice = append(actualSlice, id)
		}
		sort.Ints(actualSlice)
		sort.Ints(expected)
		t.Errorf("ID sets do not match. Expected %v, got %v", expected, actualSlice)
	}
}

func TestNewTrieNode(t *testing.T) {
	node := NewTrieNode()

	if node == nil {
		t.Fatal("NewTrieNode() should not return nil")
	}
	if node.Children == nil {
		t.Error("Children map should be initialized, not nil")
	}
	if node.IDs == nil {
		t.Error("IDs map should be initialized, not nil")
	}
	// CHANGE: The original test checked for a boolean `IsWord`.
	// The implementation uses a map `WordEndIDs`, which should be checked instead.
	if node.WordEndIDs == nil {
		t.Error("WordEndIDs map should be initialized, not nil")
	}
}

func TestTrie_Hydrate(t *testing.T) {
	// Manually build a broken trie (simulating bad deserialization)
	trie := &Trie{Root: &TrieNode{}} // Root node with nil maps
	trie.Root.Children = map[rune]*TrieNode{
		'a': {}, // Child node with nil maps
	}

	trie.Hydrate()

	// Root should have all maps initialized
	if trie.Root.Children == nil {
		t.Error("Root.Children should not be nil after Hydrate")
	}
	if trie.Root.IDs == nil {
		t.Error("Root.IDs should be initialized after Hydrate")
	}
	if trie.Root.WordEndIDs == nil {
		t.Error("Root.WordEndIDs should be initialized after Hydrate")
	}

	// Child 'a' should also have its maps initialized
	child := trie.Root.Children['a']
	if child.Children == nil {
		t.Error("Child.Children should be initialized after Hydrate")
	}
	if child.IDs == nil {
		t.Error("Child.IDs should be initialized after Hydrate")
	}
	if child.WordEndIDs == nil {
		t.Error("Child.WordEndIDs should be initialized after Hydrate")
	}
}

func TestTrie_Insert(t *testing.T) {
	trie := NewTrie(1)
	word := "hello"
	id := 101

	trie.Insert(word, id)

	// Check that the ID is present at every node along the path
	node := trie.Root
	if _, exists := node.IDs[id]; !exists {
		t.Errorf("Expected ID %d to be present at the root", id)
	}
	for _, ch := range word {
		child, ok := node.Children[ch]
		if !ok {
			t.Fatalf("Expected child node for character '%c', but it was nil", ch)
		}
		node = child
		if _, exists := node.IDs[id]; !exists {
			t.Errorf("Expected ID %d to be present at node for character '%c'", id, ch)
		}
	}

	// CHANGE: Verify the ID is in WordEndIDs at the final node.
	if _, exists := node.WordEndIDs[id]; !exists {
		t.Error("Expected final node to have the ID in its WordEndIDs map")
	}
}

func TestTrie_SearchPrefix(t *testing.T) {
	trie := NewTrie(3) // Min prefix length is 3
	trie.Insert("apple", 1)
	trie.Insert("apply", 2)
	trie.Insert("application", 3)
	trie.Insert("banana", 4)
	trie.Insert("", 5) // Empty string case

	testCases := []struct {
		name         string
		prefix       string
		expectedIDs  []int
		expectNonNil bool // To handle the new return behavior
	}{
		{
			name:        "Exact Match",
			prefix:      "apple",
			expectedIDs: []int{1},
		},
		{
			name:        "Common Prefix",
			prefix:      "app",
			expectedIDs: []int{1, 2, 3},
		},
		{
			name:   "Prefix Not Found",
			prefix: "xyz",
			// CHANGE: Expect an empty map, not nil
			expectedIDs: []int{},
		},
		{
			name:   "Prefix Too Short",
			prefix: "ap",
			// CHANGE: Expect an empty map, not nil
			expectedIDs: []int{},
		},
		{
			name:        "Empty Prefix returns all IDs",
			prefix:      "",
			expectedIDs: []int{1, 2, 3, 4, 5},
		},
		{
			name:        "Search in an empty trie",
			prefix:      "any",
			expectedIDs: []int{},
		},
	}

	// Separate test for empty trie
	emptyTrie := NewTrie(1)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var ids map[int]struct{}
			if tc.name == "Search in an empty trie" {
				ids = emptyTrie.SearchPrefix(tc.prefix)
			} else {
				ids = trie.SearchPrefix(tc.prefix)
			}

			// CHANGE: The logic is now simpler. We always expect a map.
			// We no longer need to check for `nil`.
			if ids == nil {
				t.Fatal("SearchPrefix should never return a nil map")
			}

			assertIDsMatch(t, ids, tc.expectedIDs...)
		})
	}
}

func TestTrie_Delete(t *testing.T) {
	t.Run("Delete specific ID from word with multiple IDs", func(t *testing.T) {
		trie := NewTrie(1)
		trie.Insert("flow", 10)
		trie.Insert("flow", 20)
		trie.Insert("flower", 30)

		trie.Delete("flow", 10)

		// "flow" prefix should now only contain IDs 20 and 30
		ids := trie.SearchPrefix("flow")
		assertIDsMatch(t, ids, 20, 30)

		// Check that the WordEndIDs for "flow" no longer contains 10
		node := trie.Root
		for _, r := range "flow" {
			node = node.Children[r]
		}
		if _, exists := node.WordEndIDs[10]; exists {
			t.Error("ID 10 should be removed from WordEndIDs of 'flow'")
		}
		if _, exists := node.WordEndIDs[20]; !exists {
			t.Error("ID 20 should still exist in WordEndIDs of 'flow'")
		}
	})

	t.Run("Delete completely removes node path if unused", func(t *testing.T) {
		trie := NewTrie(1)
		trie.Insert("team", 1)
		trie.Insert("tea", 2)

		trie.Delete("team", 1)

		// Node for 'm' should be pruned/deleted
		node := trie.Root
		for _, r := range "tea" {
			node = node.Children[r]
		}
		if _, exists := node.Children['m']; exists {
			t.Error("Node for 'm' should have been pruned after deletion")
		}
		assertIDsMatch(t, trie.SearchPrefix("tea"), 2)
	})

	t.Run("Delete non-existent word or ID does not panic", func(t *testing.T) {
		trie := NewTrie(1)
		trie.Insert("hello", 1)

		// Should not panic or alter the trie's state
		trie.Delete("world", 2)
		trie.Delete("hello", 99)

		assertIDsMatch(t, trie.SearchPrefix("hello"), 1)
	})

	t.Run("Delete empty string", func(t *testing.T) {
		trie := NewTrie(1)
		trie.Insert("", 1)
		trie.Insert("a", 2)

		assertIDsMatch(t, trie.SearchPrefix(""), 1, 2)

		trie.Delete("", 1)

		if _, exists := trie.Root.WordEndIDs[1]; exists {
			t.Error("ID for empty string should be removed from root's WordEndIDs")
		}
		assertIDsMatch(t, trie.SearchPrefix(""), 2)
	})
}
