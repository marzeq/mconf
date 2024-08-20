package parser

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/marzeq/mconf/tokeniser"
)

const (
	PARSER_VALUE_TYPE_STRING = "STRING"
	PARSER_VALUE_TYPE_FLOAT  = "FLOAT"
	PARSER_VALUE_TYPE_INT    = "INT"
	PARSER_VALUE_TYPE_UINT   = "UINT"
	PARSER_VALUE_TYPE_BOOL   = "BOOL"
	PARSER_VALUE_TYPE_LIST   = "LIST"
	PARSER_VALUE_TYPE_OBJECT = "OBJECT"
)

type ParserValue interface {
	GetType() string

	ValueToString(indentAndDepth ...int) string

	GetString() (string, error)
	GetFloat() (float64, error)
	GetInt() (int64, error)
	GetUInt() (uint64, error)
	GetBool() (bool, error)
	GetList() ([]ParserValue, error)
	GetObject() (map[string]ParserValue, error)
}

type Parser struct {
	tokens      []tokeniser.Token
	currIndex   int
	constants   map[string]ParserValue
	rootDir     string
	currentFile string
}

func NewParser(tokens []tokeniser.Token, rootDir string, currentFile string) Parser {
	return Parser{
		tokens:      tokens,
		currIndex:   0,
		constants:   make(map[string]ParserValue),
		rootDir:     rootDir,
		currentFile: currentFile,
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
	var prettyFile string

	if p.currentFile == "" {
		prettyFile = "(stdin)"
	} else {
		prettyFile = p.currentFile
	}

	if loc.Line == 0 && loc.Col == 0 {
		return fmt.Errorf(fmt.Sprintf("%s - Parser error at EOF: %s\n", prettyFile, message))
	}

	return fmt.Errorf(fmt.Sprintf("%s - Parser error at line %d, col %d: %s\n", prettyFile, loc.Line, loc.Col, message))
}

func (p *Parser) ParseValue() (ParserValue, error) {
	token := p.Consume()

	switch token.Type {
	case tokeniser.TOKEN_TYPE_STRING:
		return &ParserValueString{Value: token.Value}, nil
	case tokeniser.TOKEN_TYPE_NUMBER:
		if strings.Contains(token.Value, ".") {
			converted, err := strconv.ParseFloat(token.Value, 0)
			if err != nil {
				return nil, errors.Join(p.FormatErrorAtToken(fmt.Sprintf("Failed to convert `%s` to number", token.Value), token.Start), err)
			}

			return &ParserValueFloat{Value: converted}, nil
		} else if strings.Contains(token.Value, "-") {
			converted, err := strconv.ParseInt(token.Value, 10, 0)
			if err != nil {
				if err.(*strconv.NumError).Err == strconv.ErrRange {
					converted = math.MaxInt64
				} else {
					return nil, errors.Join(p.FormatErrorAtToken(fmt.Sprintf("Failed to convert `%s` to number", token.Value), token.Start), err)
				}
			}

			return &ParserValueInt{Value: converted}, nil
		} else {
			converted, err := strconv.ParseUint(token.Value, 10, 0)
			if err != nil {
				if err.(*strconv.NumError).Err == strconv.ErrRange {
					converted = math.MaxUint64
				} else {
					return nil, errors.Join(p.FormatErrorAtToken(fmt.Sprintf("Failed to convert `%s` to number", token.Value), token.Start), err)
				}
			}

			return &ParserValueUInt{Value: converted}, nil
		}
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
		if err != nil {
			return nil, err
		}

		return &ParserValueList{Value: parsedList}, nil
	case tokeniser.TOKEN_TYPE_OPEN_OBJ:
		parsedObj, err := p.ParseObject()
		if err != nil {
			return nil, err
		}

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
		case tokeniser.TOKEN_TYPE_STRING:
			fallthrough
		case tokeniser.TOKEN_TYPE_NUMBER:
			fallthrough
		case tokeniser.TOKEN_TYPE_BOOL:
			fallthrough
		case tokeniser.TOKEN_TYPE_OPEN_LIST:
			fallthrough
		case tokeniser.TOKEN_TYPE_CONSTANT:
			{
				value, err := p.ParseValue()
				if err != nil {
					return nil, err
				}

				comma_or_close := p.Peek()

				if comma_or_close.Type == tokeniser.TOKEN_TYPE_COMMA {
					p.Increment()
				} else if comma_or_close.Type != tokeniser.TOKEN_TYPE_CLOSE_LIST {
					return nil, p.FormatErrorAtToken("Expected comma or closing bracket", comma_or_close.Start)
				}

				list = append(list, value)
			}
		default:
			{
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
		case tokeniser.TOKEN_TYPE_KEY:
			fallthrough
		case tokeniser.TOKEN_TYPE_STRING:
			{
				key := token.Value

				assign := p.Consume()

				if assign.Type != tokeniser.TOKEN_TYPE_ASSIGN {
					return nil, p.FormatErrorAtToken("Expected assignment operator `=`", assign.Start)
				}

				value, err := p.ParseValue()
				if err != nil {
					return nil, err
				}

				object[key] = value

				optional_comma := p.Peek()

				if optional_comma.Type == tokeniser.TOKEN_TYPE_COMMA {
					p.Increment()
				}
			}
		default:
			{
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
		case tokeniser.TOKEN_TYPE_KEY:
			fallthrough
		case tokeniser.TOKEN_TYPE_STRING:
			{
				key := token.Value

				assign := p.Consume()

				if assign.Type != tokeniser.TOKEN_TYPE_ASSIGN {
					return nil, p.FormatErrorAtToken("Expected assignment operator `=`", assign.Start)
				}

				value, err := p.ParseValue()
				if err != nil {
					return nil, err
				}

				globalObject[key] = value
			}
		case tokeniser.TOKEN_TYPE_CONSTANT:
			{
				key := token.Value

				assign := p.Consume()

				if assign.Type != tokeniser.TOKEN_TYPE_ASSIGN {
					return nil, p.FormatErrorAtToken("Expected assignment operator `=`", assign.Start)
				}

				value, err := p.ParseValue()
				if err != nil {
					return nil, err
				}

				p.constants[key] = value
			}
		case tokeniser.TOKEN_TYPE_OPEN_OBJ:
			{
				object, err := p.ParseObject()
				if err != nil {
					return nil, err
				}

				for k, v := range object {
					globalObject[k] = v
				}
			}
		case tokeniser.TOKEN_TYPE_KEYWORD:
			{
				switch token.Value {
				case "include":
					{
						includePath := p.Consume()

						if includePath.Type != tokeniser.TOKEN_TYPE_STRING {
							return nil, p.FormatErrorAtToken("Expected string path to include", includePath.Start)
						}

						f, err := os.ReadFile(filepath.Join(p.rootDir, includePath.Value))
						if err != nil {
							return nil, err
						}

						s := string(f)

						t := tokeniser.NewTokeniser(s, includePath.Value)
						tokens, err := t.Tokenise()
						if err != nil {
							return nil, err
						}

						p2 := NewParser(tokens, p.rootDir, includePath.Value)
						object, err := p2.Parse()
						if err != nil {
							return nil, err
						}

						for k, v := range object {
							globalObject[k] = v
						}
					}
				}
			}
		default:
			{
				return nil, p.FormatErrorAtToken(fmt.Sprintf("Unexpected token %s", token.Type), token.Start)
			}
		}
	}

	return globalObject, nil
}
