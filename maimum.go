package tsplitter

import (
	"strings"
	"sync"
	"unicode/utf8"
)

var lockResult sync.RWMutex

func maximumMatch(dict Dictionary, str string) []string {
	result := []string{}

	str = removeSpecialChar(str)
	sentences := chunkStrings(str)

	for _, sentence := range sentences {
		segments := searchSegments(dict, sentence)

		result = append(result, segments...)
	}

	return result
}

func mapToArrayString(mapString map[string]struct{}) []string {
	arr := make([]string, len(mapString))
	index := 0
	for str, _ := range mapString {
		arr[index] = str
		index++
	}

	return arr
}

func addSegments(result map[string]struct{}, segments []string) {
	lockResult.Lock()
	for _, segment := range segments {
		result[segment] = struct{}{}
	}
	lockResult.Unlock()
}

func removeSpecialChar(str string) string {
	r := strings.NewReplacer(
		"!", "", "'", "", "‘", " ", "’", " ", "“", " ", "”", " ",
		"\"", " ", "-", " ", ")", " ", "(", " ", "{", " ", "}", " ",
		"...", "", "..", "", "…", "", ",", " ", ":", " ", "|", " ",
		"?", " ", "[", " ", "]", " ", "\\", " ", "\r", " ", "\r\n",
		" ", "\n", " ", "*", "", "\t", "",
	)

	return r.Replace(str)
}

func chunkStrings(str string) []string {
	return strings.Fields(str)
}

func searchSegments(dict Dictionary, sentence string) []string {
	segments := searchLeft(dict, bestSegment{sentence: sentence})

	resultIndex := 0

	//search no unknow
	index := []int{}
	for i, v := range segments {
		if v.unknowCount == 0 {
			index = append(index, i)
		}
	}

	//search no result no unknown
	if len(index) > 0 {
		min := 100000
		for _, v := range index {
			lenSeg := len(segments[v].Segments)
			if lenSeg < min {
				min = lenSeg
				resultIndex = v
			}
		}
	}

	// for i, v := range segments {
	//  fmt.Println("Posible", i, v.Segments, v.unknowCount == 0)
	// }

	// fmt.Println("Result:", resultIndex, segments[resultIndex].Segments)
	// fmt.Println("UnknownCount: ", segments[resultIndex].unknowCount)
	// fmt.Println("Len: ", len(segments[resultIndex].Segments))

	return segments[resultIndex].Segments
}

type bestSegment struct {
	sentence string
	Segments []string
	// Unknown     []string
	unknowCount int
}

func searchLeft(dict Dictionary, sourceSegment bestSegment) []bestSegment {
	var segment string
	var segments []bestSegment
	matchs := []string{}

	for len(sourceSegment.sentence) > 0 {
		ch, size := utf8.DecodeRuneInString(sourceSegment.sentence)
		sourceSegment.sentence = sourceSegment.sentence[size:]

		segment += string(ch)

		if dict.Exist(segment) {
			matchs = append(matchs, segment)
		}
	}

	if len(matchs) == 0 {
		copySourceSegment := CloneSegment(&sourceSegment)
		copySourceSegment.Segments = append(copySourceSegment.Segments, segment)
		copySourceSegment.unknowCount++

		segments = append(segments, copySourceSegment)
		return segments
	}

	for _, matched := range matchs {
		copySourceSegment := CloneSegment(&sourceSegment)

		copySourceSegment.Segments = append(copySourceSegment.Segments, matched)

		splited := splitMatchLeft(matched, segment)
		copySourceSegment.sentence = splited

		if splited != "" {
			resss := searchLeft(dict, copySourceSegment)
			segments = append(segments, resss...)
		} else {
			segments = append(segments, copySourceSegment)
		}
	}

	return segments
}

func CloneSegment(segment *bestSegment) bestSegment {
	newSegment := bestSegment{}

	newSegment.Segments = make([]string, len(segment.Segments))
	copy(newSegment.Segments, segment.Segments)

	return newSegment
}

func searchRight(dict Dictionary, sentence string) ([]string, bool) {
	var segment, matched string
	var segments []string
	isAllMatched := false

	for len(sentence) > 0 {
		ch, size := utf8.DecodeLastRuneInString(sentence)
		sentence = sentence[:len(sentence)-size]

		segment = string(ch) + segment
		if dict.Exist(segment) {
			matched = segment
		}
	}

	if matched != "" {
		isAllMatched = true
	}

	matched, segment = splitMatchRight(matched, segment)
	segments = append(segments, matched)

	if segment != "" {
		child := []string{}
		child, isAllMatched = searchRight(dict, segment)
		segments = append(segments, child...)
	}

	return segments, isAllMatched
}

func splitMatchRight(matched, segment string) (string, string) {
	matLen := len(matched)
	segLen := len(segment)

	if matched == "" {
		return segment, ""
	}

	if matLen == segLen {
		return matched, ""
	}

	splitStr := segment[:segLen-matLen]

	return matched, splitStr
}

func splitMatchLeft(matched, segment string) string {
	matLen := len(matched)

	if matLen == len(segment) {
		return ""
	}

	splitStr := segment[matLen:]

	return splitStr
}

func isThaiChar(ch rune) bool {
	if ch >= 'ก' && ch <= '๛' {
		return true
	}

	return false
}
