package tsplitter

import (
	"bufio"
	"os"
)

type Dictionary interface {
	Exist(word string) bool
}

type FileDict struct {
	dict map[string]struct{}
}

func NewFileDict(filename string) *FileDict {
	fd := &FileDict{
		dict: make(map[string]struct{}),
	}

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		word := fd.removeBOM(sc.Bytes())
		fd.dict[string(word)] = struct{}{}
	}

	return fd
}

func (f *FileDict) removeBOM(word []byte) []byte {
	bom := []byte{239, 187, 191}
	if bom[0] == word[0] && bom[1] == word[1] && bom[2] == word[2] {
		return word[3:]
	}

	return word
}

func (f *FileDict) Exist(word string) bool {
	if word == "" {
		return false
	}

	_, ok := f.dict[word]
	return ok
}
