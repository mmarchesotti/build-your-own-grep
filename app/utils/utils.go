package utils

import "unicode/utf8"

func IsDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func IsLower(r rune) bool {
	return r >= 'a' && r <= 'z'
}

func IsUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

func IsAlpha(r rune) bool {
	return IsUpper(r) || IsLower(r)
}

func IsAlphaNumeric(r rune) bool {
	return IsAlpha(r) || IsDigit(r) || r == '_'
}

func StartsWithRune(line []byte, charToMatch rune) bool {
	if len(line) == 0 {
		return false
	}

	firstRune, size := utf8.DecodeRune(line)

	if firstRune == utf8.RuneError && size <= 1 {
		return false
	}

	return firstRune == charToMatch
}
