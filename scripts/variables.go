package scripts

import "fmt"

// Variables variable pool for token execution
type Variables struct {
	parent *Variables
	values map[string]interface{}
}

// NewVariables creates new variables
//
// **Parameters**
//   parent: parent provider where to look for variable values if no variable was defined
//           in the current provider
func NewVariables(parent *Variables) *Variables {
	return &Variables{
		parent: parent,
		values: make(map[string]interface{})}
}

// GetVariable get variable from provider
func (vars *Variables) GetVariable(name string) (interface{}, error) {
	value, err := vars.values[name]
	if err {
		return value, nil
	}

	if vars.parent != nil {
		return vars.parent.GetVariable(name)
	}

	return nil, fmt.Errorf("'%s' not found", name)
}

// SetVariable set variable value
func (vars *Variables) SetVariable(name string, value interface{}) {
	vars.values[name] = value
}
