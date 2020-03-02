package scripts

import (
	"fmt"
	"strings"
)

// Interpolation combines tokens to a string
type Interpolation struct {
	tokens []Token
}

// Execute combines all tokens to a string
func (ip *Interpolation) Execute(variables *Variables) (interface{}, error) {
	var builder strings.Builder

	for _, token := range ip.tokens {
		value, err := token.Execute(variables)
		if err != nil {
			return nil, err
		}
		builder.WriteString(fmt.Sprintf("%v", value))
	}

	return builder.String(), nil
}
