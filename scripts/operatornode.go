package scripts

// OperatorNode node in an operator tree
type OperatorNode struct {
	operator  OperatorType
	character rune
	children  map[rune]*OperatorNode
}

// AddChild adds a child to the node
func (node *OperatorNode) AddChild(character rune, child *OperatorNode) {
	node.children[character] = child
}

// GetChild get child of operator node
func (node *OperatorNode) GetChild(character rune) (*OperatorNode, bool) {
	value, exists := node.children[character]
	return value, exists
}

// HasChildren determines whether this node has child nodes
func (node *OperatorNode) HasChildren() bool {
	return len(node.children) > 0
}
