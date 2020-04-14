package scripts

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

// Parser parser used to parse expressions
type Parser struct {
	operators *OperatorTree
}

// NewParser creates a new expression parser
func NewParser(operators *OperatorTree) *Parser {
	return &Parser{operators: operators}
}

// Parse parses a script expression
func (parser *Parser) Parse(data string) (Token, error) {
	index := 0
	return parser.parseStatementBlock(nil, &data, &index, true)
}

func parseCharacter(data *string, index *int) (Token, error) {
	character := (*data)[*index]
	if character == '\\' {
		(*index)++
		character = parseSpecialCharacter((*data)[*index])
	}

	(*index)++
	if (*data)[*index] != '\'' {
		return nil, errors.New("Character Literal not terminated")
	}

	(*index)++
	return &Value{Value: character}, nil
}

func parseLiteral(data *string, index *int) (Token, error) {
	var literal strings.Builder

	for ; *index < len(*data); (*index)++ {
		character := (*data)[*index]
		switch character {
		case '"':
			(*index)++
			return &Value{Value: literal.String()}, nil
		case '\\':
			(*index)++
			literal.WriteByte(parseSpecialCharacter((*data)[*index]))
		default:
			literal.WriteByte(character)
		}
	}

	return nil, errors.New("Literal not terminated")
}

func parseNumber(token string) (interface{}, error) {
	token = strings.ToLower(token)
	if strings.HasPrefix(token, "0x") {
		if strings.HasSuffix(token, "ul") {
			return strconv.ParseUint(token[2:len(token)-2], 16, 64)
		}
		if strings.HasSuffix(token, "l") {
			return strconv.ParseInt(token[2:len(token)-1], 16, 64)
		}
		if strings.HasSuffix(token, "us") {
			return strconv.ParseUint(token[2:len(token)-2], 16, 16)
		}
		if strings.HasSuffix(token, "s") {
			return strconv.ParseInt(token[2:len(token)-1], 16, 16)
		}
		if strings.HasSuffix(token, "sb") {
			return strconv.ParseUint(token[2:len(token)-2], 16, 8)
		}
		if strings.HasSuffix(token, "b") {
			return strconv.ParseInt(token[2:len(token)-1], 16, 8)
		}
		if strings.HasSuffix(token, "u") {
			return strconv.ParseUint(token[2:len(token)-1], 16, 32)
		}
		return strconv.ParseInt(token[2:], 16, 32)
	}

	if strings.HasPrefix(token, "0o") {
		if strings.HasSuffix(token, "ul") {
			return strconv.ParseUint(token[2:len(token)-2], 8, 64)
		}
		if strings.HasSuffix(token, "l") {
			return strconv.ParseInt(token[2:len(token)-1], 8, 64)
		}
		if strings.HasSuffix(token, "us") {
			return strconv.ParseUint(token[2:len(token)-2], 8, 16)
		}
		if strings.HasSuffix(token, "s") {
			return strconv.ParseInt(token[2:len(token)-1], 8, 16)
		}
		if strings.HasSuffix(token, "sb") {
			return strconv.ParseUint(token[2:len(token)-2], 8, 8)
		}
		if strings.HasSuffix(token, "b") {
			return strconv.ParseInt(token[2:len(token)-1], 8, 8)
		}
		if strings.HasSuffix(token, "u") {
			return strconv.ParseUint(token[2:len(token)-1], 8, 32)
		}
		return strconv.ParseInt(token[2:], 8, 32)
	}

	if strings.HasPrefix(token, "0b") {
		if strings.HasSuffix(token, "ul") {
			return strconv.ParseUint(token[2:len(token)-2], 2, 64)
		}
		if strings.HasSuffix(token, "l") {
			return strconv.ParseInt(token[2:len(token)-1], 2, 64)
		}
		if strings.HasSuffix(token, "us") {
			return strconv.ParseUint(token[2:len(token)-2], 2, 16)
		}
		if strings.HasSuffix(token, "s") {
			return strconv.ParseInt(token[2:len(token)-1], 2, 16)
		}
		if strings.HasSuffix(token, "sb") {
			return strconv.ParseUint(token[2:len(token)-2], 2, 8)
		}
		if strings.HasSuffix(token, "b") {
			return strconv.ParseInt(token[2:len(token)-1], 2, 8)
		}
		if strings.HasSuffix(token, "u") {
			return strconv.ParseUint(token[2:len(token)-1], 2, 32)
		}
		return strconv.ParseInt(token[2:], 2, 32)
	}

	if strings.HasSuffix(token, "ul") {
		return strconv.ParseUint(token, 10, 64)
	}
	if strings.HasSuffix(token, "l") {
		return strconv.ParseInt(token, 10, 64)
	}
	if strings.HasSuffix(token, "us") {
		return strconv.ParseUint(token, 10, 16)
	}
	if strings.HasSuffix(token, "s") {
		return strconv.ParseInt(token, 10, 16)
	}
	if strings.HasSuffix(token, "sb") {
		return strconv.ParseUint(token, 10, 8)
	}
	if strings.HasSuffix(token, "b") {
		return strconv.ParseInt(token, 10, 8)
	}
	if strings.HasSuffix(token, "u") {
		return strconv.ParseUint(token, 10, 32)
	}

	dotcount := 0
	for i := 0; i < len(token); i++ {
		if token[i] == '.' {
			dotcount++
		}
	}

	switch dotcount {
	case 0:
		return strconv.ParseInt(token, 10, 32)
	case 1:
		if strings.HasSuffix(token, "d") {
			// this would be decimal but currently we don't use decimals
			return strconv.ParseFloat(token[:len(token)-1], 64)
		}
		if strings.HasSuffix(token, "f") {
			return strconv.ParseFloat(token[:len(token)-1], 32)
		}

		return strconv.ParseFloat(token, 64)
	default:
		// if no format triggers here this token is used as a string
		return &Value{Value: token}, nil
	}
}

func (parser *Parser) parseParameters(data *string, index *int) ([]Token, error) {
	skipWhiteSpaces(data, index)
	if *index >= len(*data) || (*data)[*index] != '(' {
		return nil, errors.New("Expected parameters")
	}

	(*index)++
	var parameters []Token
	for *index < len(*data) {
		switch (*data)[*index] {
		case ')', ']':
			(*index)++
			return parameters, nil
		case ',':
			(*index)++
		default:
			parameter, err := parser.parseTokenBlock(nil, data, index, false)
			if err != nil {
				return nil, err
			}
			parameters = append(parameters, parameter)
		}
	}

	return nil, errors.New("Parameter list not terminated")
}

func (parser *Parser) parseSingleParameter(data *string, index *int) (Token, error) {
	parameters, err := parser.parseParameters(data, index)
	if err != nil {
		return nil, err
	}

	switch len(parameters) {
	case 0:
		return nil, errors.New("A parameter was expected")
	case 1:
		return parameters[0], nil
	default:
		return nil, errors.New("More than one parameter encountered while only one was expected")
	}
}

func (parser *Parser) analyseToken(token string, data *string, index *int, startofstatement bool) (Token, error) {
	if len(token) == 0 {
		return nil, errors.New("Token expected")
	}

	if token[0] >= 0x30 && token[0] <= 0x39 {
		number, err := parseNumber(token)
		if err != nil {
			return nil, err
		}

		return &Value{Value: number}, nil
	}

	if startofstatement {
		// TODO: parse possible control tokens
	}

	switch token {
	case "bool", "int", "float", "double", "string":
		parameter, err := parser.parseSingleParameter(data, index)
		if err != nil {
			return nil, err
		}

		return &Cast{
			TargetType: token,
			Data:       parameter}, nil
	}

	switch token {
	case "true":
		return &Value{Value: true}, nil
	case "false":
		return &Value{Value: false}, nil
	case "null":
		return &Value{Value: nil}, nil
	}

	return &Variable{Name: token}, nil
}

func (parser *Parser) parseToken(data *string, index *int, startofstatement bool) (Token, error) {
	skipWhiteSpaces(data, index)

	parsenumber := false
	if *index < len(*data) {
		switch (*data)[*index] {
		case '.', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			parsenumber = true
		case '"':
			(*index)++
			return parseLiteral(data, index)
		case '\'':
			return parseCharacter(data, index)
		}
	}

	var tokenname strings.Builder
	for ; *index < len(*data); (*index)++ {
		character := (*data)[*index]
		if character >= 0x30 && character <= 0x39 || character >= 0x41 && character <= 0x5A || character >= 0x61 && character <= 0x7A || parsenumber && character == '.' {
			tokenname.WriteByte(character)
		} else if character == '"' || character == '\\' {
			(*index)++
			tokenname.WriteByte(parseSpecialCharacter((*data)[*index]))
		} else {
			return parser.analyseToken(tokenname.String(), data, index, startofstatement)
		}
	}

	if tokenname.Len() > 0 {
		return parser.analyseToken(tokenname.String(), data, index, startofstatement)
	}

	return &Value{}, nil
}

func (parser *Parser) parseInterpolation(data *string, index *int) (Token, error) {
	var tokens []Token
	var literal strings.Builder

	for ; *index < len(*data); (*index)++ {
		character := (*data)[*index]

		switch character {
		case '"':
			(*index)++
			tokens = append(tokens, &Value{Value: literal.String()})
			return &Interpolation{tokens: tokens}, nil
		case '{':
			(*index)++
			if peek(data, *index) == '{' {
				literal.WriteRune('{')
			} else {
				if literal.Len() > 0 {
					tokens = append(tokens, &Value{Value: literal.String()})
					literal.Reset()
				}

				block, err := parser.parseStatementBlock(nil, data, index, false)
				if err != nil {
					return nil, err
				}
				tokens = append(tokens, block)

				// for loop automatically increases index, but index should remain here
				(*index)--
			}
		case '\\':
			(*index)++
			literal.WriteByte(parseSpecialCharacter((*data)[*index]))
		default:
			literal.WriteByte(character)
		}
	}

	return nil, errors.New("Literal not terminated")
}

func (parser *Parser) parseBlock(data *string, index *int) (Token, error) {
	block, err := parser.parseTokenBlock(nil, data, index, false)
	if err != nil {
		return nil, err
	}

	if (*data)[*index] != ')' {
		return nil, errors.New("Block not terminated")
	}

	(*index)++
	return &StatementBlock{
		Body: []Token{
			block,
		},
		IsMethod: false,
	}, nil
}

func parseMember(host Token, data *string, index *int) (Token, error) {
	var membername strings.Builder

	for ; *index < len(*data); (*index)++ {
		character := (*data)[*index]
		if character >= 0x30 && character <= 0x39 || character >= 0x41 && character <= 0x5A || character >= 0x61 && character <= 0x7A || character == '_' {
			membername.WriteByte(character)
			continue
		}

		break
	}

	if membername.Len() > 0 {
		return &Member{host: host, member: membername.String()}, nil
	}

	return nil, errors.New("Membername expected")
}

func (parser *Parser) parseTokenBlock(parent Token, data *string, index *int, startofstatement bool) (Token, error) {
	skipWhiteSpaces(data, index)

	var tokens []Token
	var operators []*operatorIndex
	concat := true
	done := false
	start := *index

	for *index < len(*data) && !done {
		switch (*data)[*index] {
		case '=', '!', '~', '<', '>', '/', '+', '*', '-', '%', '&', '|', '^':
			operator, err := parser.operators.ParseOperator(data, index)
			if err != nil {
				return nil, err
			}

			if operator.Class == OP_PostUnary && !concat {
				(*index) -= 2
				done = true
				break
			}

			if operator.Type == OP_Sub {
				var isop bool
				if len(tokens) > 0 {
					_, isop = tokens[len(tokens)-1].(*Operator)
				}

				if len(tokens) == 0 || isop {
					operator = &Operator{Class: OP_PreUnary, Type: OP_Neg}
				}
			}

			operators = append(operators, &operatorIndex{index: len(tokens), operator: operator})
			tokens = append(tokens, operator)

			if !(operator.Class == OP_PreUnary || operator.Class == OP_PostUnary) {
				concat = true
			}
		case '.':
			member, err := parseMember(tokens[len(tokens)-1], data, index)
			if err != nil {
				return nil, err
			}
			tokens[len(tokens)-1] = member
		case '$':
			if !concat {
				done = true
				break
			}

			(*index)++
			if peek(data, *index) == '"' {
				(*index)++
				str, err := parser.parseInterpolation(data, index)
				if err != nil {
					return nil, err
				}
				tokens = append(tokens, str)
			} else {
				token, err := parser.parseToken(data, index, startofstatement)
				if err != nil {
					return nil, err
				}

				if _, ok := token.(*Variable); !ok {
					return nil, errors.New("Variable was expected")
				}

				tokens = append(tokens, token)
			}
			concat = false
		case '(':
			(*index)++
			block, err := parser.parseBlock(data, index)
			if err != nil {
				return nil, err
			}

			tokens = append(tokens, block)
			concat = false
		case '[':
			// TODO: parse indexer
			return nil, errors.New("Array indexer not yet implemented")
		case ',', ']', '}', ')':
			done = true
		default:
			if !concat {
				done = true
				break
			}

			token, err := parser.parseToken(data, index, startofstatement)
			if err != nil {
				return nil, err
			}

			// TODO: check for return and control tokens
			tokens = append(tokens, token)
		}

		if *index == start && !done {
			return nil, errors.New("Unable to parse code")
		}

		if !done {
			skipWhiteSpaces(data, index)
		}
	}

	if len(tokens) > 1 {
		sort.SliceStable(operators, func(i, k int) bool {
			return operators[i].operator.Type < operators[k].operator.Type
		})

		for _, opindex := range operators {
			switch opindex.operator.Class {
			case OP_PostUnary:
				if opindex.index == 0 {
					return nil, errors.New("Posttoken at beginning of tokenlist")
				}

				opindex.operator.LHS = tokens[opindex.index-1]
				tokens = removeAt(tokens, opindex.index-1)
				adjustOperatorIndices(operators, opindex.index, 1)
			case OP_PreUnary:
				opindex.operator.LHS = tokens[opindex.index+1]
				tokens = removeAt(tokens, opindex.index+1)
				adjustOperatorIndices(operators, opindex.index-1, 1)
			case OP_Binary:
				if opindex.index == 0 {
					return nil, errors.New("Left Hand Side operand expected")
				}
				if opindex.index >= len(tokens)-1 {
					return nil, errors.New("Right Hand Side operand expected")
				}
				opindex.operator.LHS = tokens[opindex.index-1]
				opindex.operator.RHS = tokens[opindex.index+1]
				tokens = removeAt(tokens, opindex.index+1)
				tokens = removeAt(tokens, opindex.index-1)
				adjustOperatorIndices(operators, opindex.index-1, 2)
			default:
				return nil, errors.New("Invalid operator class")
			}
		}
	}

	if len(tokens) > 1 {
		return nil, errors.New("Too many tokens left to make up a meaningful statement")
	}

	return tokens[0], nil
}

func adjustOperatorIndices(operators []*operatorIndex, index int, count int) {
	for _, opindex := range operators {
		if opindex.index > index {
			opindex.index -= count
		}
	}
}
func removeAt(tokens []Token, index int) []Token {
	return append(tokens[:index], tokens[index+1:]...)
}

func (parser *Parser) parseStatementBlock(parent Token, data *string, index *int, methodblock bool) (Token, error) {
	skipWhiteSpaces(data, index)

	var statements []Token

	terminated := false
	for *index < len(*data) {
		if (*data)[*index] == '}' {
			*index++
			terminated = true
			break
		}

		token, err := parser.parseTokenBlock(parent, data, index, true)
		if err != nil {
			return nil, err
		}

		if token != nil {
			statements = append(statements, token)
		}

		skipWhiteSpaces(data, index)
	}

	if !terminated && !methodblock {
		return nil, errors.New("Unterminated statement block")
	}

	if len(statements) <= 1 {
		return &StatementBlock{
			Body:     statements,
			IsMethod: methodblock}, nil
	}

	// TODO: fetch control tokens as soon as they are supported
	return &StatementBlock{
		Body:     statements,
		IsMethod: methodblock}, nil
}
