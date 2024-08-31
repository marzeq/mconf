package parser

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
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
	GetFloat() (*big.Float, error)
	GetInt() (*big.Int, error)
	GetBool() (bool, error)
	GetList() ([]ParserValue, error)
	GetObject() (map[string]ParserValue, error)
}

type importCacheEntry struct {
	values    map[string]ParserValue
	constants map[string]ParserValue
}

type Parser struct {
	tokens      []tokeniser.Token
	currIndex   int
	rootDir     string
	currentFile string
	importCache *map[string]importCacheEntry
}

func NewParser(tokens []tokeniser.Token, rootDir string, currentFile string) Parser {
	importCache := make(map[string]importCacheEntry)

	fullFile := filepath.Join(rootDir, currentFile)

	importCache[fullFile] = importCacheEntry{
		values:    make(map[string]ParserValue),
		constants: make(map[string]ParserValue),
	}

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		importCache[fullFile].constants[pair[0]] = &ParserValueString{Value: pair[1]}
	}

	return Parser{
		tokens:      tokens,
		currIndex:   0,
		rootDir:     rootDir,
		currentFile: currentFile,
		importCache: &importCache,
	}
}

func (p *Parser) childParser(tokens []tokeniser.Token, currentFile string) Parser {
	fullFile := filepath.Join(p.rootDir, currentFile)

	(*p.importCache)[fullFile] = importCacheEntry{
		values:    make(map[string]ParserValue),
		constants: make(map[string]ParserValue),
	}

	return Parser{
		tokens:      tokens,
		currIndex:   0,
		rootDir:     p.rootDir,
		currentFile: currentFile,
		importCache: p.importCache,
	}
}

func (p *Parser) GetValues() map[string]ParserValue {
	return (*p.importCache)[filepath.Join(p.rootDir, p.currentFile)].values
}

func (p *Parser) GetConstants() map[string]ParserValue {
	return (*p.importCache)[filepath.Join(p.rootDir, p.currentFile)].constants
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
		return fmt.Errorf(fmt.Sprintf("%s - Parser error at EOF: %s", prettyFile, message))
	}

	return fmt.Errorf(fmt.Sprintf("%s:%d:%d - Parser error: %s", prettyFile, loc.Line, loc.Col, message))
}

func (p *Parser) ParseDeepKey() ([]string, error) {
	key := make([]string, 0)

	for {
		token := p.Consume()

		if token.Type == tokeniser.TOKEN_TYPE_KEY || token.Type == tokeniser.TOKEN_TYPE_STRING {
			key = append(key, token.Value)
		}

		next := p.Peek()

		if next.Type == tokeniser.TOKEN_TYPE_DOT {
			p.Increment()
		} else {
			break
		}
	}

	return key, nil
}

func (p *Parser) ParseValue() (ParserValue, error) {
	token := p.Consume()

	switch token.Type {
	case tokeniser.TOKEN_TYPE_STRING:
		sb := ""
		for i, v := range token.Values {
			sb += v

			if i < len(token.StringSubs) {
				constantName := token.StringSubs[i]
				constantValue, ok := p.GetConstants()[constantName]
				if !ok {
					return nil, p.FormatErrorAtToken(fmt.Sprintf("Constant in string substitution `%s` not found", constantName), token.Start)
				}

				if constantValue.GetType() != PARSER_VALUE_TYPE_STRING {
					sb += constantValue.ValueToString()
				} else {
					constantStr, err := constantValue.GetString()
					if err != nil {
						return nil, err
					}

					sb += constantStr
				}
			}
		}

		return &ParserValueString{Value: sb}, nil
	case tokeniser.TOKEN_TYPE_NUMBER_DECIMAL:
		if strings.Contains(token.Value, ".") {
			bigFl, _, err := big.ParseFloat(token.Value, 10, 0, big.ToNearestEven)
			if err != nil {
				return nil, p.FormatErrorAtToken(fmt.Sprintf("Failed to convert `%s` to float", token.Value), token.Start)
			}

			return &ParserValueFloat{Value: bigFl}, nil
		} else {
			intVal, success := new(big.Int).SetString(token.Value, 10)
			if !success {
				return nil, p.FormatErrorAtToken(fmt.Sprintf("Failed to convert `%s` to decimal int", token.Value), token.Start)
			}

			return &ParserValueInt{Value: intVal}, nil
		}
	case tokeniser.TOKEN_TYPE_NUMBER_HEX:
		intVal, success := new(big.Int).SetString(token.Value, 16)
		if !success {
			return nil, p.FormatErrorAtToken(fmt.Sprintf("Failed to convert `%s` to hex int", token.Value), token.Start)
		}

		return &ParserValueInt{Value: intVal}, nil
	case tokeniser.TOKEN_TYPE_NUMBER_BINARY:
		intVal, success := new(big.Int).SetString(token.Value, 2)
		if !success {
			return nil, p.FormatErrorAtToken(fmt.Sprintf("Failed to convert `%s` to binary int", token.Value), token.Start)
		}

		return &ParserValueInt{Value: intVal}, nil
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
		value, ok := p.GetConstants()[token.Value]

		peeked := p.Peek()

		backup := ""

		if peeked.Type == tokeniser.TOKEN_TYPE_QUESTION_MARK {
			p.Increment()

			b := p.Consume()

			if b.Type != tokeniser.TOKEN_TYPE_CONSTANT {
				return nil, p.FormatErrorAtToken("Expected constant after `?`", b.Start)
			}

			backup = b.Value
		}

		if !ok {
			if backup != "" {
				value, ok = p.GetConstants()[backup]
			}
		}

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
		case tokeniser.TOKEN_TYPE_NUMBER_DECIMAL:
			fallthrough
		case tokeniser.TOKEN_TYPE_NUMBER_HEX:
			fallthrough
		case tokeniser.TOKEN_TYPE_NUMBER_BINARY:
			fallthrough
		case tokeniser.TOKEN_TYPE_BOOL:
			fallthrough
		case tokeniser.TOKEN_TYPE_OPEN_LIST:
			fallthrough
		case tokeniser.TOKEN_TYPE_OPEN_OBJ:
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

func (p *Parser) SmartlySetValuesAndConstants(importEverything bool, importPaths [][]string, importConstants []string, ic importCacheEntry, errorLoc tokeniser.Location, importPath string) error {
	if importEverything {
		for k, v := range ic.values {
			p.GetValues()[k] = v
		}

		for k, v := range ic.constants {
			p.GetConstants()[k] = v
		}
	} else {
		for _, path := range importPaths {
			current := ic.values

			for i, key := range path {
				indexedVal, ok := current[key]
				if !ok {
					joinedPath := strings.Join(path[:i+1], ".")
					return p.FormatErrorAtToken(fmt.Sprintf("Path `%s` not found in imported file %s", joinedPath, importPath), errorLoc)
				}

				if i == len(path)-1 {
					p.GetValues()[key] = indexedVal
					break
				}

				got, err := indexedVal.GetObject()
				if err != nil {
					return p.FormatErrorAtToken(fmt.Sprintf("Path `%s` in imported file %s is not an object", strings.Join(path[:i+1], "."), importPath), errorLoc)
				}
				current = got
			}
		}

		for k, v := range ic.constants {
			for _, constant := range importConstants {
				if k == constant {
					p.GetConstants()[k] = v
				}
			}
		}
	}

	return nil
}

func (p *Parser) Parse() (map[string]ParserValue, error) {
	for {
		token := p.Consume()

		switch token.Type {
		case tokeniser.TOKEN_TYPE_EOF:
			return p.GetValues(), nil
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

				p.GetValues()[key] = value
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

				p.GetConstants()[key] = value
			}
		case tokeniser.TOKEN_TYPE_OPEN_OBJ:
			{
				object, err := p.ParseObject()
				if err != nil {
					return nil, err
				}

				for k, v := range object {
					p.GetValues()[k] = v
				}
			}
		case tokeniser.TOKEN_TYPE_DIRECTIVE:
			{
				switch token.Value {
				case "import":
					{
						nextUnknown := p.Peek()

						importPaths := [][]string{}
						importConstants := []string{}
						importEverything := true

						if nextUnknown.Type == tokeniser.TOKEN_TYPE_OPEN_OBJ {
							p.Increment()
							importEverything = false
							for {
								tok := p.Peek()

								if tok.Type == tokeniser.TOKEN_TYPE_CLOSE_OBJ {
									p.Increment()
									break
								}

								if tok.Type == tokeniser.TOKEN_TYPE_CONSTANT {
									p.Increment()
									importConstants = append(importConstants, tok.Value)
								} else {
									key, err := p.ParseDeepKey()
									if err != nil {
										return nil, err
									}

									importPaths = append(importPaths, key)
								}

								comma_or_close := p.Peek()

								if comma_or_close.Type == tokeniser.TOKEN_TYPE_COMMA {
									p.Increment()
								} else if comma_or_close.Type != tokeniser.TOKEN_TYPE_CLOSE_OBJ {
									return nil, p.FormatErrorAtToken("Expected comma or closing bracket", comma_or_close.Start)
								}
							}
						}

						importPath := p.Consume()

						if importPath.Type != tokeniser.TOKEN_TYPE_STRING {
							return nil, p.FormatErrorAtToken("Expected string path to import", importPath.Start)
						}

						if importPath.Value == p.currentFile {
							return nil, p.FormatErrorAtToken("Cannot import the same file", importPath.Start)
						}

						fullFilePath := filepath.Join(p.rootDir, importPath.Value)
						relative, err := filepath.Rel(p.rootDir, fullFilePath)

						ic, icOk := (*p.importCache)[fullFilePath]

						if icOk {
							err := p.SmartlySetValuesAndConstants(importEverything, importPaths, importConstants, ic, importPath.Start, importPath.Value)
							if err != nil {
								return nil, err
							}
							continue
						}

						f, err := os.ReadFile(fullFilePath)
						if err != nil {
							return nil, err
						}

						s := string(f)

						t := tokeniser.NewTokeniser(s, relative)
						tokens, errTokenise := t.Tokenise()
						if errTokenise != nil {
							return nil, errTokenise
						}

						p2 := p.childParser(tokens, relative)
						_, errParse := p2.Parse()
						if errParse != nil {
							return nil, errParse
						}

						ic, icOk = (*p.importCache)[fullFilePath]

						if !icOk {
							return nil, fmt.Errorf("Unreachable code reached, please report this as a bug")
						}

						err = p.SmartlySetValuesAndConstants(importEverything, importPaths, importConstants, ic, importPath.Start, importPath.Value)
						if err != nil {
							return nil, err
						}
					}
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

	return p.GetValues(), nil
}
