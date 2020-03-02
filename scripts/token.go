package scripts

// Token token in a script
type Token interface {

	// executes the token
	Execute(variables *Variables) (interface{}, error)
}
