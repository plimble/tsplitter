package tsplitter

import (
	"strings"
	"unicode/utf8"
)

const (
	noneType = iota
	knownType
	ambiguousType
	unknownType
)

type WordBreak struct {
	Known     *OrderSet
	Ambiguous *OrderSet
	Unknown   *OrderSet
	lastType  int
}

func (w *WordBreak) All() []string {

	return append(w.Known.All(), append(w.Ambiguous.All(), w.Unknown.All()...)...)
}

func Split(dict Dictionary, str string) *WordBreak {

	str = strings.Replace(str, "ํา", "ำ", -1)

	str = removeSpecialChar(str)
	sentences := chunkStrings(str)

	w := &WordBreak{
		Known:     NewOrderSet(),
		Unknown:   NewOrderSet(),
		Ambiguous: NewOrderSet(),
	}

	for _, sentence := range sentences {
		// fmt.Println(sentence)
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

func removeSpecialCharSecond(str string) string {
	r := strings.NewReplacer(
		".", " ", "-", " ",
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

func wordbreakLeftFirst(w *WordBreak, dict Dictionary, sentence string) {
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
			w.Known.Add(sentence[match:pos])
			w.lastType = knownType
			match = pos
		} else {
			pos = wordBreakLeft(w, dict, sentence, match)
			match = pos

			// for pos < sentlen && isThaiChar(ch) {
			// 	ch, size = utf8.DecodeRuneInString(sentence[pos:])
			// 	pos += size
			// }
			// if pos < sentlen {
			// 	pos -= size
			// }
			// w.Known.Add(sentence[lastmatch:pos])
			// lastmatch = pos
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

func wordBreakLeft(w *WordBreak, dict Dictionary, sentence string, beginPos int) int {
	pos := beginPos
	match := 0
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
		ch, _ = utf8.DecodeLastRuneInString(sentence[:beginPos])
		prevRune = ch
	}

	if match == 0 {
		size = w.Unknown.Size() + w.Known.Size() + w.Ambiguous.Size()
		if size > 0 && (isFrontDep(beginRune) || isTonal(beginRune) || isRearDep(prevRune) || w.lastType == unknownType) {
			switch w.lastType {
			case unknownType:
				w.Unknown.ConcatLast(sentence[beginPos:nextBeginPos])
			case knownType:
				w.Unknown.Add(w.Known.RemoveLast() + sentence[beginPos:nextBeginPos])
			case ambiguousType:
				w.Unknown.Add(w.Ambiguous.RemoveLast() + sentence[beginPos:nextBeginPos])
			}
			w.lastType = unknownType
		} else {
			w.Unknown.Add(sentence[beginPos:nextBeginPos])
			w.lastType = unknownType
		}
		return nextBeginPos
	} else {
		if longestMatch == 0 {
			if isRearDep(prevRune) {
				switch w.lastType {
				case unknownType:
					w.Unknown.ConcatLast(sentence[beginPos:match])
				case knownType:
					w.Unknown.Add(w.Known.RemoveLast() + sentence[beginPos:match])
				case ambiguousType:
					w.Unknown.Add(w.Ambiguous.RemoveLast() + sentence[beginPos:match])
				}
				w.lastType = unknownType
			} else {
				w.Known.Add(sentence[beginPos:match])
				w.lastType = knownType
			}
			return match
		} else {
			if isRearDep(prevRune) {

				switch w.lastType {
				case unknownType:
					w.Unknown.ConcatLast(sentence[beginPos:longestMatch])
				case knownType:
					w.Unknown.Add(w.Known.RemoveLast() + sentence[beginPos:longestMatch])
				case ambiguousType:
					w.Unknown.Add(w.Ambiguous.RemoveLast() + sentence[beginPos:longestMatch])
				}
				w.lastType = unknownType
			} else if numValidPos == 1 {
				w.Known.Add(sentence[beginPos:longestMatch])
				w.lastType = knownType
			} else {
				w.Ambiguous.Add(sentence[beginPos:longestMatch])
				w.lastType = ambiguousType
			}

			return longestMatch
		}
	}
}

func wordBreakRight(dict Dictionary, sentence string) ([]string, string, int) {
	var isThai bool = true
	var match, pos, lastMatch, maxPos int = 0, 0, 0, len(sentence)
	var fullSentence string = sentence
	words := []string{}

	for len(sentence) > 0 {
		ch, size := utf8.DecodeLastRuneInString(sentence)
		sentence = sentence[:len(sentence)-size]

		pos += size

		if isThai && !isThaiChar(ch) {
			isThai = false
		}

		if isThai && dict.Exist(fullSentence[maxPos-lastMatch-pos:maxPos-lastMatch]) {
			match = pos
		}

		if len(sentence) == 0 {
			if match == 0 {
				words = mergeAmbiguous(dict, words)
				if !isThai {
					words = append(words, fullSentence[:maxPos-lastMatch])
					return words, "", maxPos
				}

				return words, fullSentence[:maxPos-lastMatch], lastMatch
			} else {
				words = append(words, fullSentence[maxPos-lastMatch-match:maxPos-lastMatch])
				lastMatch += match
				sentence = fullSentence[:maxPos-lastMatch]
				isThai = true
				match = 0
				pos = 0
			}
		}
	}

	words = mergeAmbiguous(dict, words)

	return words, "", lastMatch
}

func mergeAmbiguous(dict Dictionary, words []string) []string {
	maxWordlen := len(words) - 1
	for i := 0; i < maxWordlen; i++ {
		_, size := utf8.DecodeLastRuneInString(words[i])
		newWord := words[i][len(words[i])-size:] + words[i+1]
		if dict.Exist(newWord) {
			//merge
			newWord = words[i] + words[i+1]
			words = append(words[:i], words[i+2:]...)
			words = append(words, newWord)
			maxWordlen = len(words) - 1
		}
	}

	return words
}

func isLastCharSara(word string) bool {
	ch, _ := utf8.DecodeLastRuneInString(word)
	switch ch {
	case '์', 'ุ', 'ู', 'ึ', 'ะ', '๊', '็', '้', '่', 'า', 'ิ':
		return true
	}

	return false
}

// func (m *longestMatch) Match(dict Dictionary, str string) map[string]struct{} {
// 	str = m.removeSpecialCharFast(str)
// 	m.leftSegments = make(map[string]struct{})

// 	for _, sentence := range m.chunkStrings(str) {
// 		m.sentence = sentence
// 		m.searchLeft(dict)
// 	}

// 	return m.leftSegments
// }

// func (m *longestMatch) searchRight(dict Dictionary) []string {
// 	m.segment = ""
// 	m.match = ""
// 	m.isThai = true
// 	matches := []string{}

// 	for len(m.sentence) > 0 {
// 		m.ch, m.size = utf8.DecodeLastRuneInString(m.sentence)
// 		m.sentence = m.sentence[:len(m.sentence)-m.size]

// 		m.segment = string(m.ch) + m.segment

// 		if m.isThai && !isThaiChar(m.ch) {
// 			m.isThai = false
// 		}

// 		if m.isThai && dict.Exist(m.segment) {
// 			m.match = m.segment
// 		}

// 		if len(m.sentence) == 0 {
// 			if m.match == "" {
// 				// m.rightSegments = append(m.rightSegments, maximumMatch(dict, m.segment)...)
// 				matches = append(matches, m.segment)
// 				m.unknownRightCount = len(m.segment)
// 				m.unkniwnRightSegment = m.segment
// 				m.match = ""
// 				m.segment = ""
// 				m.isThai = true
// 				return matches
// 			} else {
// 				// m.rightSegments = append(m.rightSegments, m.match)
// 				matches = append(matches, m.match)
// 				m.sentence = m.segment[:len(m.segment)-len(m.match)]
// 				m.match = ""
// 				m.segment = ""
// 				m.isThai = true
// 			}
// 		}
// 	}

// 	return matches
// }

// func (m *longestMatch) addLeftSegment(segments ...string) {
// 	for _, segment := range segments {
// 		m.leftSegments[segment] = struct{}{}
// 	}
// }

// func (m *longestMatch) searchLeft(dict Dictionary) {
// 	m.isThai = true
// 	fullSentence := m.sentence
// 	fullSentenceSize := len(fullSentence)
// 	matches := []string{}

// 	for len(m.sentence) > 0 {
// 		m.ch, m.size = utf8.DecodeRuneInString(m.sentence)
// 		m.sentence = m.sentence[m.size:]

// 		m.segment += string(m.ch)
// 		if m.isThai && !isThaiChar(m.ch) {
// 			m.isThai = false
// 		}

// 		if m.isThai && dict.Exist(m.segment) {
// 			m.match = m.segment
// 		}

// 		if len(m.sentence) == 0 {
// 			if m.match == "" {
// fmt.Println("all match", matches)

// 				// if m.isThai {
// 				// 	if len(matches) > 0 {
// 				// 		m.segment = matches[len(matches)-1] + m.segment[len(m.match):]
// 				// 		matches = matches[:len(matches)-1]
// 				// 	}
// 				// 	matches = append(matches, maximumMatch(dict, m.segment)...)
// 				// } else {
// 				// 	matches = append(matches, m.segment)
// 				// }

// 				if m.isThai {

// 					// fmt.Println("Percent", (len(m.segment) * 100 / fullSentenceSize), "%", m.segment)
// 					if len(m.segment)*100/fullSentenceSize > 60 {
// 						// fmt.Println("go right", (len(m.segment) * 100 / fullSentenceSize), "%", m.segment)
// 						//search right
// 						m.sentence = fullSentence
// 						rMatch := m.searchRight(dict)
// 						// fmt.Println("Right Match", rMatch)

// 						if m.unknownRightCount == 0 {
// 							matches = rMatch
// 							m.match = ""
// 							m.segment = ""
// 							m.isThai = true
// 							break
// 						} else {
// 							rLen := len(rMatch)
// 							switch {
// 							case rLen > 2:
// 								m.segment = m.unkniwnRightSegment + rMatch[len(rMatch)-3]
// 								matches = rMatch[:len(rMatch)-3]
// 							case rLen == 2:
// 								m.segment = m.unkniwnRightSegment + rMatch[len(rMatch)-2]
// 								matches = rMatch[:len(rMatch)-2]
// 							case rLen == 1:
// 								m.segment = m.unkniwnRightSegment + rMatch[len(rMatch)-1]
// 								matches = rMatch[:len(rMatch)-1]
// 							}
// 						}
// 					} else {

// 						lLen := len(matches)
// 						switch {
// 						case lLen > 2:
// 							m.segment = matches[len(matches)-3] + m.segment[len(m.match):]
// 							matches = matches[:len(matches)-3]
// 						case lLen == 2:
// 							m.segment = matches[len(matches)-2] + m.segment[len(m.match):]
// 							matches = matches[:len(matches)-2]
// 						case lLen == 1:
// 							m.segment = matches[len(matches)-1] + m.segment[len(m.match):]
// 							matches = matches[:len(matches)-1]
// 						}
// 						fmt.Println("max left", m.segment, matches)
// 					}
// 					fmt.Println("GO Maximum", fullSentence, m.segment, matches)
// 					matches = append(matches, maximumMatch(dict, m.segment)...)
// 					// fmt.Println("Max", matches)
// 				} else {
// 					matches = append(matches, m.segment)
// 				}

// 				m.match = ""
// 				m.segment = ""
// 				m.isThai = true
// 				// m.leftSegments = append(m.leftSegments, maximumMatch(dict, m.segment)...)
// 				break
// 			} else {
// 				matches = append(matches, m.match)
// 				// if utf8.RuneCount([]byte(m.match)) > 1 {

// 				// }
// 				// m.addLeftSegment(m.match)
// 				m.sentence = m.segment[len(m.match):]
// 				m.match = ""
// 				m.segment = ""
// 				m.isThai = true
// 				// m.leftSegments = append(m.leftSegments, m.match)
// 			}
// 		}
// 	}
// 	m.addLeftSegment(matches...)
// }
