package search

import "sync"

type TrieNode struct {
	Children   map[rune]*TrieNode
	IDs        map[int]struct{} // Questions having this prefix
	WordEndIDs map[int]struct{}
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		Children:   make(map[rune]*TrieNode),
		IDs:        make(map[int]struct{}),
		WordEndIDs: make(map[int]struct{}),
	}
}

type Trie struct {
	Root            *TrieNode
	MinPrefixLength int
	mu              sync.RWMutex
}

func NewTrie(minPrefixLength int) *Trie {
	return &Trie{Root: NewTrieNode(), MinPrefixLength: minPrefixLength}
}

func (t *Trie) Insert(word string, id int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	node := t.Root
	// Always add ID to root for empty string case
	node.IDs[id] = struct{}{}

	for _, ch := range word {
		if _, ok := node.Children[ch]; !ok {
			node.Children[ch] = NewTrieNode()
		}
		node = node.Children[ch]
		node.IDs[id] = struct{}{} // Mark question ID at every node for partial match
	}
	node.WordEndIDs[id] = struct{}{}
}

func (t *Trie) SearchPrefix(prefix string) map[int]struct{} {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// Special case: empty prefix should return all IDs
	if prefix == "" {
		return t.Root.IDs
	}

	if len([]rune(prefix)) < t.MinPrefixLength {
		return nil
	}

	node := t.Root
	for _, ch := range prefix {
		if _, ok := node.Children[ch]; !ok {
			return nil
		}
		node = node.Children[ch]
	}
	return node.IDs
}

func (t *Trie) Delete(word string, id int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Convert string to a slice of runes once to safely handle Unicode.
	runes := []rune(word)

	// The ID must also be removed from the root's ID set.
	delete(t.Root.IDs, id)

	if len(runes) == 0 {
		delete(t.Root.WordEndIDs, id)
		return
	}

	// i: next index to check
	// node: current node
	// return: returning "true" if the current node should "be deleted by its parent"
	var dfs func(node *TrieNode, i int) bool
	dfs = func(node *TrieNode, i int) bool {
		if i == len(runes) {
			delete(node.WordEndIDs, id)
			// If node is an unused leaf node, delete it
			return len(node.Children) == 0 && len(node.WordEndIDs) == 0 && len(node.IDs) == 0
		}

		ch := runes[i]
		child, ok := node.Children[ch]
		if !ok {
			// If the word is not in the trie, do nothing
			return false
		}

		// Remove the ID from the child node as we traverse down the path.
		delete(child.IDs, id)

		// If child is an unused leaf node, delete it
		if shouldDelete := dfs(child, i+1); shouldDelete {
			delete(node.Children, ch)
		}

		// After deleting child, determine if the CURRENT node is now prunable
		return len(node.Children) == 0 && len(node.WordEndIDs) == 0 && len(node.IDs) == 0
	}
	dfs(t.Root, 0)
}
