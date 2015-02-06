package tsplitter

import (
	// "fmt"
	"strings"
	"unicode/utf8"
)

type bestSegment struct {
	sentence string
	Segments []string
	// Unknown     []string
	unknowCount int
	unknowLen   int
}

func maximumMatch(dict Dictionary, str string) []string {
	result := []string{}

	str = strings.Replace(str, ".", " ", -1)
	sentences := strings.Fields(str)

	// fmt.Println("")
	// fmt.Println("start", sentences)
	for index := 0; index < len(sentences); index++ {
		// if len(sentences[index]) > 150 {
		// 	sentences = append(sentences, sentences[index][150:])
		// 	sentences[index] = sentences[index][:150]
		// 	// fmt.Println("split", sentences[index], "|", sentences)
		// }

		segments := searchSegments(dict, sentences[index])
		// fmt.Println("Get", segments.Segments, len(segments.Segments), segments.unknowCount, segments.unknowLen)

		// if len(segments.Segments) == 1 {
		// 	result = append(result, sentences...)
		// 	break
		// }

		// if segments.unknowLen > 0 && index < len(sentences)-1 {
		// 	for i := 1; i <= segments.unknowLen; i++ {
		// 		sentences[index+1] = segments.Segments[len(segments.Segments)-i] + sentences[index+1]
		// 	}
		// 	// fmt.Println("merge", sentences[index+1], len(sentences[index+1]))
		// 	segments.Segments = segments.Segments[:len(segments.Segments)-segments.unknowLen]
		// }

		result = append(result, segments.Segments...)
	}

	// fmt.Println("Result", result)

	return result
}

func searchSegments(dict Dictionary, sentence string) *bestSegment {
	segments := searchLeft(dict, &bestSegment{sentence: sentence})

	oldUnknow := 1000
	oldMatch := 1000
	fullMatch := 1000
	resultIndex := 0
	hasUnknown := true

	for index, v := range segments {
		// fmt.Println("POS", v.Segments)
		if v.unknowCount == 0 {
			hasUnknown = false
			if len(v.Segments) < fullMatch {
				resultIndex = index
				fullMatch = len(v.Segments)
			}
		} else if hasUnknown {
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

	// fmt.Println("Result:", resultIndex, segments[resultIndex].Segments)
	// fmt.Println("UnknownCount: ", segments[resultIndex].unknowCount)
	// fmt.Println("Len: ", len(segments[resultIndex].Segments))

	return segments[resultIndex]
}

func searchLeft(dict Dictionary, sourceSegment *bestSegment) []*bestSegment {
	var segment string
	var segments []*bestSegment
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
		copySourceSegment := CloneSegment(sourceSegment)
		copySourceSegment.Segments = append(copySourceSegment.Segments, segment)
		copySourceSegment.unknowCount += len(segment)
		copySourceSegment.unknowLen++

		segments = append(segments, copySourceSegment)
		return segments
	}

	for _, matched := range matchs {
		copySourceSegment := CloneSegment(sourceSegment)

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

func CloneSegment(segment *bestSegment) *bestSegment {
	newSegment := &bestSegment{}

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
