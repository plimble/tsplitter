package tsplitter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWordsAdd(t *testing.T) {
	w := newWords()
	w.addKnown("123")
	assert.Equal(t, "123", w.words[0])
	assert.Equal(t, knownType, w.wordTypes[0])
	assert.Equal(t, 1, w.size)
	_, has := w.knownDeDup["123"]
	assert.True(t, has)
	assert.Len(t, w.knownDeDup, 1)

	w.addUnKnown("456")
	assert.Equal(t, "456", w.words[1])
	assert.Equal(t, unknownType, w.wordTypes[1])
	assert.Equal(t, 2, w.size)
	_, has = w.unknownDeDup["456"]
	assert.True(t, has)
	assert.Len(t, w.unknownDeDup, 1)

	w.addKnown("123")
	assert.Equal(t, "123", w.words[2])
	assert.Equal(t, knownType, w.wordTypes[2])
	assert.Equal(t, 3, w.size)
	_, has = w.knownDeDup["123"]
	assert.True(t, has)
	assert.Len(t, w.knownDeDup, 1)
}

func TestWordsConcatLast(t *testing.T) {
	w := newWords()
	w.addKnown("123")
	w.addKnown("456")

	w.concatLast("789", unknownType)

	assert.Equal(t, w.words[1], "456789")
	assert.Len(t, w.words, 2)
	assert.Len(t, w.knownDeDup, 1)
	assert.Len(t, w.unknownDeDup, 1)
	_, has := w.unknownDeDup["456789"]
	assert.True(t, has)
	_, has = w.knownDeDup["456"]
	assert.False(t, has)
	_, has = w.knownDeDup["123"]
	assert.True(t, has)
}

func TestWordsIsLastType(t *testing.T) {
	w := newWords()
	w.addKnown("123")
	w.addUnKnown("456")
	ok := w.isLastType(unknownType)
	assert.True(t, ok)
}

func TestWordsAll(t *testing.T) {
	w := newWords()
	w.addKnown("123")
	w.addKnown("456")
	w.addKnown("789")

	words := w.All()
	assert.Len(t, words, 3)
	assert.Equal(t, words[0], "123")
	assert.Equal(t, words[1], "456")
	assert.Equal(t, words[2], "789")
}

func TestWordsAllDeDup(t *testing.T) {
	w := newWords()
	w.addKnown("123")
	w.addKnown("123")
	w.addKnown("456")
	w.addUnKnown("abc")
	w.addUnKnown("abc")
	w.addUnKnown("def")

	words := w.AllDedup()
	assert.Len(t, words, 4)
}

func TestWordsKnown(t *testing.T) {
	w := newWords()
	w.addKnown("123")
	w.addKnown("123")
	w.addKnown("456")

	words := w.Known()
	assert.Len(t, words, 2)
}

func TestWordsUnKnown(t *testing.T) {
	w := newWords()
	w.addUnKnown("123")
	w.addUnKnown("123")
	w.addUnKnown("456")

	words := w.Unknown()
	assert.Len(t, words, 2)
}
