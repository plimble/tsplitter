package tsplitter

const (
	noneType = iota
	knownType
	ambiguousType
	unknownType
)

type Words struct {
	words     []string
	wordTypes []int
	size      int

	knownDeDup     map[string]struct{}
	unknownDeDup   map[string]struct{}
	ambiguousDeDup map[string]struct{}
}

func newWords() *Words {
	return &Words{
		knownDeDup:     make(map[string]struct{}),
		unknownDeDup:   make(map[string]struct{}),
		ambiguousDeDup: make(map[string]struct{}),
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
	case ambiguousType:
		w.ambiguousDeDup[word] = struct{}{}
	}
}

func (w *Words) removeDedup(word string, wordType int) {
	switch wordType {
	case knownType:
		delete(w.knownDeDup, word)
	case unknownType:
		delete(w.unknownDeDup, word)
	case ambiguousType:
		delete(w.ambiguousDeDup, word)
	}
}

func (w *Words) addKnown(word string) {
	w.add(word, knownType)
}

func (w *Words) addUnKnown(word string) {
	w.add(word, unknownType)
}

func (w *Words) addAmbiguous(word string) {
	w.add(word, ambiguousType)
}

func (w *Words) concatLast(word string, newWordType int) {
	last := w.size - 1

	w.removeDedup(w.words[last], w.wordTypes[last])
	w.words[last] += word
	w.wordTypes[last] = newWordType
	w.addDedup(w.words[last], w.wordTypes[last])

	// old := w.words[last]
	// newWord := old + word
	// delete(w.keys, word)
	// fmt.Println(word, newWord, w.keys, w.words)
	// if _, has := w.keys[newWord]; !has {
	// 	w.words[last] = newWord
	// 	w.keys[newWord] = last
	// 	if newWordType != w.wordTypes[last] {
	// 		w.wordTypes[last].RemoveLast(old)
	// 		newWordType.Add(last, newWord)
	// 		w.wordTypes[last] = newWordType
	// 	}
	// } else {
	// 	w.wordTypes[last].RemoveLast(old)
	// 	w.words = w.words[:last]
	// 	w.size--
	// }
}

func (w *Words) isLastType(wordType int) bool {
	return w.wordTypes[w.size-1] == wordType
}

func (w *Words) All() []string {
	return w.words
}

func (w *Words) Ambiguous() []string {
	result := make([]string, len(w.ambiguousDeDup))
	i := 0
	for k, _ := range w.ambiguousDeDup {
		result[i] = k
		i++
	}

	return result
}

func (w *Words) Unknown() []string {
	result := make([]string, len(w.unknownDeDup))
	i := 0
	for k, _ := range w.unknownDeDup {
		result[i] = k
		i++
	}

	return result
}

func (w *Words) Known() []string {
	result := make([]string, len(w.knownDeDup))
	i := 0
	for k, _ := range w.knownDeDup {
		result[i] = k
		i++
	}

	return result
}

func (w *Words) Size() int {
	return w.size
}
