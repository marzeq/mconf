package main

import (
	// "fmt"
	"os"

	"github.com/marzeq/mconf/parser"
	"github.com/marzeq/mconf/tokeniser"
)

func check(e error) {
  if e != nil {
    panic(e)
  }
}

func ParseFromString(s string) (map[string]parser.ParserValue, error) {
  t := tokeniser.NewTokeniser(s)
  tokens, err := t.Tokenise()
  
  if err != nil { return nil, err }

  p := parser.NewParser(tokens)

  return p.Parse()
}

func ParseFromFile(filename string) (map[string]parser.ParserValue, error) {
  f, err := os.ReadFile(filename)
  
  if err != nil {
    return nil, err
  }

  s := string(f)

  return ParseFromString(s)
}

func main() {

}
