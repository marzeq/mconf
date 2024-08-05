package parser

import (
	"errors"
	"fmt"
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

  ValueToString() string

  GetString() (string, error)
  GetNumber() (float64, error)
  GetBool() (bool, error)
  GetList() ([]ParserValue, error)
  GetObject() (map[string]ParserValue, error)
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

func (p *Parser) FormatErrorAtToken(message string, loc tokeniser.Location) error {
  if loc.Line == 0 && loc.Col == 0 {
    return fmt.Errorf(fmt.Sprintf("Parser error at EOF: %s\n", message))
  }

  return fmt.Errorf(fmt.Sprintf("Parser error at line %d, col %d: %s\n", loc.Line, loc.Col, message))
}

func (p *Parser) ParseValue() (ParserValue, error) {
  token := p.Consume()

  switch token.Type {
    case tokeniser.TOKEN_TYPE_STRING:
      return &ParserValueString{Value: token.Value}, nil
    case tokeniser.TOKEN_TYPE_NUMBER:
      converted, err := strconv.ParseFloat(token.Value, 64)

      if err != nil {
        return nil, errors.Join(p.FormatErrorAtToken(fmt.Sprintf("Failed to convert `%s` to number", token.Value), token.Start), err)
      }

      return &ParserValueNumber{Value: converted}, nil
    case tokeniser.TOKEN_TYPE_BOOL:
      var converted bool

      if token.Value == "true" {
        converted = true
      } else if token.Value == "false" {
        converted = false
      } else {
        return nil, p.FormatErrorAtToken(fmt.Sprintf("Failed to convert `%s` to bool", token.Value), token.Start)
      }

      return &ParserValueBool{Value: converted}, nil
    case tokeniser.TOKEN_TYPE_CONSTANT:
      value, ok := p.constants[token.Value]
      
      if !ok {
        return nil, p.FormatErrorAtToken(fmt.Sprintf("Constant `%s` not found", token.Value), token.Start)
      }
      
      return value, nil
    case tokeniser.TOKEN_TYPE_OPEN_LIST:
      parsedList, err := p.ParseList()
      
      if err != nil { return nil, err }

      return &ParserValueList{Value: parsedList}, nil
    case tokeniser.TOKEN_TYPE_OPEN_OBJ:
      parsedObj, err := p.ParseObject()
      
      if err != nil { return nil, err }

      return &ParserValueObject{Value: parsedObj}, nil
    default:
      return nil, p.FormatErrorAtToken(fmt.Sprintf("Unexpected token %s", token.Type), token.Start)
  }

  return nil, fmt.Errorf("Unreachable code reached, please report this as a bug")
}

func (p *Parser) ParseList() ([]ParserValue, error) {
  list := make([]ParserValue, 0)

  for {
    token := p.Peek()

    switch token.Type {
      case tokeniser.TOKEN_TYPE_CLOSE_LIST:
        p.Increment()
        return list, nil
      case tokeniser.TOKEN_TYPE_STRING: fallthrough
      case tokeniser.TOKEN_TYPE_NUMBER: fallthrough
      case tokeniser.TOKEN_TYPE_BOOL: fallthrough
      case tokeniser.TOKEN_TYPE_OPEN_LIST: fallthrough
      case tokeniser.TOKEN_TYPE_CONSTANT: {
        value, err := p.ParseValue()

        if err != nil { return nil, err }

        comma_or_close := p.Peek()

        if comma_or_close.Type == tokeniser.TOKEN_TYPE_COMMA {
          p.Increment()
        } else if comma_or_close.Type != tokeniser.TOKEN_TYPE_CLOSE_LIST {
          return nil, p.FormatErrorAtToken("Expected comma or closing bracket", comma_or_close.Start)
        }

        list = append(list, value)
      }
      default: {
        return nil, p.FormatErrorAtToken("Unexpected token", token.Start)
      }
    }
  }

  return list, nil
}

func (p *Parser) ParseObject() (map[string]ParserValue, error) {
  object := make(map[string]ParserValue)

  for {
    token := p.Consume()

    switch token.Type {
      case tokeniser.TOKEN_TYPE_CLOSE_OBJ:
        return object, nil
      case tokeniser.TOKEN_TYPE_KEY: fallthrough
      case tokeniser.TOKEN_TYPE_STRING: {
        key := token.Value

        assign := p.Consume()
        
        if assign.Type != tokeniser.TOKEN_TYPE_ASSIGN {
          return nil, p.FormatErrorAtToken("Expected assignment operator `=`", assign.Start)
        }

        value, err := p.ParseValue()

        if err != nil { return nil, err }
        
        object[key] = value
        
        optional_comma := p.Peek()

        if optional_comma.Type == tokeniser.TOKEN_TYPE_COMMA {
          p.Increment()
        }
      }
      default: {
        return nil, p.FormatErrorAtToken(fmt.Sprintf("Unexpected token %s", token.Type), token.Start)
      }
    }
  }

  return object, nil
}

func (p *Parser) Parse() (map[string]ParserValue, error) {
  globalObject := make(map[string]ParserValue)

  for {
    token := p.Consume()

    switch token.Type {
      case tokeniser.TOKEN_TYPE_EOF:
        return globalObject, nil
      case tokeniser.TOKEN_TYPE_KEY: fallthrough
      case tokeniser.TOKEN_TYPE_STRING: {
        key := token.Value

        assign := p.Consume()
        
        if assign.Type != tokeniser.TOKEN_TYPE_ASSIGN {
          return nil, p.FormatErrorAtToken("Expected assignment operator `=`", assign.Start)
        }

        value, err := p.ParseValue()

        if err != nil { return nil, err }
        
        globalObject[key] = value
      }
      case tokeniser.TOKEN_TYPE_CONSTANT: {
        key := token.Value
        
        assign := p.Consume()
        
        if assign.Type != tokeniser.TOKEN_TYPE_ASSIGN {
          return nil, p.FormatErrorAtToken("Expected assignment operator `=`", assign.Start)
        }

        value, err := p.ParseValue()
        
        if err != nil { return nil, err }

        p.constants[key] = value
      }
      case tokeniser.TOKEN_TYPE_OPEN_OBJ: {
        object, err := p.ParseObject()
        
        if err != nil { return nil, err }

        for k, v := range object {
          globalObject[k] = v
        }
      }
      default: {
        return nil, p.FormatErrorAtToken(fmt.Sprintf("Unexpected token %s", token.Type), token.Start)
      }
    }
  }

  return globalObject, nil
}
