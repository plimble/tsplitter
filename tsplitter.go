package tsplitter

import (
	"strings"
	"unicode/utf8"
)

func Split(dict Dictionary, str string) *Words {

	str = strings.Replace(str, "ํา", "ำ", -1)

	str = removeSpecialChar(str)
	sentences := chunkStrings(str)

	w := newWords()

	for _, sentence := range sentences {
		wordbreakLeftFirst(w, dict, sentence)
	}

	return w
}

func removeSpecialChar(str string) string {
	r := strings.NewReplacer(
		"!", "", "'", "", "‘", " ", "’", " ", "“", " ", "”", " ",
		"\"", " ", ")", " ", "(", " ", "{", " ", "}", " ",
		"...", "", "..", "", "…", "", ",", " ", ":", " ", "|", " ",
		"?", " ", "[", " ", "]", " ", "\\", " ", "\r", " ", "\r\n",
		" ", "\n", " ", "*", "", "\t", "", "|", " ", "/", " ", "+", " ", "ๆ", "",
		"~", " ", "=", " ", ">", " ", "<", " ",
	)

	return r.Replace(str)
}

func chunkStrings(str string) []string {
	return strings.Fields(str)
}

func isFrontDep(s rune) bool {
	switch s {
	case 'ะ', '้', 'า', 'ำ', 'ิ', 'ี', 'ึ', 'ื', 'ุ', 'ู', 'ๅ', '็', '์', 'ํ':
		return true
	}

	return false
}

func isRearDep(s rune) bool {
	switch s {
	case 'ั', 'ื', 'เ', 'แ', 'โ', 'ใ', 'ไ', 'ํ':
		return true
	}

	return false
}

func isTonal(s rune) bool {
	switch s {
	case '่', '้', '๊', '๋':
		return true
	}

	return false
}

func isEnding(s rune) bool {
	switch s {
	case 'ๆ', 'ฯ':
		return true
	}

	return false
}

func isThaiChar(ch rune) bool {
	return ch >= 'ก' && ch <= '๛' || ch == '.'
}

func isEnglish(ch rune) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')
}

func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9') || (ch >= '๑' && ch <= '๙')
}

func isSpecialChar(ch rune) bool {
	return ch <= '~' || ch == 'ๆ' || ch == 'ฯ' || ch == '“' || ch == '”' || ch == ','
}

func wordbreakLeftFirst(w *Words, dict Dictionary, sentence string) {
	match := 0
	pos := 0
	sentlen := len(sentence)

	for pos < sentlen {
		ch, size := utf8.DecodeRuneInString(sentence[pos:])

		pos += size

		if !isThaiChar(ch) {
			for pos < sentlen && !isThaiChar(ch) {
				ch, size = utf8.DecodeRuneInString(sentence[pos:])
				pos += size
			}
			if pos < sentlen {
				pos -= size
			}
			w.addKnown(sentence[match:pos])
			match = pos
		} else {
			pos = wordBreakLeft(w, dict, sentence, match)
			match = pos
		}
	}
}

func nextWordValid(dict Dictionary, beginPos int, sentence string) bool {
	pos := beginPos
	sentLen := len(sentence)

	if beginPos == sentLen {
		return true
	} else if ch, _ := utf8.DecodeRuneInString(sentence[beginPos:]); ch < '~' {
		return true
	} else {
		for pos < sentLen {
			_, size := utf8.DecodeRuneInString(sentence[pos:])
			pos += size

			if dict.Exist(sentence[beginPos:pos]) {
				return true
			}
		}
	}

	return false
}

func wordBreakLeft(w *Words, dict Dictionary, sentence string, beginPos int) int {
	pos := beginPos
	match := -1
	longestMatch := 0
	sentlen := len(sentence)
	numValidPos := 0
	nextBeginPos := beginPos
	var beginRune rune = 0
	var ch rune
	var size int
	var prevRune rune = 0

	for pos < sentlen {
		ch, size = utf8.DecodeRuneInString(sentence[pos:])

		if pos == beginPos {
			nextBeginPos += size
			beginRune = ch
		}

		pos += size
		if dict.Exist(sentence[beginPos:pos]) {
			match = pos
			if nextWordValid(dict, pos, sentence) {
				longestMatch = pos
				numValidPos++
			}
		}
	}

	if beginPos > 0 {
		prevRune, _ = utf8.DecodeLastRuneInString(sentence[:beginPos])
	}

	if match == -1 {
		if w.size > 0 && (isFrontDep(beginRune) || isTonal(beginRune) || isRearDep(prevRune) || w.isLastType(unknownType)) {
			w.concatLast(sentence[beginPos:nextBeginPos], unknownType)
		} else {
			w.addUnKnown(sentence[beginPos:nextBeginPos])
		}
		return nextBeginPos
	} else {
		if longestMatch == 0 {
			if isRearDep(prevRune) {
				w.concatLast(sentence[beginPos:match], unknownType)
			} else {
				w.add(sentence[beginPos:match], knownType)
			}
			return match
		} else {
			if isRearDep(prevRune) {
				w.concatLast(sentence[beginPos:longestMatch], unknownType)
			} else if numValidPos == 1 {
				w.add(sentence[beginPos:longestMatch], knownType)
			} else {
				w.add(sentence[beginPos:longestMatch], ambiguousType)
			}

			return longestMatch
		}
	}
}

// func mergeAmbiguous(dict Dictionary, prevRune rune, word string) string {
// 	newWord := string(prevRune) + word
// 	dict.Exist(newWord)
// 	if dict.Exist(newWord){
// 		return newWord
// 	}

// 	return ""

// 	maxWordlen := len(words) - 1

// 	for i := 0; i < maxWordlen; i++ {
// 		_, size := utf8.DecodeLastRuneInString(words[i])
// 		newWord := words[i][len(words[i])-size:] + words[i+1]
// 		if dict.Exist(newWord) {
// 			//merge
// 			newWord = words[i] + words[i+1]
// 			words = append(words[:i], words[i+2:]...)
// 			words = append(words, newWord)
// 			maxWordlen = len(words) - 1
// 		}
// 	}

// 	return words
// }
