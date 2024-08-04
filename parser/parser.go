package parser

import (
  "fmt"
  "github.com/marzeq/mconf/tokeniser"
)

type Parser struct {
  tokens []tokeniser.Token
}

func NewParser(tokens []tokeniser.Token) Parser {
  return Parser{
    tokens: tokens,
  }
}

func (p *Parser) Parse() {
  fmt.Println("Parser.Parse: not implemented")
  for _, token := range p.tokens {
    fmt.Println(token.Repr())
  }
}
