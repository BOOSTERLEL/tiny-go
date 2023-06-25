package lexer

import "unicode"

func isSpace(r rune) bool {
	switch r {
	case ' ', '\t', '\r':
		return true
	}
	return false
}

func isAlpha(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func isAlphaNumberic(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
