package scripts

import "errors"

// OperatorTree tree containing all supported operators
type OperatorTree struct {
	OperatorNode
}

// AddOperator adds an operator to the tree
func (tree *OperatorTree) AddOperator(literal string, operator OperatorType) {
	current := &tree.OperatorNode

	for _, character := range literal {
		node, found := current.GetChild(character)
		if !found {
			node = &OperatorNode{
				character: character,
				children:  make(map[rune]*OperatorNode, 0)}
			current.children[character] = node
		}

		current = node
	}

	current.operator = operator
}

// ParseOperator parses an operator from tree
func (tree *OperatorTree) ParseOperator(data *string, index *int) (*Operator, error) {
	current := &tree.OperatorNode
	parsestart := *index

loop:
	for *index < len(*data) {
		character := (*data)[*index]
		switch character {
		case '=', '!', '~', '<', '>', '/', '+', '-', '*', '%', '&', '|', '^':
			(*index)++
		default:
			break loop
		}

		child, found := current.GetChild(rune(character))
		if !found {
			(*index)--
			break
		}

		current = child
		if !current.HasChildren() {
			break
		}
	}

	if current == nil {
		return nil, errors.New("Operator expected but nothing found")
	}

	switch current.operator {
	case OP_Inc, OP_Dec:
		if *index-parsestart >= 3 && !isWhiteSpace((*data)[*index-3]) {
			return &Operator{Type: current.operator, Class: OP_PostUnary}, nil
		}
		if *index < len(*data) && !isWhiteSpace((*data)[*index]) {
			return &Operator{Type: current.operator, Class: OP_PreUnary}, nil
		}
		return nil, errors.New("Increment/Decrement without connected operand detected")
	case OP_Neg, OP_Not:
		return &Operator{Type: current.operator, Class: OP_PreUnary}, nil
	default:
		return &Operator{Type: current.operator, Class: OP_Binary}, nil
	}
}

// NewExpressionOperators creates a new operator tree containing all operators
//                        used for expression evaluation
func NewExpressionOperators() *OperatorTree {
	tree := &OperatorTree{}
	tree.children = make(map[rune]*OperatorNode, 0)

	tree.AddOperator("+", OP_Add)
	tree.AddOperator("-", OP_Sub)
	tree.AddOperator("*", OP_Mul)
	tree.AddOperator("/", OP_Div)
	tree.AddOperator("=", OP_Assign)
	tree.AddOperator("!", OP_Not)
	tree.AddOperator("++", OP_Inc)
	tree.AddOperator("--", OP_Dec)
	tree.AddOperator("%", OP_Mod)
	tree.AddOperator("<", OP_Less)
	tree.AddOperator("<=", OP_LessEqual)
	tree.AddOperator(">", OP_Greater)
	tree.AddOperator(">=", OP_GreaterEqual)
	tree.AddOperator("==", OP_Equal)
	tree.AddOperator("!=", OP_NotEqual)
	tree.AddOperator("~~", OP_Match)
	tree.AddOperator("!~", OP_NotMatch)
	tree.AddOperator("&", OP_BitAnd)
	tree.AddOperator("|", OP_BitOr)
	tree.AddOperator("^", OP_BitXor)
	tree.AddOperator("&&", OP_And)
	tree.AddOperator("||", OP_Or)
	tree.AddOperator("^^", OP_Xor)
	tree.AddOperator(">>", OP_Shr)
	tree.AddOperator(">>>", OP_Ror)
	tree.AddOperator("<<", OP_Shl)
	tree.AddOperator("<<<", OP_Rol)
	return tree
}
