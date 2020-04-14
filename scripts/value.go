package scripts

// Value value in a script
type Value struct {
	Value interface{}
}

// Execute returns the value of this token
func (value *Value) Execute(variables *Variables) (interface{}, error) {
	return value.Value, nil
}
