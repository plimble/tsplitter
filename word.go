package tsplitter

const (
	noneType = iota
	knownType
	ambiguousType
	unknownType
)

//Word after splited
type Words struct {
	words     []string
	wordTypes []int
	size      int

	knownDeDup   map[string]struct{}
	unknownDeDup map[string]struct{}
}

func newWords() *Words {
	return &Words{
		knownDeDup:   make(map[string]struct{}),
		unknownDeDup: make(map[string]struct{}),
	}
}

func (w *Words) add(word string, wordType int) {
	w.words = append(w.words, word)
	w.wordTypes = append(w.wordTypes, wordType)
	w.addDedup(word, wordType)
	w.size++
}

func (w *Words) addDedup(word string, wordType int) {
	switch wordType {
	case knownType:
		w.knownDeDup[word] = struct{}{}
	case unknownType:
		w.unknownDeDup[word] = struct{}{}
	}
}

func (w *Words) removeDedup(word string, wordType int) {
	switch wordType {
	case knownType:
		delete(w.knownDeDup, word)
	case unknownType:
		delete(w.unknownDeDup, word)
	}
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

//return all words
func (w *Words) All() []string {
	return w.words
}

//return deduplicate words which not found in dictionary
func (w *Words) Unknown() []string {
	result := make([]string, len(w.unknownDeDup))
	i := 0
	for k := range w.unknownDeDup {
		result[i] = k
		i++
	}

	return result
}

//return deduplicate and ambiguous words which found in dictionary
func (w *Words) Known() []string {
	result := make([]string, len(w.knownDeDup))
	i := 0
	for k := range w.knownDeDup {
		result[i] = k
		i++
	}

	return result
}

//return size of words
func (w *Words) Size() int {
	return w.size
}
