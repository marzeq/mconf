package main

import (
	"fmt"
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

  parsed := p.Parse()

  for k, v := range parsed {
    if v.GetType() == parser.PARSER_VALUE_TYPE_STRING {
      fmt.Printf("%s: %s\n", k, v.GetString())
    } else if v.GetType() == parser.PARSER_VALUE_TYPE_NUMBER {
      fmt.Printf("%s: %f\n", k, v.GetNumber())
    } else if v.GetType() == parser.PARSER_VALUE_TYPE_BOOL {
      fmt.Printf("%s: %t\n", k, v.GetBool())
    } else if v.GetType() == parser.PARSER_VALUE_TYPE_LIST {
      fmt.Printf("%s: %v\n", k, v.GetList())
    } else if v.GetType() == parser.PARSER_VALUE_TYPE_OBJECT {
      fmt.Printf("%s: %v\n", k, v.GetObject())
    }
  }
}
