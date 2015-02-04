package tsplitter

import (
	"unicode/utf8"
)

type bestSegment struct {
	sentence string
	Segments []string
	// Unknown     []string
	unknowCount int
}

func maximumMatch(dict Dictionary, str string) []string {
	return searchSegments(dict, str)
}

func searchSegments(dict Dictionary, sentence string) []string {
	segments := searchLeft(dict, bestSegment{sentence: sentence})

	oldUnknow := 1000
	oldMatch := 1000
	resultIndex := 0

	for index, v := range segments {
		if v.unknowCount == 0 {
			if len(v.Segments) < oldMatch {
				resultIndex = index
				oldMatch = len(v.Segments)
				oldUnknow = v.unknowCount
			}
		} else {
			if v.unknowCount == oldUnknow {
				if len(v.Segments) < oldMatch {
					resultIndex = index
					oldMatch = len(v.Segments)
					oldUnknow = v.unknowCount
				}
			} else if v.unknowCount < oldUnknow {
				resultIndex = index
				oldMatch = len(v.Segments)
				oldUnknow = v.unknowCount
			}
		}
	}

	// for i, v := range segments {
	// 	fmt.Println("Posible", i, v.Segments, v.unknowCount == 0, v.unknowCount)
	// }

	// fmt.Println("Result:", resultIndex, segments[resultIndex].Segments)
	// fmt.Println("UnknownCount: ", segments[resultIndex].unknowCount)
	// fmt.Println("Len: ", len(segments[resultIndex].Segments))

	return segments[resultIndex].Segments
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
		copySourceSegment.unknowCount += len(segment)

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

func splitMatchLeft(matched, segment string) string {
	matLen := len(matched)

	if matLen == len(segment) {
		return ""
	}

	splitStr := segment[matLen:]

	return splitStr
}

func isThaiChar(ch rune) bool {
	return ch >= 'ก' && ch <= '๛'
}
