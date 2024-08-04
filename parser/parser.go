package parser

import (
	"fmt"
	"os"
	"strconv"

	"github.com/marzeq/mconf/tokeniser"
)

const (
  PARSER_VALUE_TYPE_STRING = "STRING"
  PARSER_VALUE_TYPE_NUMBER = "NUMBER"
  PARSER_VALUE_TYPE_BOOL   = "BOOL"
  PARSER_VALUE_TYPE_LIST   = "LIST"
  PARSER_VALUE_TYPE_OBJECT = "OBJECT"
)

type ParserValue interface {
  GetType() string

  GetString() string
  GetNumber() float64
  GetBool() bool
  GetList() []ParserValue
  GetObject() map[string]ParserValue
}

type Parser struct {
  tokens []tokeniser.Token
  currIndex int
  constants map[string]ParserValue
}

func NewParser(tokens []tokeniser.Token) Parser {
  return Parser{
    tokens: tokens,
    currIndex: 0,
    constants: make(map[string]ParserValue),
  }
}

func (p *Parser) PeekAhead(i int) tokeniser.Token {
  if p.currIndex+i >= len(p.tokens) {
    return tokeniser.EOFToken()
  }

  return p.tokens[p.currIndex+i]
}

func (p *Parser) Peek() tokeniser.Token {
  return p.PeekAhead(0)
}

func (p *Parser) Increment() {
  p.currIndex++
}

func (p *Parser) Consume() tokeniser.Token {
  t := p.Peek()

  p.Increment()

  return t
}

func (p *Parser) GoBack() {
  p.currIndex--
}

func (p *Parser) FormatErrorAtToken(message string, loc tokeniser.Location) string {
  return fmt.Sprintf("Parser error at line %d, col %d: %s\n", loc.Line, loc.Col, message)
}

func (p *Parser) ParseValue() ParserValue {
  token := p.Consume()

  switch token.Type {
    case tokeniser.TOKEN_TYPE_STRING:
      return &ParserValueString{Value: token.Value}
    case tokeniser.TOKEN_TYPE_NUMBER:
      converted, err := strconv.ParseFloat(token.Value, 64)

      if err != nil {
        fmt.Printf(p.FormatErrorAtToken("Failed to convert `%s` to number", token.Start))
        os.Exit(1)
      }

      return &ParserValueNumber{Value: converted}
    case tokeniser.TOKEN_TYPE_BOOL:
      var converted bool

      if token.Value == "true" {
        converted = true
      } else if token.Value == "false" {
        converted = false
      } else {
        fmt.Printf(p.FormatErrorAtToken("Failed to convert `%s` to bool", token.Start))
        os.Exit(1)
      }

      return &ParserValueBool{Value: converted}
    case tokeniser.TOKEN_TYPE_CONSTANT:
      value, ok := p.constants[token.Value]
      
      if !ok {
        fmt.Printf(p.FormatErrorAtToken("Constant `%s` not found", token.Start))
        os.Exit(1)
      }
      
      return value
    case tokeniser.TOKEN_TYPE_OPEN_LIST:
      fmt.Println("List values not implemented yet")
    case tokeniser.TOKEN_TYPE_OPEN_OBJ:
      fmt.Println("Object values not implemented yet")
    default:
      fmt.Printf(p.FormatErrorAtToken("Unexpected token", token.Start))
      os.Exit(1)
  }

  return nil
}

func (p *Parser) Parse() map[string]ParserValue {
  globalObject := make(map[string]ParserValue)

  for {
    token := p.Consume()

    switch token.Type {
      case tokeniser.TOKEN_TYPE_EOF:
        return globalObject
      case tokeniser.TOKEN_TYPE_KEY: fallthrough
      case tokeniser.TOKEN_TYPE_STRING: {
        key := token.Value

        assign := p.Consume()
        
        if assign.Type != tokeniser.TOKEN_TYPE_ASSIGN {
          fmt.Printf(p.FormatErrorAtToken("Expected assignment operator `=`", assign.Start))
          os.Exit(1)
        }

        value := p.ParseValue()
        
        globalObject[key] = value
      }
      case tokeniser.TOKEN_TYPE_CONSTANT: {
        key := token.Value
        
        assign := p.Consume()
        
        if assign.Type != tokeniser.TOKEN_TYPE_ASSIGN {
          fmt.Printf(p.FormatErrorAtToken("Expected assignment operator `=`", assign.Start))
          os.Exit(1)
        }

        value := p.ParseValue()

        p.constants[key] = value
      }
      case tokeniser.TOKEN_TYPE_OPEN_OBJ: {
        // this is a object with no key that's on the root. it's optional and you can have as many of them,
        // so you can use them to kinda group things together without having to use a key.
        //
        // effecively, the opening and closing tokens for this object will be ignored and
        // every key inside will be stuck onto the root object.
        
        /*
           see this example:

           ---
           {
             key1 = "value1"
             key2 = "value2"
           }
           
           { key3 = "value3" }
           ---

           is equivalent to:

           ---
           key1 = "value1"
           key2 = "value2"
           key3 = "value3"
           ---
        */

        fmt.Println("Object values not implemented yet")
      }
    }
  }

  return globalObject
}
