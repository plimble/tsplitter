package tsplitter

import (
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

func (w *WordBreak) AddWordAndUnknown(words []string, unknown string) {
	for _, v := range words {
		w.mapWord[v] = struct{}{}
	}

	if unknown != "" {
		w.mapUnknown[unknown] = struct{}{}
	}
}

func Split(dict Dictionary, str string) *WordBreak {

	str = strings.Replace(str, "ํา", "ำ", -1)

	str = removeSpecialChar(str)
	sentences := chunkStrings(str)

	w := &WordBreak{
		mapWord:    make(map[string]struct{}),
		mapUnknown: make(map[string]struct{}),
	}

	for _, sentence := range sentences {
		// fmt.Println(sentence)
		w.AddWordAndUnknown(wordbreakLeftFirst(dict, sentence))
	}

	w.toArray()

	return w
}

func removeSpecialChar(str string) string {
	r := strings.NewReplacer(
		"!", "", "'", "", "‘", " ", "’", " ", "“", " ", "”", " ",
		"\"", " ", "-", " ", ")", " ", "(", " ", "{", " ", "}", " ",
		"...", "", "..", "", "…", "", ",", " ", ":", " ", "|", " ",
		"?", " ", "[", " ", "]", " ", "\\", " ", "\r", " ", "\r\n",
		" ", "\n", " ", "*", "", "\t", "", "|", " ", "/", " ", "+", " ", "ๆ", "",
	)

	return r.Replace(str)
}

func chunkStrings(str string) []string {
	return strings.Fields(str)
}

func isThaiChar(ch rune) bool {
	return ch >= 'ก' && ch <= '๛' || ch == '.'
}

func wordbreakRightLeft(dict Dictionary, sentence string) ([]string, string) {
	lwords, lunknown, _ := wordBreakRight(dict, sentence)
	for index, word := range lwords {
		lwords = append(lwords[:index], lwords[index+1:]...)
		newwords, newunknown, _ := wordBreakLeft(dict, word)
		lwords = append(lwords, newwords...)
		lwords = append(lwords, newunknown)
	}

	if len(lunknown) > 10 {
		newwords, newunknown, _ := wordBreakLeft(dict, lunknown)
		lwords = append(lwords, newwords...)
		lunknown = newunknown
	}

	return lwords, lunknown
}

func wordbreakLeftFirst(dict Dictionary, sentence string) ([]string, string) {

	var rwords, lwords, xwords []string
	var runknown, lunknown, xunknown string
	for {
		lwords, lunknown, _ = wordBreakLeft(dict, sentence)

		lwordsLen := len(lwords)
		// fmt.Println("Left", lwords, lunknown)
		switch {
		case len(lunknown) == 3 && lwordsLen > 0:
			lwords[lwordsLen-1] = lwords[lwordsLen-1] + lunknown
			return lwords, ""
		case lunknown == "":
			return lwords, lunknown
		case lwordsLen == 0:
			rwords, runknown, _ = wordBreakRight(dict, lunknown)
		case lwordsLen > 2:
			word1 := lwords[lwordsLen-1:][0]
			word2 := lwords[lwordsLen-2:][0]
			lwords = lwords[:lwordsLen-2]
			// fmt.Println("Case Right", word2+word1+lunknown)
			rwords, runknown, _ = wordBreakRight(dict, word2+word1+lunknown)
		case lwordsLen == 2:
			word1 := lwords[lwordsLen-1:][0]
			lwords = lwords[:lwordsLen-1]
			// fmt.Println("Case Right", word1+lunknown)
			rwords, runknown, _ = wordBreakRight(dict, word1+lunknown)
		case lwordsLen == 1:
			rwords, runknown, _ = wordBreakRight(dict, lunknown)
		}
		// fmt.Println("Right", rwords, runknown)

		runknown = strings.Replace(runknown, ".", " ", -1)

		rwordLen := len(rwords)
		switch {
		case len(runknown) == 3 && rwordLen > 0:
			rwords[rwordLen-1] = runknown + rwords[lwordsLen-1]

			return rwords, ""
		case rwordLen == 0:
			return lwords, runknown
		case runknown == "":
			return lwords, runknown
		case len(lunknown) == len(runknown):
			return lwords, runknown
		case rwordLen > 2:
			word1 := rwords[rwordLen-1:][0]
			word2 := rwords[rwordLen-2:][0]
			lwords = append(lwords, rwords[:rwordLen-2]...)
			// fmt.Println("Merge", lwords, runknown+word1+word2)
			xwords, xunknown = wordbreakLeftFirst(dict, runknown+word1+word2)
		case rwordLen == 2:
			word1 := rwords[rwordLen-1:][0]
			lwords = append(lwords, rwords[:rwordLen-1]...)
			// fmt.Println("Merge", lwords, runknown+word1)
			xwords, xunknown = wordbreakLeftFirst(dict, runknown+word1)
		case rwordLen == 1:
			lwords = append(lwords, rwords[:rwordLen]...)
			xwords, xunknown = wordbreakLeftFirst(dict, runknown)
		}

		return append(lwords, xwords...), xunknown
	}
}

func wordbreakRightFirst(dict Dictionary, sentence string) ([]string, string) {

	var rwords, lwords, xwords []string
	var runknown, lunknown, xunknown string
	for {
		rwords, runknown, _ = wordBreakRight(dict, sentence)
		rwordsLen := len(rwords)
		// fmt.Println("Right", rwords, runknown)
		switch {
		case runknown == "":
			return rwords, runknown
		case rwordsLen == 0:
			lwords, lunknown, _ = wordBreakLeft(dict, runknown)
		case rwordsLen > 2:
			word1 := rwords[rwordsLen-1:][0]
			word2 := rwords[rwordsLen-2:][0]
			rwords = rwords[:rwordsLen-2]
			lwords, lunknown, _ = wordBreakLeft(dict, runknown+word1+word2)
		case rwordsLen == 2:
			word1 := rwords[rwordsLen-1:][0]
			rwords = rwords[:rwordsLen-1]
			// fmt.Println("Case Left", runknown+word1)
			lwords, lunknown, _ = wordBreakLeft(dict, runknown+word1)
		case rwordsLen == 1:
			// fmt.Println("Case Left", runknown)
			lwords, lunknown, _ = wordBreakLeft(dict, runknown)
		}
		// fmt.Println("Left", lwords, lunknown)

		lunknown = strings.Replace(lunknown, ".", " ", -1)

		lwordLen := len(lwords)
		switch {
		case lwordLen == 0:
			return rwords, lunknown
		case lunknown == "":
			return rwords, lunknown
		case len(lunknown) == len(runknown):
			return rwords, lunknown
		case lwordLen > 2:
			word1 := lwords[lwordLen-1:][0]
			word2 := lwords[lwordLen-2:][0]
			rwords = append(rwords, lwords[:lwordLen-2]...)
			// fmt.Println("Merge", rwords, word2+word1+lunknown)
			xwords, xunknown = wordbreakRightFirst(dict, word2+word1+lunknown)
		case lwordLen == 2:
			word1 := lwords[lwordLen-1:][0]
			rwords = append(rwords, lwords[:lwordLen-1]...)
			// fmt.Println("Merge", rwords, word1+lunknown)
			xwords, xunknown = wordbreakRightFirst(dict, word1+lunknown)
		case lwordLen == 1:
			rwords = append(rwords, lwords[:lwordLen]...)
			xwords, xunknown = wordbreakRightFirst(dict, lunknown)
		}

		return append(rwords, xwords...), xunknown
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
