/*
Thai word break written in GO

Example

Get all words

  func main(){
    dict := NewFileDict("dictionary.txt")
    txt := "ตัดคำไทย"
    words := Split(dict, txt)
    fmt.Println(words.All()) //ตัด, คำ, ไทย
  }

Get deduplicate word

  func main(){
    dict := NewFileDict("dictionary.txt")
    txt := "ตัดคำไทย"
    words := Split(dict, txt)
    fmt.Println(words.Known())
    fmt.Println(words.Unknown())
  }
*/
package tsplitter
