package search

type TrieNode struct {
	Children map[rune]*TrieNode
	IDs      map[int]struct{} // Questions having this prefix
	IsWord   bool
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		Children: make(map[rune]*TrieNode),
		IDs:      make(map[int]struct{}),
	}
}

type Trie struct {
	Root            *TrieNode
	MinPrefixLength int
}

func NewTrie(minPrefixLength int) *Trie {
	return &Trie{Root: NewTrieNode(), MinPrefixLength: minPrefixLength}
}

func (t *Trie) Insert(word string, id int) {
	node := t.Root
	for _, ch := range word {
		if _, ok := node.Children[ch]; !ok {
			node.Children[ch] = NewTrieNode()
		}
		node = node.Children[ch]
		node.IDs[id] = struct{}{} // Mark question ID at every node for partial match
	}
	node.IsWord = true
}

func (t *Trie) SearchPrefix(prefix string) map[int]struct{} {
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

	// i: next index to check
	// node: current node
	// return true if the node is a leaf node
	var dfs func(node *TrieNode, i int) bool
	dfs = func(node *TrieNode, i int) bool {
		if i == len(word) {
			delete(node.IDs, id)
			// If node is a leaf node, delete it
			return len(node.Children) == 0 && !node.IsWord
		}

		ch := rune(word[i])
		child, ok := node.Children[ch]
		if !ok {
			// If the word is not in the trie, do nothing
			return false
		}
		delete(child.IDs, id)

		// If child is a leaf node, delete it
		shouldDelete := dfs(child, i+1)
		if shouldDelete {
			delete(node.Children, ch)
		}

		// After deleting child, if the node is a leaf node, delete it
		return len(child.Children) == 0 && !child.IsWord
	}
	dfs(t.Root, 0)
}
