package bayesian

// WordsTrie is the trie data structure which contains WordData for words.
type WordsTrie struct {
	WordData *WordData
	Children map[rune]*WordsTrie
}

// NewWordsTrie returns new WordsTrie.
func NewWordsTrie() WordsTrie {
	return WordsTrie{
		WordData: nil,
		Children: make(map[rune]*WordsTrie),
	}
}

// Get returns the WordData of the given word if it exists.
func (trie *WordsTrie) Get(word string) *WordData {
	node := trie
	for _, r := range word {
		node = node.Children[r]
		if node == nil {
			return nil
		}
	}
	return node.WordData
}

// GetOrNew returns the WordData of the given word.
// If WordData doesn't exist, make it and return.
func (trie *WordsTrie) GetOrNew(word string) *WordData {
	node := trie
	for _, r := range word {
		trie := NewWordsTrie()
		if node.Children[r] == nil {
			node.Children[r] = &trie
		}
		node = node.Children[r]
	}

	if node.WordData == nil {
		node.WordData = &WordData{
			Total: 0,
			Count: make(map[Class]int64),
		}
	}

	return node.WordData
}
