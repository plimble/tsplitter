package tsplitter

import (
	"bufio"
	"os"
)

//Dictionary interface
type Dictionary interface {
	Exist(word string) bool
}

//File System Dictionary
//File format should be one word per line
//word should not have any space
type FileDict struct {
	dict map[string]struct{}
}

// Create new file dictionary where filename is dictionary path
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

//Check word is exist in dictionary
func (f *FileDict) Exist(word string) bool {
	if word == "" {
		return false
	}

	_, ok := f.dict[word]
	return ok
}
