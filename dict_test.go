package tsplitter

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestFileDict(t *testing.T) {
	dict := NewFileDict("dictionary.txt")
	assert.True(t, len(dict.dict) > 70000)
}

func TestFileDictNotExist(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "open dictionary123.txt: no such file or directory", r.(error).Error())
		}
	}()
	NewFileDict("dictionary123.txt")
}

func TestFileDictExist(t *testing.T) {
	dict := NewFileDict("dictionary.txt")
	assert.True(t, dict.Exist("นั้น"))
	assert.False(t, dict.Exist("นั้น1"))
	assert.False(t, dict.Exist(""))
}
