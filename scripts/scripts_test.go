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

func Test_ParseEquation(t *testing.T) {
	parser := NewParser(NewExpressionOperators())
	script, err := parser.Parse("x*x+3=(y-x)*4")
	require.NoError(t, err)
	require.NotNil(t, script)
}

func Test_UnaryNeg(t *testing.T) {
	parser := NewParser(NewExpressionOperators())
	script, err := parser.Parse("-x")
	require.NoError(t, err)
	require.NotNil(t, script)

	body := script.(*StatementBlock)
	op := body.Body[0].(*Operator)
	require.Equal(t, OP_PreUnary, op.Class)
	require.Equal(t, OP_Neg, op.Type)
}

func Test_UnaryNegEquation(t *testing.T) {
	parser := NewParser(NewExpressionOperators())
	script, err := parser.Parse("x=-x")
	require.NoError(t, err)
	require.NotNil(t, script)

	body := script.(*StatementBlock)
	op := body.Body[0].(*Operator)
	require.Equal(t, OP_Binary, op.Class)
	require.Equal(t, OP_Assign, op.Type)

	rhsop := op.RHS.(*Operator)
	require.Equal(t, OP_PreUnary, rhsop.Class)
	require.Equal(t, OP_Neg, rhsop.Type)
}

func Test_UnaryNegExecution(t *testing.T) {
	parser := NewParser(NewExpressionOperators())
	script, err := parser.Parse("x-(-x)")
	require.NoError(t, err)
	require.NotNil(t, script)

	vars := NewVariables(nil)
	vars.SetVariable("x", 1.0)
	result, err := script.Execute(vars)
	require.NotNil(t, result)
}
