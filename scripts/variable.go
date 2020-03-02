package scripts

// Variable holds reference to a script variable
type Variable struct {
	Name string
}

// Execute returns value of variable
func (variable *Variable) Execute(variables *Variables) (interface{}, error) {
	return variables.GetVariable(variable.Name)
}
