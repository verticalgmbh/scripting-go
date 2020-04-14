package scripts

import (
	"fmt"
	"reflect"

	"github.com/spf13/cast"
)

type OperatorType int8
type OperatorClass int8

const (
	OP_Not OperatorType = iota
	OP_Neg
	OP_Com
	OP_Inc
	OP_Dec
	OP_Mul
	OP_Div
	OP_Mod
	OP_Add
	OP_Sub
	OP_BitAnd
	OP_BitOr
	OP_BitXor
	OP_Shl
	OP_Shr
	OP_Rol
	OP_Ror
	OP_And
	OP_Or
	OP_Xor
	OP_Assign
	OP_Less
	OP_LessEqual
	OP_Greater
	OP_GreaterEqual
	OP_Equal
	OP_NotEqual
	OP_Match
	OP_NotMatch
	OP_AddAssign
	OP_SubAssign
	OP_DivAssign
	OP_MulAssign
	OP_ModAssign
	OP_ShlAssign
	OP_ShrAssign
	OP_AndAssign
	OP_OrAssign
	OP_XorAssign
	OP_Lambda
)

const (
	OP_PreUnary OperatorClass = iota
	OP_PostUnary
	OP_Binary
)

// Operator operates on one or two tokens depending on class to produce a result
type Operator struct {
	Type  OperatorType
	Class OperatorClass
	LHS   Token
	RHS   Token
}

// Execute executes the operator
func (op *Operator) Execute(variables *Variables) (interface{}, error) {
	var value interface{}
	var err error

	switch op.Type {
	case OP_Neg:
		value, err = op.neg(variables)
	case OP_Equal:
		value, err = op.equal(variables)
	case OP_NotEqual:
		value, err = op.equal(variables)
		if err == nil {
			value = !value.(bool)
		}
	case OP_Less:
		value, err = op.less(variables)
	case OP_LessEqual:
		value, err = op.lessEqual(variables)
	case OP_Greater:
		value, err = op.lessEqual(variables)
		if err == nil {
			value = !value.(bool)
		}
	case OP_GreaterEqual:
		value, err = op.less(variables)
		if err == nil {
			value = !value.(bool)
		}
	case OP_Add:
		value, err = op.add(variables)
	case OP_Sub:
		value, err = op.sub(variables)
	case OP_Mul:
		value, err = op.mul(variables)
	case OP_Div:
		value, err = op.div(variables)
	default:
		return nil, fmt.Errorf("Operator '%v' not implemented", op.Type)
	}

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (op *Operator) equal(variables *Variables) (bool, error) {
	lhs, err := op.LHS.Execute(variables)
	if err != nil {
		return false, err
	}

	rhs, err := op.RHS.Execute(variables)
	if err != nil {
		return false, err
	}

	return fmt.Sprintf("%v", lhs) == fmt.Sprintf("%v", rhs), nil
}

func (op *Operator) less(variables *Variables) (bool, error) {
	lhs, err := op.LHS.Execute(variables)
	if err != nil {
		return false, err
	}

	rhs, err := op.RHS.Execute(variables)
	if err != nil {
		return false, err
	}

	switch lhs.(type) {
	case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		switch rhs.(type) {
		case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			lhs, err = castValue(lhs, CAST_INT)
			if err == nil {
				rhs, err = castValue(rhs, CAST_INT)
				if err == nil {
					return lhs.(int64) < rhs.(int64), nil
				}
			}
		case float32, float64:
			lhs, err = castValue(lhs, CAST_DOUBLE)
			if err == nil {
				rhs, err = castValue(rhs, CAST_DOUBLE)
				if err == nil {
					return lhs.(float64) < rhs.(float64), nil
				}
			}
		case string:
			lhs, err = castValue(lhs, CAST_INT)
			if err == nil {
				cmp, err := castValue(rhs, CAST_INT)
				if err != nil {
					cmp, err = castValue(rhs, CAST_STRING)
					return lhs.(string) < cmp.(string), nil
				}
				return lhs.(int64) < cmp.(int64), nil
			}
		default:
			return false, fmt.Errorf("Comparision not supported for '%v'", reflect.ValueOf(rhs))
		}
	case float32, float64:
		switch rhs.(type) {
		case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			lhs, err = castValue(lhs, CAST_DOUBLE)
			if err == nil {
				rhs, err = castValue(rhs, CAST_DOUBLE)
				if err == nil {
					return lhs.(float64) < rhs.(float64), nil
				}
			}
		case float32, float64:
			lhs, err = castValue(lhs, CAST_DOUBLE)
			if err == nil {
				rhs, err = castValue(rhs, CAST_DOUBLE)
				if err == nil {
					return lhs.(float64) < rhs.(float64), nil
				}
			}
		case string:
			lhs, err = castValue(lhs, CAST_DOUBLE)
			if err == nil {
				cmp, err := castValue(rhs, CAST_DOUBLE)
				if err != nil {
					cmp, err = castValue(rhs, CAST_STRING)
					return lhs.(string) < cmp.(string), nil
				}
				return lhs.(float64) < cmp.(float64), nil
			}
		default:
			return false, fmt.Errorf("Comparision not supported for '%v'", reflect.ValueOf(rhs))
		}
	case string:
		switch rhs.(type) {
		case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			rhs, err = castValue(rhs, CAST_INT)
			if err == nil {
				cmp, err := castValue(lhs, CAST_INT)
				if err != nil {
					cmp, err = castValue(lhs, CAST_STRING)
					return lhs.(string) < cmp.(string), nil
				}
				return lhs.(int64) < cmp.(int64), nil
			}
		case float32, float64:
			rhs, err = castValue(rhs, CAST_DOUBLE)
			if err == nil {
				cmp, err := castValue(lhs, CAST_DOUBLE)
				if err != nil {
					cmp, err = castValue(lhs, CAST_STRING)
					return lhs.(string) < cmp.(string), nil
				}
				return lhs.(float64) < cmp.(float64), nil
			}
		case string:
			return lhs.(string) < rhs.(string), nil
		default:
			return false, fmt.Errorf("Comparision not supported for '%v'", reflect.ValueOf(rhs))
		}
	default:
		return false, fmt.Errorf("Comparision not supported for '%v'", reflect.ValueOf(lhs))
	}

	return false, err
}

func (op *Operator) neg(variables *Variables) (interface{}, error) {
	lhs, err := op.LHS.Execute(variables)
	if err != nil {
		return false, err
	}

	switch v := lhs.(type) {
	case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		return -cast.ToInt64(v), nil
	case float32, float64:
		return -cast.ToFloat64(v), nil
	default:
		return false, fmt.Errorf("Negation not supported for '%v'", reflect.ValueOf(lhs))
	}
}

func (op *Operator) lessEqual(variables *Variables) (bool, error) {
	lhs, err := op.LHS.Execute(variables)
	if err != nil {
		return false, err
	}

	rhs, err := op.RHS.Execute(variables)
	if err != nil {
		return false, err
	}

	switch lhs.(type) {
	case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		switch rhs.(type) {
		case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			lhs, err = castValue(lhs, CAST_INT)
			if err == nil {
				rhs, err = castValue(rhs, CAST_INT)
				if err == nil {
					return lhs.(int64) <= rhs.(int64), nil
				}
			}
		case float32, float64:
			lhs, err = castValue(lhs, CAST_DOUBLE)
			if err == nil {
				rhs, err = castValue(rhs, CAST_DOUBLE)
				if err == nil {
					return lhs.(float64) <= rhs.(float64), nil
				}
			}
		case string:
			lhs, err = castValue(lhs, CAST_INT)
			if err == nil {
				cmp, err := castValue(rhs, CAST_INT)
				if err != nil {
					cmp, err = castValue(rhs, CAST_STRING)
					return lhs.(string) <= cmp.(string), nil
				}
				return lhs.(int64) <= cmp.(int64), nil
			}
		default:
			return false, fmt.Errorf("Comparision not supported for '%v'", reflect.ValueOf(rhs))
		}
	case float32, float64:
		switch rhs.(type) {
		case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			lhs, err = castValue(lhs, CAST_DOUBLE)
			if err == nil {
				rhs, err = castValue(rhs, CAST_DOUBLE)
				if err == nil {
					return lhs.(float64) <= rhs.(float64), nil
				}
			}
		case float32, float64:
			lhs, err = castValue(lhs, CAST_DOUBLE)
			if err == nil {
				rhs, err = castValue(rhs, CAST_DOUBLE)
				if err == nil {
					return lhs.(float64) <= rhs.(float64), nil
				}
			}
		case string:
			lhs, err = castValue(lhs, CAST_DOUBLE)
			if err == nil {
				cmp, err := castValue(rhs, CAST_DOUBLE)
				if err != nil {
					cmp, err = castValue(rhs, CAST_STRING)
					return lhs.(string) <= cmp.(string), nil
				}
				return lhs.(float64) <= cmp.(float64), nil
			}
		default:
			return false, fmt.Errorf("Comparision not supported for '%v'", reflect.ValueOf(rhs))
		}
	case string:
		switch rhs.(type) {
		case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			rhs, err = castValue(rhs, CAST_INT)
			if err == nil {
				cmp, err := castValue(lhs, CAST_INT)
				if err != nil {
					cmp, err = castValue(lhs, CAST_STRING)
					return lhs.(string) <= cmp.(string), nil
				}
				return lhs.(int64) <= cmp.(int64), nil
			}
		case float32, float64:
			rhs, err = castValue(rhs, CAST_DOUBLE)
			if err == nil {
				cmp, err := castValue(lhs, CAST_DOUBLE)
				if err != nil {
					cmp, err = castValue(lhs, CAST_STRING)
					return lhs.(string) <= cmp.(string), nil
				}
				return lhs.(float64) <= cmp.(float64), nil
			}
		case string:
			return lhs.(string) <= rhs.(string), nil
		default:
			return false, fmt.Errorf("Comparision not supported for '%v'", reflect.ValueOf(rhs))
		}
	default:
		return false, fmt.Errorf("Comparision not supported for '%v'", reflect.ValueOf(lhs))
	}

	return false, err
}

func (op *Operator) add(variables *Variables) (interface{}, error) {
	lhs, err := op.LHS.Execute(variables)
	if err != nil {
		return false, err
	}

	rhs, err := op.RHS.Execute(variables)
	if err != nil {
		return false, err
	}

	switch lhs.(type) {
	case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		switch rhs.(type) {
		case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			lhs, err = castValue(lhs, CAST_INT)
			if err == nil {
				rhs, err = castValue(rhs, CAST_INT)
				if err == nil {
					return lhs.(int64) + rhs.(int64), nil
				}
			}
		case float32, float64:
			lhs, err = castValue(lhs, CAST_DOUBLE)
			if err == nil {
				rhs, err = castValue(rhs, CAST_DOUBLE)
				if err == nil {
					return lhs.(float64) + rhs.(float64), nil
				}
			}
		case string:
			lhs, err = castValue(lhs, CAST_INT)
			if err == nil {
				cmp, err := castValue(rhs, CAST_INT)
				if err != nil {
					cmp, err = castValue(rhs, CAST_STRING)
					return lhs.(string) + cmp.(string), nil
				}
				return lhs.(int64) + cmp.(int64), nil
			}
		default:
			return false, fmt.Errorf("Comparision not supported for '%v'", reflect.ValueOf(rhs))
		}
	case float32, float64:
		switch rhs.(type) {
		case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			lhs, err = castValue(lhs, CAST_DOUBLE)
			if err == nil {
				rhs, err = castValue(rhs, CAST_DOUBLE)
				if err == nil {
					return lhs.(float64) + rhs.(float64), nil
				}
			}
		case float32, float64:
			lhs, err = castValue(lhs, CAST_DOUBLE)
			if err == nil {
				rhs, err = castValue(rhs, CAST_DOUBLE)
				if err == nil {
					return lhs.(float64) + rhs.(float64), nil
				}
			}
		case string:
			lhs, err = castValue(lhs, CAST_DOUBLE)
			if err == nil {
				cmp, err := castValue(rhs, CAST_DOUBLE)
				if err != nil {
					cmp, err = castValue(rhs, CAST_STRING)
					return lhs.(string) + cmp.(string), nil
				}
				return lhs.(float64) + cmp.(float64), nil
			}
		default:
			return false, fmt.Errorf("Comparision not supported for '%v'", reflect.ValueOf(rhs))
		}
	case string:
		switch rhs.(type) {
		case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			rhs, err = castValue(rhs, CAST_INT)
			if err == nil {
				cmp, err := castValue(lhs, CAST_INT)
				if err != nil {
					cmp, err = castValue(lhs, CAST_STRING)
					return lhs.(string) + cmp.(string), nil
				}
				return lhs.(int64) + cmp.(int64), nil
			}
		case float32, float64:
			rhs, err = castValue(rhs, CAST_DOUBLE)
			if err == nil {
				cmp, err := castValue(lhs, CAST_DOUBLE)
				if err != nil {
					cmp, err = castValue(lhs, CAST_STRING)
					return lhs.(string) + cmp.(string), nil
				}
				return lhs.(float64) + cmp.(float64), nil
			}
		case string:
			return lhs.(string) + rhs.(string), nil
		default:
			return false, fmt.Errorf("Comparision not supported for '%v'", reflect.ValueOf(rhs))
		}
	default:
		return false, fmt.Errorf("Comparision not supported for '%v'", reflect.ValueOf(lhs))
	}

	return false, err
}

func (op *Operator) sub(variables *Variables) (interface{}, error) {
	lhs, err := op.LHS.Execute(variables)
	if err != nil {
		return false, err
	}

	rhs, err := op.RHS.Execute(variables)
	if err != nil {
		return false, err
	}

	switch lhs.(type) {
	case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		switch rhs.(type) {
		case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			lhs, err = castValue(lhs, CAST_INT)
			if err == nil {
				rhs, err = castValue(rhs, CAST_INT)
				if err == nil {
					return lhs.(int64) - rhs.(int64), nil
				}
			}
		case float32, float64:
			lhs, err = castValue(lhs, CAST_DOUBLE)
			if err == nil {
				rhs, err = castValue(rhs, CAST_DOUBLE)
				if err == nil {
					return lhs.(float64) - rhs.(float64), nil
				}
			}
		case string:
			lhs, err = castValue(lhs, CAST_INT)
			if err == nil {
				cmp, err := castValue(rhs, CAST_INT)
				if err == nil {
					return lhs.(int64) - cmp.(int64), nil
				}
			}
		default:
			return false, fmt.Errorf("Comparision not supported for '%v'", reflect.ValueOf(rhs))
		}
	case float32, float64:
		switch rhs.(type) {
		case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			lhs, err = castValue(lhs, CAST_DOUBLE)
			if err == nil {
				rhs, err = castValue(rhs, CAST_DOUBLE)
				if err == nil {
					return lhs.(float64) - rhs.(float64), nil
				}
			}
		case float32, float64:
			lhs, err = castValue(lhs, CAST_DOUBLE)
			if err == nil {
				rhs, err = castValue(rhs, CAST_DOUBLE)
				if err == nil {
					return lhs.(float64) - rhs.(float64), nil
				}
			}
		case string:
			lhs, err = castValue(lhs, CAST_DOUBLE)
			if err == nil {
				cmp, err := castValue(rhs, CAST_DOUBLE)
				if err == nil {
					return lhs.(float64) - cmp.(float64), nil
				}
			}
		default:
			return false, fmt.Errorf("Comparision not supported for '%v'", reflect.ValueOf(rhs))
		}
	case string:
		switch rhs.(type) {
		case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			rhs, err = castValue(rhs, CAST_INT)
			if err == nil {
				cmp, err := castValue(lhs, CAST_INT)
				if err == nil {
					return lhs.(int64) - cmp.(int64), nil
				}
			}
		case float32, float64:
			rhs, err = castValue(rhs, CAST_DOUBLE)
			if err == nil {
				cmp, err := castValue(lhs, CAST_DOUBLE)
				if err != nil {
					return lhs.(float64) - cmp.(float64), nil
				}
			}
		default:
			return false, fmt.Errorf("Comparision not supported for '%v'", reflect.ValueOf(rhs))
		}
	default:
		return false, fmt.Errorf("Comparision not supported for '%v'", reflect.ValueOf(lhs))
	}

	return false, err
}

func (op *Operator) mul(variables *Variables) (interface{}, error) {
	lhs, err := op.LHS.Execute(variables)
	if err != nil {
		return false, err
	}

	rhs, err := op.RHS.Execute(variables)
	if err != nil {
		return false, err
	}

	switch lhs.(type) {
	case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		switch rhs.(type) {
		case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			lhs, err = castValue(lhs, CAST_INT)
			if err == nil {
				rhs, err = castValue(rhs, CAST_INT)
				if err == nil {
					return lhs.(int64) * rhs.(int64), nil
				}
			}
		case float32, float64:
			lhs, err = castValue(lhs, CAST_DOUBLE)
			if err == nil {
				rhs, err = castValue(rhs, CAST_DOUBLE)
				if err == nil {
					return lhs.(float64) * rhs.(float64), nil
				}
			}
		case string:
			lhs, err = castValue(lhs, CAST_INT)
			if err == nil {
				cmp, err := castValue(rhs, CAST_INT)
				if err == nil {
					return lhs.(int64) * cmp.(int64), nil
				}
			}
		default:
			return false, fmt.Errorf("Comparision not supported for '%v'", reflect.ValueOf(rhs))
		}
	case float32, float64:
		switch rhs.(type) {
		case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			lhs, err = castValue(lhs, CAST_DOUBLE)
			if err == nil {
				rhs, err = castValue(rhs, CAST_DOUBLE)
				if err == nil {
					return lhs.(float64) * rhs.(float64), nil
				}
			}
		case float32, float64:
			lhs, err = castValue(lhs, CAST_DOUBLE)
			if err == nil {
				rhs, err = castValue(rhs, CAST_DOUBLE)
				if err == nil {
					return lhs.(float64) * rhs.(float64), nil
				}
			}
		case string:
			lhs, err = castValue(lhs, CAST_DOUBLE)
			if err == nil {
				cmp, err := castValue(rhs, CAST_DOUBLE)
				if err == nil {
					return lhs.(float64) * cmp.(float64), nil
				}
			}
		default:
			return false, fmt.Errorf("Comparision not supported for '%v'", reflect.ValueOf(rhs))
		}
	case string:
		switch rhs.(type) {
		case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			rhs, err = castValue(rhs, CAST_INT)
			if err == nil {
				cmp, err := castValue(lhs, CAST_INT)
				if err == nil {
					return lhs.(int64) * cmp.(int64), nil
				}
			}
		case float32, float64:
			rhs, err = castValue(rhs, CAST_DOUBLE)
			if err == nil {
				cmp, err := castValue(lhs, CAST_DOUBLE)
				if err != nil {
					return lhs.(float64) * cmp.(float64), nil
				}
			}
		default:
			return false, fmt.Errorf("Comparision not supported for '%v'", reflect.ValueOf(rhs))
		}
	default:
		return false, fmt.Errorf("Comparision not supported for '%v'", reflect.ValueOf(lhs))
	}

	return false, err
}

func (op *Operator) div(variables *Variables) (interface{}, error) {
	lhs, err := op.LHS.Execute(variables)
	if err != nil {
		return false, err
	}

	rhs, err := op.RHS.Execute(variables)
	if err != nil {
		return false, err
	}

	switch lhs.(type) {
	case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		switch rhs.(type) {
		case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			lhs, err = castValue(lhs, CAST_INT)
			if err == nil {
				rhs, err = castValue(rhs, CAST_INT)
				if err == nil {
					return lhs.(int64) / rhs.(int64), nil
				}
			}
		case float32, float64:
			lhs, err = castValue(lhs, CAST_DOUBLE)
			if err == nil {
				rhs, err = castValue(rhs, CAST_DOUBLE)
				if err == nil {
					return lhs.(float64) / rhs.(float64), nil
				}
			}
		case string:
			lhs, err = castValue(lhs, CAST_INT)
			if err == nil {
				cmp, err := castValue(rhs, CAST_INT)
				if err == nil {
					return lhs.(int64) / cmp.(int64), nil
				}
			}
		default:
			return false, fmt.Errorf("Comparision not supported for '%v'", reflect.ValueOf(rhs))
		}
	case float32, float64:
		switch rhs.(type) {
		case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			lhs, err = castValue(lhs, CAST_DOUBLE)
			if err == nil {
				rhs, err = castValue(rhs, CAST_DOUBLE)
				if err == nil {
					return lhs.(float64) / rhs.(float64), nil
				}
			}
		case float32, float64:
			lhs, err = castValue(lhs, CAST_DOUBLE)
			if err == nil {
				rhs, err = castValue(rhs, CAST_DOUBLE)
				if err == nil {
					return lhs.(float64) / rhs.(float64), nil
				}
			}
		case string:
			lhs, err = castValue(lhs, CAST_DOUBLE)
			if err == nil {
				cmp, err := castValue(rhs, CAST_DOUBLE)
				if err == nil {
					return lhs.(float64) / cmp.(float64), nil
				}
			}
		default:
			return false, fmt.Errorf("Comparision not supported for '%v'", reflect.ValueOf(rhs))
		}
	case string:
		switch rhs.(type) {
		case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			rhs, err = castValue(rhs, CAST_INT)
			if err == nil {
				cmp, err := castValue(lhs, CAST_INT)
				if err == nil {
					return lhs.(int64) / cmp.(int64), nil
				}
			}
		case float32, float64:
			rhs, err = castValue(rhs, CAST_DOUBLE)
			if err == nil {
				cmp, err := castValue(lhs, CAST_DOUBLE)
				if err != nil {
					return lhs.(float64) / cmp.(float64), nil
				}
			}
		default:
			return false, fmt.Errorf("Comparision not supported for '%v'", reflect.ValueOf(rhs))
		}
	default:
		return false, fmt.Errorf("Comparision not supported for '%v'", reflect.ValueOf(lhs))
	}

	return false, err
}
