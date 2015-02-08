package tsplitter

import (
	"unicode/utf8"
)

type Strings string

func (s Strings) CharAtString(index int) (string, int) {
	ch, size := s.CharAtRune(index)
	return string(ch), size
}

func (s Strings) CharAtRune(index int) (rune, int) {
	if !utf8.RuneStart([]byte(s[index:])[0]) {
		panic("index is not rune")
	}

	ch, size := utf8.DecodeRuneInString(string(s[index:]))
	return ch, size
}

func (s Strings) Range(start, end int) string {
	return string(s[start:end])
}
