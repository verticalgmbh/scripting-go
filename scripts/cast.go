package scripts

import (
	"fmt"
	"strconv"
)

const CAST_BOOL = "bool"
const CAST_INT = "int"
const CAST_FLOAT = "float"
const CAST_DOUBLE = "double"
const CAST_STRING = "string"

// Cast casts/converts data to another type
//
// Some of these conversion implementation are pretty hacky and not really optimized for
// performance. The main goal was to get it to work fast.
type Cast struct {
	TargetType string
	Data       Token
}

func castValue(value interface{}, targettype string) (interface{}, error) {
	switch targettype {
	case "bool":
		switch value.(type) {
		case bool:
			return value, nil
		case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			return fmt.Sprintf("%d", value) != "0", nil
		case float32, float64:
			return fmt.Sprintf("%f", value) != "0", nil
		default:
			return value != nil, nil
		}
	case "int":
		switch value.(type) {
		case bool:
			bvalue := value.(bool)
			if bvalue {
				return int64(1), nil
			}
			return int64(0), nil
		case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			return strconv.ParseInt(fmt.Sprintf("%d", value), 10, 32)
		case float32, float64:
			return strconv.ParseInt(fmt.Sprintf("%f", value), 10, 32)
		case string:
			return strconv.ParseInt(fmt.Sprintf("%s", value), 10, 32)
		default:
			return nil, fmt.Errorf("Unsupported cast from '%v' to '%s'", value, targettype)
		}
	case "float":
		switch value.(type) {
		case bool:
			bvalue := value.(bool)
			if bvalue {
				return float32(1.0), nil
			}
			return float32(0.0), nil
		case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			return strconv.ParseFloat(fmt.Sprintf("%d", value), 32)
		case float32, float64:
			return strconv.ParseFloat(fmt.Sprintf("%f", value), 32)
		case string:
			return strconv.ParseFloat(fmt.Sprintf("%s", value), 32)
		default:
			return nil, fmt.Errorf("Unsupported cast from '%v' to '%s'", value, targettype)
		}
	case "double":
		switch value.(type) {
		case bool:
			bvalue := value.(bool)
			if bvalue {
				return 1.0, nil
			}
			return 0.0, nil
		case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			return strconv.ParseFloat(fmt.Sprintf("%d", value), 64)
		case float32, float64:
			return strconv.ParseFloat(fmt.Sprintf("%f", value), 64)
		case string:
			return strconv.ParseFloat(fmt.Sprintf("%s", value), 64)
		default:
			return nil, fmt.Errorf("Unsupported cast from '%v' to '%s'", value, targettype)
		}
	case "string":
		return fmt.Sprintf("%v", value), nil
	default:
		return nil, fmt.Errorf("Unsupported cast target type '%s'", targettype)
	}
}

// Execute executes the type cast
func (cast *Cast) Execute(variables *Variables) (interface{}, error) {
	value, err := cast.Data.Execute(variables)

	if err != nil {
		return nil, err
	}

	return castValue(value, cast.TargetType)
}
