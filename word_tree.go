package main

type WordTreeNode struct {
	is_terminal bool
	letters     []WordTreeNode
}

func NewWordTreeNode() WordTreeNode {
	return WordTreeNode{false, nil}
}

func (root *WordTreeNode) insert(word string) {
	currentNode := root
	for _, char := range word {
		if currentNode.letters == nil {
			// Should be initialized to WordTreeNode{false, nil} by default.
			currentNode.letters = make([]WordTreeNode, 26)
		}
		currentNode = &currentNode.letters[char-'a']
	}
	currentNode.is_terminal = true
}

func (root *WordTreeNode) get_child(char byte) *WordTreeNode {
	if root.letters == nil {
		return nil
	}
	child := &root.letters[char-'a']
	if child.is_terminal || child.letters != nil {
		return child
	}
	return nil
}
