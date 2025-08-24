package mconf_parser

import (
  "fmt"
  "math/big"
  "os"
  "path"
  "strings"

  "github.com/marzeq/mconf/mconf_tokeniser"
  "github.com/marzeq/mconf/mconf_values"
)

type ValueType int

type Parser struct {
  tokens      []mconf_tokeniser.Token
  currIndex   int
  rootDir     string
  relativeDir string
  currentFile string
  Values      map[string]mconf_values.MconfValue
  ValuesOrder []string
  Constants   map[string]mconf_values.MconfValue
}

func NewParser(tokens []mconf_tokeniser.Token, rootDir string, currentFile string, relativeDir string, constants map[string]mconf_values.MconfValue) Parser {
	env := GetEnv()
	
	for k, v := range constants{
		env[k] = v
	}

  return Parser{
    tokens:      tokens,
    currIndex:   0,
    rootDir:     rootDir,
    relativeDir: relativeDir,
    currentFile: currentFile,
    Values:       make(map[string]mconf_values.MconfValue),
    ValuesOrder: []string{},
    Constants:   env,
  }
}

func GetEnv() map[string]mconf_values.MconfValue {
  env := make(map[string]mconf_values.MconfValue)

  for _, e := range os.Environ() {
    pair := strings.SplitN(e, "=", 2)
    env[pair[0]] = &mconf_values.MconfString{Value: pair[1]}
  }

  return env
}

func (p *Parser) PeekAhead(i int) mconf_tokeniser.Token {
  if p.currIndex+i >= len(p.tokens) {
    return mconf_tokeniser.EOFToken()
  }

  return p.tokens[p.currIndex+i]
}

func (p *Parser) Peek() mconf_tokeniser.Token {
  return p.PeekAhead(0)
}

func (p *Parser) Increment() {
  p.currIndex++
}

func (p *Parser) Consume() mconf_tokeniser.Token {
  t := p.Peek()

  p.Increment()

  return t
}

func (p *Parser) GoBack() {
  p.currIndex--
}

func (p *Parser) FormatErrorAtToken(message string, loc mconf_tokeniser.Location) error {
  var prettyFile string

  if p.currentFile == "" {
    prettyFile = "(stdin)"
  } else {
    prettyFile = path.Join(p.relativeDir, p.currentFile)
  }

  if loc.Line == 0 && loc.Col == 0 {
    return fmt.Errorf(fmt.Sprintf("%s (EOF) - Parser error: %s", prettyFile, message))
  }

  return fmt.Errorf(fmt.Sprintf("%s:%d:%d - Parser error: %s", prettyFile, loc.Line, loc.Col, message))
}

func (p *Parser) EvaluateStringValue(token mconf_tokeniser.Token) (string, error) {
  sb := ""
  for i, v := range token.Values {
    sb += v

    if i < len(token.StringSubs) {
      constantName := token.StringSubs[i]
      constantValue, ok := p.Constants[constantName]
      if !ok {
        return "", p.FormatErrorAtToken(fmt.Sprintf("Constant in string substitution `%s` not found", constantName), token.Start)
      }

      switch constantValue := constantValue.(type) {
      case *mconf_values.MconfString:
        constantStr := constantValue.Value
        sb += constantStr
      default:
        sb += constantValue.ValueToString()
      }
    }
  }

  return sb, nil
}

func (p *Parser) ParseConstantWithBackup() (mconf_values.MconfValue, error) {
  token := p.Consume()

  if token.Type == mconf_tokeniser.TOKEN_TYPE_CONSTANT {
    value, ok := p.Constants[token.Value]

    if ok {
      for {
        next := p.Peek()

        if next.Type == mconf_tokeniser.TOKEN_TYPE_QUESTION_MARK {
          p.Increment()

          unusedBackup := p.Consume()

          if unusedBackup.Type == mconf_tokeniser.TOKEN_TYPE_CONSTANT {
            continue
          } else {
            p.GoBack()
            _, err := p.ParseValue()
            if err != nil {
              return nil, err
            }
            break
          }
        } else {
          break
        }
      }

      return value, nil
    }

    peeked := p.Peek()

    if peeked.Type == mconf_tokeniser.TOKEN_TYPE_QUESTION_MARK {
      p.Increment()
      return p.ParseConstantWithBackup()
    }

    return nil, p.FormatErrorAtToken(fmt.Sprintf("Constant `%s` not found", token.Value), token.Start)
  } else {
    p.GoBack()
    return p.ParseValue()
  }
}

func (p *Parser) ParseValue() (mconf_values.MconfValue, error) {
  token := p.Consume()

  switch token.Type {
  case mconf_tokeniser.TOKEN_TYPE_WORD:
    return &mconf_values.MconfString{Value: token.Value}, nil
  case mconf_tokeniser.TOKEN_TYPE_STRING:
    sb, err := p.EvaluateStringValue(token)
    if err != nil {
      return nil, err
    }

    return &mconf_values.MconfString{Value: sb}, nil
  case mconf_tokeniser.TOKEN_TYPE_NUMBER_DECIMAL:
    if strings.Contains(token.Value, ".") || strings.Contains(token.Value, "e") || strings.Contains(token.Value, "E") {
      bigFl, _, err := big.ParseFloat(token.Value, 10, 0, big.ToNearestEven)
      if err != nil {
        return nil, p.FormatErrorAtToken(fmt.Sprintf("Failed to convert `%s` to float", token.Value), token.Start)
      }

      return &mconf_values.MconfFloat{Value: bigFl}, nil
    } else {
      intVal, success := new(big.Int).SetString(token.Value, 10)
      if !success {
        return nil, p.FormatErrorAtToken(fmt.Sprintf("Failed to convert `%s` to decimal int", token.Value), token.Start)
      }

      return &mconf_values.MconfInt{Value: intVal}, nil
    }
  case mconf_tokeniser.TOKEN_TYPE_NUMBER_HEX:
    intVal, success := new(big.Int).SetString(token.Value, 16)
    if !success {
      return nil, p.FormatErrorAtToken(fmt.Sprintf("Failed to convert `%s` to hex int", token.Value), token.Start)
    }

    return &mconf_values.MconfInt{Value: intVal}, nil
  case mconf_tokeniser.TOKEN_TYPE_NUMBER_BINARY:
    intVal, success := new(big.Int).SetString(token.Value, 2)
    if !success {
      return nil, p.FormatErrorAtToken(fmt.Sprintf("Failed to convert `%s` to binary int", token.Value), token.Start)
    }

    return &mconf_values.MconfInt{Value: intVal}, nil
  case mconf_tokeniser.TOKEN_TYPE_BOOL:
    var converted bool

		switch token.Value {
		case "true":
			converted = true
		case "false":
			converted = false
		default:
			return nil, p.FormatErrorAtToken(fmt.Sprintf("Failed to convert `%s` to bool", token.Value), token.Start)
		}

    return &mconf_values.MconfBool{Value: converted}, nil
  case mconf_tokeniser.TOKEN_TYPE_NULL:
    return &mconf_values.MconfNull{}, nil
  case mconf_tokeniser.TOKEN_TYPE_CONSTANT:
    p.GoBack()
    value, err := p.ParseConstantWithBackup()
    if err != nil {
      return nil, err
    }

    return value, nil
  case mconf_tokeniser.TOKEN_TYPE_OPEN_LIST:
    parsedList, err := p.ParseList()
    if err != nil {
      return nil, err
    }

    return &mconf_values.MconfList{Value: parsedList}, nil
  case mconf_tokeniser.TOKEN_TYPE_OPEN_OBJ:
    parsedObj, keysOrder, err := p.ParseObject()
    if err != nil {
      return nil, err
    }

    return &mconf_values.MconfObject{Value: parsedObj, KeysOrder: keysOrder}, nil
  default:
    return nil, p.FormatErrorAtToken(fmt.Sprintf("Unexpected token %s", token.Type), token.Start)
  }
}

func (p *Parser) ParseList() ([]mconf_values.MconfValue, error) {
  list := make([]mconf_values.MconfValue, 0)

  for {
    token := p.Peek()

    switch token.Type {
    case mconf_tokeniser.TOKEN_TYPE_CLOSE_LIST:
      p.Increment()
      return list, nil
    case mconf_tokeniser.TOKEN_TYPE_WORD:
      fallthrough
    case mconf_tokeniser.TOKEN_TYPE_STRING:
      fallthrough
    case mconf_tokeniser.TOKEN_TYPE_NUMBER_DECIMAL:
      fallthrough
    case mconf_tokeniser.TOKEN_TYPE_NUMBER_HEX:
      fallthrough
    case mconf_tokeniser.TOKEN_TYPE_NUMBER_BINARY:
      fallthrough
    case mconf_tokeniser.TOKEN_TYPE_BOOL:
      fallthrough
    case mconf_tokeniser.TOKEN_TYPE_NULL:
      fallthrough
    case mconf_tokeniser.TOKEN_TYPE_OPEN_LIST:
      fallthrough
    case mconf_tokeniser.TOKEN_TYPE_OPEN_OBJ:
      fallthrough
    case mconf_tokeniser.TOKEN_TYPE_CONSTANT:
      {
        value, err := p.ParseValue()
        if err != nil {
          return nil, err
        }

        comma_or_close := p.Peek()

        if comma_or_close.Type == mconf_tokeniser.TOKEN_TYPE_COMMA {
          p.Increment()
        } else if comma_or_close.Type != mconf_tokeniser.TOKEN_TYPE_CLOSE_LIST {
          return nil, p.FormatErrorAtToken("Expected comma or closing bracket", comma_or_close.Start)
        }

        list = append(list, value)
      }
    default:
      {
        return nil, p.FormatErrorAtToken(fmt.Sprintf("Unexpected token %s", token.Type), token.Start)
      }
    }
  }
}

func (p *Parser) ParseObject() (map[string]mconf_values.MconfValue, []string, error) {
  object := make(map[string]mconf_values.MconfValue)
  keysOrder := []string{}

  for {
    token := p.Consume()

    switch token.Type {
    case mconf_tokeniser.TOKEN_TYPE_CLOSE_OBJ:
      return object, keysOrder, nil
    case mconf_tokeniser.TOKEN_TYPE_WORD:
      fallthrough
    case mconf_tokeniser.TOKEN_TYPE_STRING:
      {
        var key string

        if token.Type == mconf_tokeniser.TOKEN_TYPE_WORD {
          key = token.Value
        } else {
          evkey, err := p.EvaluateStringValue(token)
          if err != nil {
            return nil, nil, err
          }

          key = evkey
        }

        assign := p.Consume()

        if assign.Type != mconf_tokeniser.TOKEN_TYPE_ASSIGN {
          return nil, nil, p.FormatErrorAtToken("Expected assignment operator `=`", assign.Start)
        }

        value, err := p.ParseValue()
        if err != nil {
          return nil, nil, err
        }

        object[key] = value
        keysOrder = append(keysOrder, key)

        optional_comma := p.Peek()

        if optional_comma.Type == mconf_tokeniser.TOKEN_TYPE_COMMA {
          p.Increment()
        }
      }
    default:
      {
        return nil, nil, p.FormatErrorAtToken(fmt.Sprintf("Unexpected token %s", token.Type), token.Start)
      }
    }
  }
}

func (p *Parser) Parse() (*mconf_values.MconfObject, error) {
  for {
    token := p.Consume()

    switch token.Type {
    case mconf_tokeniser.TOKEN_TYPE_EOF:
      return &mconf_values.MconfObject{Value: p.Values, KeysOrder: p.ValuesOrder}, nil
    case mconf_tokeniser.TOKEN_TYPE_WORD:
      fallthrough
    case mconf_tokeniser.TOKEN_TYPE_STRING:
      {
        var key string

        if token.Type == mconf_tokeniser.TOKEN_TYPE_WORD {
          key = token.Value
        } else {
          evkey, err := p.EvaluateStringValue(token)
          if err != nil {
            return nil,  err
          }

          key = evkey
        }

        assign := p.Consume()

        if assign.Type != mconf_tokeniser.TOKEN_TYPE_ASSIGN {
          return nil,  p.FormatErrorAtToken("Expected assignment operator `=`", assign.Start)
        }

        value, err := p.ParseValue()
        if err != nil {
          return nil,  err
        }

        p.Values[key] = value
        p.ValuesOrder = append(p.ValuesOrder, key)
      }
    case mconf_tokeniser.TOKEN_TYPE_CONSTANT:
      {
        key := token.Value

        assignOrQmark := p.Consume()
        if assignOrQmark.Type == mconf_tokeniser.TOKEN_TYPE_QUESTION_MARK {
          assign := p.Consume()
          if assign.Type != mconf_tokeniser.TOKEN_TYPE_ASSIGN {
            return nil,  p.FormatErrorAtToken("Expected assignment operator after `${key} ? at the top level`", assignOrQmark.Start)
          }
          value, err := p.ParseValue()
          if err != nil {
            return nil, err
          }
          if _, exists := p.Constants[key]; !exists {
            p.Constants[key] = value
          }
        } else {
          if assignOrQmark.Type != mconf_tokeniser.TOKEN_TYPE_ASSIGN {
            return nil, p.FormatErrorAtToken("Expected assignment operator `=`", assignOrQmark.Start)
          }

          value, err := p.ParseValue()
          if err != nil {
            return nil, err
          }

          p.Constants[key] = value
        }
      }
    case mconf_tokeniser.TOKEN_TYPE_OPEN_OBJ:
      {
        object, keysOrder, err := p.ParseObject()
        if err != nil {
          return nil, err
        }

        for _, k := range keysOrder {
          v := object[k]
          p.Values[k] = v
          p.ValuesOrder = append(p.ValuesOrder, k)
        }
      }
    case mconf_tokeniser.TOKEN_TYPE_DIRECTIVE:
      {
        switch token.Value {
        default:
          {
            return nil, p.FormatErrorAtToken(fmt.Sprintf("Unknown directive `%s`", token.Value), token.Start)
          }
        }
      }
    default:
      {
        return nil, p.FormatErrorAtToken(fmt.Sprintf("Unexpected token %s", token.Type), token.Start)
      }
    }
  }
}
