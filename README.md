tsplitter [![godoc badge](http://godoc.org/github.com/plimble/tsplitter?status.png)](http://godoc.org/github.com/plimble/tsplitter)   [![gocover badge](http://gocover.io/_badge/github.com/plimble/tsplitter?t=2)](http://gocover.io/github.com/plimble/tsplitter) [![Build Status](https://api.travis-ci.org/plimble/tsplitter.svg?branch=master&t=2)](https://travis-ci.org/plimble/tsplitter) [![Go Report Card](http://goreportcard.com/badge/plimble/tsplitter?t=2)](http:/goreportcard.com/report/plimble/tsplitter)
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

###Contributing

If you'd like to help out with the project. You can put up a Pull Request.
