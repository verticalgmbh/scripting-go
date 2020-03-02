package scripts

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Equal_Int_Int(t *testing.T) {
	parser := NewParser(NewExpressionOperators())
	script, err := parser.Parse("1022==1022")
	require.NoError(t, err)

	result, err := script.Execute(nil)
	require.NoError(t, err)
	require.Equal(t, true, result)
}

func Test_Less_Int_Int(t *testing.T) {
	parser := NewParser(NewExpressionOperators())
	script, err := parser.Parse("855<1022")
	require.NoError(t, err)

	result, err := script.Execute(nil)
	require.NoError(t, err)
	require.Equal(t, true, result)
}

func Test_ComplexComparision(t *testing.T) {
	parser := NewParser(NewExpressionOperators())
	script, err := parser.Parse("88.0+702>=4096/8.0")
	require.NoError(t, err)

	result, err := script.Execute(nil)
	require.NoError(t, err)
	require.Equal(t, true, result)
}

func Test_InterpolationPlusString(t *testing.T) {
	parser := NewParser(NewExpressionOperators())
	script, err := parser.Parse("$\"Hello {name}, \"+\"may i help you?\"")
	require.NoError(t, err)

	variables := NewVariables(nil)
	variables.SetVariable("name", "Gangolf")
	result, err := script.Execute(variables)
	require.NoError(t, err)
	require.Equal(t, "Hello Gangolf, may i help you?", result)
}
