package tsplitter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileDict(t *testing.T) {
	dict := NewFileDict("dictionary.txt")
	assert.True(t, len(dict.dict) > 70000)
}

func TestFileDictExist(t *testing.T) {
	dict := NewFileDict("dictionary.txt")
	ok := dict.Exist("นั้น")

	assert.True(t, ok)
}

func BenchmarkFileDictExist(b *testing.B) {
	dict := NewFileDict("dictionary.txt")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dict.Exist("นั้น")
	}
}
