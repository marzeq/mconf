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

func main() {
  file, err := os.ReadFile("example.mconf")

  check(err)

  contents := string(file)

  t := tokeniser.NewTokeniser(contents)

  tokens := t.Tokenise()

  p := parser.NewParser(tokens)

  p.Parse()
}
