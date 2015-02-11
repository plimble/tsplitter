package tsplitter

import (
	"strings"
	"sync"
)

const (
	knownType = iota
	unknownType
)

//Words after splited
type Words struct {
	words     []string
	wordTypes []int
	size      int
	deDup     []map[string]struct{}
	lockDedup sync.RWMutex
}

func newWords() *Words {
	return &Words{
		deDup: []map[string]struct{}{
			make(map[string]struct{}),
			make(map[string]struct{}),
		},
	}
}

func (w *Words) add(word string, wordType int) {
	w.words = append(w.words, word)
	w.wordTypes = append(w.wordTypes, wordType)
	w.addDedup(word, wordType)
	w.size++
}

func (w *Words) addDedup(word string, wordType int) {
	w.deDup[wordType][word] = struct{}{}
}

func (w *Words) removeDedup(word string, wordType int) {
	w.lockDedup.Lock()
	delete(w.deDup[wordType], word)
	w.lockDedup.Unlock()
}

func (w *Words) addKnown(word string) {
	w.add(word, knownType)
}

func (w *Words) addUnKnown(word string) {
	w.add(word, unknownType)
}

func (w *Words) concatLast(word string, newWordType int) {
	last := w.size - 1

	w.removeDedup(w.words[last], w.wordTypes[last])
	w.words[last] += word
	w.wordTypes[last] = newWordType
	w.addDedup(w.words[last], w.wordTypes[last])
}

func (w *Words) isLastType(wordType int) bool {
	return w.wordTypes[w.size-1] == wordType
}

//All return all words
func (w *Words) All() []string {
	return w.words
}

//AllDedup return all deduplicate words
func (w *Words) AllDedup() []string {
	allLen := 0
	for _, v := range w.deDup {
		allLen += len(v)
	}

	result := make([]string, allLen)
	i := 0
	for _, v := range w.deDup {
		for k := range v {
			result[i] = k
			i++
		}
	}

	return result
}

//AllDedupDelim return all deduplicate words with delimiter
func (w *Words) AllDedupDelim(delim string) string {
	allLen := 0
	for _, v := range w.deDup {
		allLen += len(v)
	}

	result := make([]string, allLen)
	i := 0
	for _, v := range w.deDup {
		for k := range v {
			result[i] = k
			i++
		}
	}

	return strings.Join(result, delim)
}

//AllDedupInterface return all deduplicate words in terface type
func (w *Words) AllDedupInterface() []interface{} {
	allLen := 0
	for _, v := range w.deDup {
		allLen += len(v)
	}

	result := make([]interface{}, allLen)
	i := 0
	for _, v := range w.deDup {
		for k := range v {
			result[i] = k
			i++
		}
	}

	return result
}

func (w *Words) getDedup(wordType int) []string {
	result := make([]string, len(w.deDup[wordType]))
	i := 0
	for k := range w.deDup[wordType] {
		result[i] = k
		i++
	}

	return result
}

//Unknown return deduplicate words which not found in dictionary
func (w *Words) Unknown() []string {
	return w.getDedup(unknownType)
}

//Known return deduplicate and ambiguous words which found in dictionary
func (w *Words) Known() []string {
	return w.getDedup(knownType)
}

//Size return size of words
func (w *Words) Size() int {
	return w.size
}
