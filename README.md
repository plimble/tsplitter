tsplitter [![godoc badge](http://godoc.org/github.com/plimble/tsplitter?status.png)](http://godoc.org/github.com/plimble/tsplitter)   [![gocover badge](http://gocover.io/_badge/github.com/plimble/tsplitter)](http://gocover.io/github.com/plimble/tsplitter) [![Build Status](https://travis-ci.org/plimble/tsplitter.svg?branch=master)](https://travis-ci.org/plimble/tsplitter)
=========

Thai word break written in GO

### Installation
`go get -u github.com/plimble/tsplitter`

### Example

#####Get all words
```go
  func main(){
    dict := NewFileDict("dictionary.txt")
    txt := "ตัดคำไทย"
    words := Split(dict, txt)
    fmt.Println(words.All()) //ตัด, คำ, ไทย
  }
```

#####Get deduplicate word
```go
  func main(){
    dict := NewFileDict("dictionary.txt")
    txt := "ตัดคำไทย"
    words := Split(dict, txt)
    fmt.Println(words.Known())
    fmt.Println(words.Unknown())
  }
```

### Documentation
 - [GoDoc](http://godoc.org/github.com/plimble/tsplitter)