package search

import (
	"testing"
)

func TestNewTrieNode(t *testing.T) {
	node := NewTrieNode()

	if node == nil {
		t.Fatal("NewTrieNode() returned nil")
	}

	if node.Children == nil {
		t.Error("Children map should not be nil")
	}

	if node.IDs == nil {
		t.Error("IDs map should not be nil")
	}

	if node.IsWord != false {
		t.Error("IsWord should be false for new node")
	}
}

func TestNewTrie(t *testing.T) {
	minPrefixLength := 3
	trie := NewTrie(minPrefixLength)

	if trie == nil {
		t.Fatal("NewTrie() returned nil")
	}

	if trie.Root == nil {
		t.Error("Root should not be nil")
	}

	if trie.MinPrefixLength != minPrefixLength {
		t.Errorf("MinPrefixLength is %d, expected %d", trie.MinPrefixLength, minPrefixLength)
	}
}

func TestTrie_Insert(t *testing.T) {
	trie := NewTrie(1)

	// Test inserting a single word
	trie.Insert("hello", 1)

	// Check that the word is marked as complete
	node := trie.Root
	for _, ch := range "hello" {
		if node.Children[ch] == nil {
			t.Errorf("Expected child node for character %c", ch)
		}
		node = node.Children[ch]
		if _, exists := node.IDs[1]; !exists {
			t.Errorf("Expected ID 1 to be present at character %c", ch)
		}
	}

	if !node.IsWord {
		t.Error("Expected final node to be marked as word")
	}
}

func TestTrie_InsertMultipleWords(t *testing.T) {
	trie := NewTrie(1)

	// Insert multiple words
	trie.Insert("hello", 1)
	trie.Insert("world", 2)
	trie.Insert("help", 3)

	// Check that all words are present
	helloIDs := trie.SearchPrefix("hello")
	if _, exists := helloIDs[1]; !exists {
		t.Error("Expected ID 1 for 'hello'")
	}

	worldIDs := trie.SearchPrefix("world")
	if _, exists := worldIDs[2]; !exists {
		t.Error("Expected ID 2 for 'world'")
	}

	helpIDs := trie.SearchPrefix("help")
	if _, exists := helpIDs[3]; !exists {
		t.Error("Expected ID 3 for 'help'")
	}
}

func TestTrie_InsertSameWordMultipleTimes(t *testing.T) {
	trie := NewTrie(1)

	// Insert the same word with different IDs
	trie.Insert("hello", 1)
	trie.Insert("hello", 2)

	// Check that both IDs are present
	ids := trie.SearchPrefix("hello")
	if _, exists := ids[1]; !exists {
		t.Error("Expected ID 1 for 'hello'")
	}
	if _, exists := ids[2]; !exists {
		t.Error("Expected ID 2 for 'hello'")
	}
}

func TestTrie_SearchPrefix(t *testing.T) {
	trie := NewTrie(2)

	// Insert test data
	trie.Insert("hello", 1)
	trie.Insert("world", 2)
	trie.Insert("help", 3)

	// Test exact prefix match
	ids := trie.SearchPrefix("hello")
	if len(ids) != 1 {
		t.Errorf("Expected 1 ID for 'hello', got %d", len(ids))
	}
	if _, exists := ids[1]; !exists {
		t.Error("Expected ID 1 for 'hello'")
	}

	// Test partial prefix match
	ids = trie.SearchPrefix("hel")
	if len(ids) != 2 {
		t.Errorf("Expected 2 IDs for 'hel', got %d", len(ids))
	}
	if _, exists := ids[1]; !exists {
		t.Error("Expected ID 1 for 'hel'")
	}
	if _, exists := ids[3]; !exists {
		t.Error("Expected ID 3 for 'hel'")
	}

	// Test prefix shorter than MinPrefixLength
	ids = trie.SearchPrefix("h")
	if ids != nil {
		t.Errorf("Expected nil for prefix shorter than MinPrefixLength, got %v", ids)
	}

	// Test non-existent prefix
	ids = trie.SearchPrefix("xyz")
	if ids != nil {
		t.Errorf("Expected nil for non-existent prefix, got %v", ids)
	}
}

func TestTrie_SearchPrefix_EdgeCases(t *testing.T) {
	trie := NewTrie(1)

	// Test empty prefix
	ids := trie.SearchPrefix("")
	if ids == nil {
		t.Error("Empty prefix should return all IDs")
	}

	// Test with no words inserted
	ids = trie.SearchPrefix("hello")
	if ids != nil {
		t.Errorf("Expected nil for non-existent prefix, got %v", ids)
	}
}

func TestTrie_Delete(t *testing.T) {
	trie := NewTrie(1)

	// Insert test data
	trie.Insert("hello", 1)
	trie.Insert("hello", 2)
	trie.Insert("world", 3)

	// Delete one occurrence of "hello"
	trie.Delete("hello", 1)

	// Check that ID 1 is removed but ID 2 remains
	ids := trie.SearchPrefix("hello")
	if _, exists := ids[1]; exists {
		t.Error("ID 1 should be removed")
	}
	if _, exists := ids[2]; !exists {
		t.Error("ID 2 should remain")
	}

	// Check that "world" is unaffected
	worldIDs := trie.SearchPrefix("world")
	if _, exists := worldIDs[3]; !exists {
		t.Error("ID 3 for 'world' should remain")
	}
}

func TestTrie_Delete_NonExistentWord(t *testing.T) {
	trie := NewTrie(1)

	// Try to delete a word that doesn't exist
	trie.Delete("nonexistent", 1)

	// Should not cause any errors or panics
	// Just verify the trie is still functional
	_ = trie.SearchPrefix("hello")
}

func TestTrie_Delete_NonExistentID(t *testing.T) {
	trie := NewTrie(1)

	// Insert a word
	trie.Insert("hello", 1)

	// Try to delete a non-existent ID
	trie.Delete("hello", 999)

	// Check that the original ID is still present
	ids := trie.SearchPrefix("hello")
	if _, exists := ids[1]; !exists {
		t.Error("ID 1 should remain after deleting non-existent ID")
	}
}

func TestTrie_Delete_AllOccurrences(t *testing.T) {
	trie := NewTrie(1)

	// Insert the same word multiple times
	trie.Insert("hello", 1)
	trie.Insert("hello", 2)
	trie.Insert("hello", 3)

	// Delete all occurrences
	trie.Delete("hello", 1)
	trie.Delete("hello", 2)
	trie.Delete("hello", 3)

	// Check that the word is completely removed
	ids := trie.SearchPrefix("hello")
	if len(ids) != 0 {
		t.Errorf("Expected 0 IDs after deleting all occurrences, got %d", len(ids))
	}
}

func TestTrie_ComplexOperations(t *testing.T) {
	trie := NewTrie(1)

	// Insert various words
	words := []struct {
		word string
		id   int
	}{
		{"algorithm", 1},
		{"algorithms", 2},
		{"algo", 3},
		{"binary", 4},
		{"binarysearch", 5},
		{"tree", 6},
		{"trie", 7},
	}

	for _, w := range words {
		trie.Insert(w.word, w.id)
	}

	// Test prefix search for "algo"
	ids := trie.SearchPrefix("algo")
	expectedIDs := map[int]struct{}{1: {}, 2: {}, 3: {}}

	if len(ids) != len(expectedIDs) {
		t.Errorf("Expected %d IDs for 'algo', got %d", len(expectedIDs), len(ids))
	}

	for id := range expectedIDs {
		if _, exists := ids[id]; !exists {
			t.Errorf("Expected ID %d for 'algo'", id)
		}
	}

	// Test prefix search for "bin"
	ids = trie.SearchPrefix("bin")
	expectedIDs = map[int]struct{}{4: {}, 5: {}}

	if len(ids) != len(expectedIDs) {
		t.Errorf("Expected %d IDs for 'bin', got %d", len(expectedIDs), len(ids))
	}

	for id := range expectedIDs {
		if _, exists := ids[id]; !exists {
			t.Errorf("Expected ID %d for 'bin'", id)
		}
	}

	// Delete some words and verify
	trie.Delete("algorithms", 2)
	trie.Delete("binarysearch", 5)

	// Check that "algo" now only has 2 IDs
	ids = trie.SearchPrefix("algo")
	if len(ids) != 2 {
		t.Errorf("Expected 2 IDs for 'algo' after deletion, got %d", len(ids))
	}

	// Check that "bin" now only has 1 ID
	ids = trie.SearchPrefix("bin")
	if len(ids) != 1 {
		t.Errorf("Expected 1 ID for 'bin' after deletion, got %d", len(ids))
	}
}

func TestTrie_UnicodeSupport(t *testing.T) {
	trie := NewTrie(1)

	// Test with Unicode characters
	trie.Insert("café", 1)
	trie.Insert("naïve", 2)
	trie.Insert("résumé", 3)

	// Test prefix search
	ids := trie.SearchPrefix("caf")
	if _, exists := ids[1]; !exists {
		t.Error("Expected ID 1 for 'caf'")
	}

	ids = trie.SearchPrefix("naï")
	if _, exists := ids[2]; !exists {
		t.Error("Expected ID 2 for 'naï'")
	}

	ids = trie.SearchPrefix("rés")
	if _, exists := ids[3]; !exists {
		t.Error("Expected ID 3 for 'rés'")
	}
}

func TestTrie_EmptyString(t *testing.T) {
	trie := NewTrie(1)

	// Test inserting empty string
	trie.Insert("", 1)

	// Test searching empty string
	ids := trie.SearchPrefix("")
	if _, exists := ids[1]; !exists {
		t.Error("Expected ID 1 for empty string")
	}

	// Test deleting empty string
	trie.Delete("", 1)
	ids = trie.SearchPrefix("")
	if len(ids) != 0 {
		t.Errorf("Expected 0 IDs after deleting empty string, got %d", len(ids))
	}
}
