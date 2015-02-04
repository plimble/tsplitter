package tsplitter

import (
	"strings"
	"unicode/utf8"
)

type longestMatch struct {
	sentence      string
	leftSegments  map[string]struct{}
	rightSegments []string
	Unknown       []string
	match         string
	segment       string
	ch            rune
	size          int
	isThai        bool
	isSkip        bool
}

func Split(dict Dictionary, str string) []string {
	l := &longestMatch{}

	mapSegments := l.Match(dict, str)

	segments := make([]string, len(mapSegments))
	i := 0
	for k, _ := range mapSegments {
		segments[i] = k
		i++
	}

	return segments
}

func (m *longestMatch) Match(dict Dictionary, str string) map[string]struct{} {
	str = m.removeSpecialCharFast(str)
	m.leftSegments = make(map[string]struct{})

	for _, sentence := range m.chunkStrings(str) {
		m.sentence = sentence
		m.searchLeft(dict)
	}

	return m.leftSegments
}

func (m *longestMatch) chunkStrings(str string) []string {
	return strings.Fields(str)
}

func (m *longestMatch) removeSpecialCharFast(str string) string {
	r := strings.NewReplacer(
		"!", "", "'", "", "‘", " ", "’", " ", "“", " ", "”", " ",
		"\"", " ", "-", " ", ")", " ", "(", " ", "{", " ", "}", " ",
		"...", "", "..", "", "…", "", ",", " ", ":", " ", "|", " ",
		"?", " ", "[", " ", "]", " ", "\\", " ", "\r", " ", "\r\n",
		" ", "\n", " ", "*", "", "\t", "",
	)

	return r.Replace(str)
}

func (m *longestMatch) searchRight(dict Dictionary) {
	m.isThai = true

	for {
		for len(m.sentence) > 0 {
			m.ch, m.size = utf8.DecodeLastRuneInString(m.sentence)
			m.sentence = m.sentence[:len(m.sentence)-m.size]

			m.segment = string(m.ch) + m.segment
			if m.isThai && utf8.RuneLen(m.ch) < 3 {
				m.isThai = false
			}

			if m.isThai && dict.Exist(m.segment) {
				m.match = m.segment
			}

			if len(m.sentence) == 0 {
				if m.match == "" {
					m.rightSegments = append(m.rightSegments, maximumMatch(dict, m.segment)...)
					m.match = ""
					m.segment = ""
					m.isThai = true
					return
				} else {
					m.rightSegments = append(m.rightSegments, m.match)
				}
			}
		}

		m.sentence = m.segment[:len(m.segment)-len(m.match)]
		m.match = ""
		m.segment = ""
		m.isThai = true
		if m.sentence == "" {
			break
		}
	}

}

func (m *longestMatch) addLeftSegment(segments ...string) {
	for _, segment := range segments {
		m.leftSegments[segment] = struct{}{}
	}
}

func (m *longestMatch) searchLeft(dict Dictionary) {
	m.isThai = true
	fullSentence := m.sentence
	matches := []string{}

	for len(m.sentence) > 0 {
		m.ch, m.size = utf8.DecodeRuneInString(m.sentence)
		m.sentence = m.sentence[m.size:]

		m.segment += string(m.ch)
		if m.isThai && utf8.RuneLen(m.ch) < 3 {
			m.isThai = false
		}

		if m.isThai && dict.Exist(m.segment) {
			m.match = m.segment
		}

		if len(m.sentence) == 0 {
			if m.match == "" {
				matches = maximumMatch(dict, fullSentence)
				m.match = ""
				m.segment = ""
				m.isThai = true
				// m.leftSegments = append(m.leftSegments, maximumMatch(dict, m.segment)...)
				break
			} else {
				matches = append(matches, m.match)
				// m.addLeftSegment(m.match)
				m.sentence = m.segment[len(m.match):]
				m.match = ""
				m.segment = ""
				m.isThai = true
				// m.leftSegments = append(m.leftSegments, m.match)
			}
		}
	}
	m.addLeftSegment(matches...)
}
