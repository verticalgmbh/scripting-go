package scripts

// StatementBlock series of statements
type StatementBlock struct {
	Body     []Token
	IsMethod bool
}

// Execute executes the statement block
func (block *StatementBlock) Execute(variables *Variables) (interface{}, error) {
	blockvariables := NewVariables(variables)
	var result interface{}
	var err error

	for _, token := range block.Body {
		result, err = token.Execute(blockvariables)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
