package utils

import "unicode/utf8"

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
