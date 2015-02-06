package tsplitter

import (

	// "fmt"
	"strings"
	"unicode/utf8"
)

type WordBreak struct {
	Words      []string
	Unknown    []string
	mapWord    map[string]struct{}
	mapUnknown map[string]struct{}
}

func (w *WordBreak) All() []string {
	return append(w.Words, w.Unknown...)
}

func (w *WordBreak) toArray() {
	for k, _ := range w.mapWord {
		w.Words = append(w.Words, k)
	}

	for k, _ := range w.mapUnknown {
		w.Unknown = append(w.Unknown, k)
	}
}

func Split(dict Dictionary, str string) *WordBreak {

	str = strings.Replace(str, "ํา", "ำ", -1)

	str = removeSpecialChar(str)
	sentences := chunkStrings(str)

	wb := &WordBreak{
		mapWord:    make(map[string]struct{}),
		mapUnknown: make(map[string]struct{}),
	}

	for _, sentence := range sentences {
		wordbreakLongest(dict, sentence, wb)
	}

	wb.toArray()

	return wb
}

func removeSpecialChar(str string) string {
	r := strings.NewReplacer(
		"!", "", "'", "", "‘", " ", "’", " ", "“", " ", "”", " ",
		"\"", " ", "-", " ", ")", " ", "(", " ", "{", " ", "}", " ",
		"...", "", "..", "", "…", "", ",", " ", ":", " ", "|", " ",
		"?", " ", "[", " ", "]", " ", "\\", " ", "\r", " ", "\r\n",
		" ", "\n", " ", "*", "", "\t", "", "|", " ", "/", " ", "+", " ",
	)

	return r.Replace(str)
}

func chunkStrings(str string) []string {
	return strings.Fields(str)
}

func isThaiChar(ch rune) bool {
	return ch >= 'ก' && ch <= '๛' || ch == '.'
}

func wordbreakLongest(dict Dictionary, sentence string, wb *WordBreak) {
	// words, unknown, _ := wordBreakRight(dict, sentence)
	for {
		words, unknown, _ := wordBreakRight(dict, sentence)
		wordsLen := len(words)
		switch {
		case wordsLen > 2:
			wordBreakLeft(dict, words[wordsLen-1]+unknown)
		}

		for _, v := range words {
			wb.mapWord[v] = struct{}{}
		}

		if unknown != "" {
			wb.mapUnknown[unknown] = struct{}{}
		}
	}
}

func wordBreakLeft(dict Dictionary, sentence string) ([]string, string, int) {
	var isThai bool = true
	var match, pos, lastMatch int = 0, 0, 0
	var fullSentence string = sentence
	words := []string{}

	for len(sentence) > 0 {
		ch, size := utf8.DecodeRuneInString(sentence)
		sentence = sentence[size:]

		pos += size

		if isThai && !isThaiChar(ch) {
			isThai = false
		}

		if isThai && dict.Exist(fullSentence[lastMatch:lastMatch+pos]) {
			match = pos
		}

		if len(sentence) == 0 {
			if match == 0 {
				return words, fullSentence[lastMatch:], lastMatch
			} else {
				words = append(words, fullSentence[lastMatch:lastMatch+match])
				lastMatch += match
				sentence = fullSentence[lastMatch:]
				match = 0
				pos = 0
				isThai = true
			}
		}
	}

	return words, "", lastMatch
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
				return words, fullSentence[maxPos-lastMatch-pos:], lastMatch
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

	return words, "", lastMatch
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
// 				fmt.Println("all match", matches)

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
