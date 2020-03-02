package scripts

func isWhiteSpace(character byte) bool {
	switch character {
	case ' ', '\t', '\r', '\n':
		return true
	default:
		return false
	}
}

func peek(data *string, index int) byte {
	for index < len(*data) && isWhiteSpace((*data)[index]) {
		index++
	}

	if index < len(*data) {
		return (*data)[index]
	}

	return 0
}

func parseSpecialCharacter(character byte) byte {
	switch character {
	case 't':
		return '\t'
	case 'n':
		return '\n'
	case 'r':
		return '\r'
	default:
		return character
	}
}

func skipWhiteSpaces(data *string, index *int) {
	for *index < len(*data) {
		switch (*data)[*index] {
		case ' ', '\t', '\r', '\n':
			(*index)++
		default:
			return
		}
	}
}
