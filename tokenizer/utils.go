package tokenizer

import "unicode"

func isAlphaNumeric(b rune) bool {
	return unicode.IsLetter(b) || unicode.IsDigit(b) || b == '_'
}

func isWhitespace(b rune) bool {
	return b == '\t' || b == '\n' || b == ' ' || b == '\r'
}
