package tokeniser

import (
	"fmt"
	"os"
	"unicode"
)

const (
	TOKEN_TYPE_KEY        = "KEY"
	TOKEN_TYPE_CONSTANT   = "CONSTANT"
	TOKEN_TYPE_ASSIGN     = "ASSIGN"
	TOKEN_TYPE_NUMBER     = "NUMBER"
	TOKEN_TYPE_STRING     = "STRING"
	TOKEN_TYPE_BOOL       = "BOOL"
	TOKEN_TYPE_OPEN_LIST  = "OPEN_LIST"
	TOKEN_TYPE_CLOSE_LIST = "CLOSE_LIST"
	TOKEN_TYPE_COMMA      = "COMMA"
	TOKEN_TYPE_OPEN_OBJ   = "OPEN_OBJ"
	TOKEN_TYPE_CLOSE_OBJ  = "CLOSE_OBJ"
)

const (
	NO_VALUE = "NO_VALUE"
)

type Location struct {
	Line int
	Col  int
}

func (l Location) Repr() string {
	return fmt.Sprintf("Location{Line: %d, Col: %d}", l.Line, l.Col)
}

type Token struct {
	Type string
	Value     string
	Start     Location
}

func (t Token) Repr() string {
	if t.Value == NO_VALUE {
		return fmt.Sprintf(`Token{
  Type: %s,
  Location: %s
}`, t.Type, t.Start.Repr())
	}
	return fmt.Sprintf(`Token{
  Type: %s,
  Value: %s,
  Location: %s
}`, t.Type, t.Value, t.Start.Repr())
}

func KeyToken(value string, start Location) Token {
	return Token{
		Type: TOKEN_TYPE_KEY,
		Value:     value,
		Start:     start,
	}
}

func ConstantToken(value string, start Location) Token {
	return Token{
		Type: TOKEN_TYPE_CONSTANT,
		Value:     value,
		Start:     start,
	}
}

func AssignToken(start Location) Token {
	return Token{
		Type: TOKEN_TYPE_ASSIGN,
		Value:     NO_VALUE,
		Start:     start,
	}
}

func NumberToken(value string, start Location) Token {
	return Token{
		Type: TOKEN_TYPE_NUMBER,
		Value:     value,
		Start:     start,
	}
}

func StringToken(value string, start Location) Token {
	return Token{
		Type: TOKEN_TYPE_STRING,
		Value:     value,
		Start:     start,
	}
}

func BoolToken(value string, start Location) Token {
	return Token{
		Type: TOKEN_TYPE_BOOL,
		Value:     value,
		Start:     start,
	}
}

func OpenListToken(start Location) Token {
	return Token{
		Type: TOKEN_TYPE_OPEN_LIST,
		Value:     NO_VALUE,
		Start:     start,
	}
}

func CloseListToken(start Location) Token {
	return Token{
		Type: TOKEN_TYPE_CLOSE_LIST,
		Value:     NO_VALUE,
		Start:     start,
	}
}

func CommaToken(start Location) Token {
	return Token{
		Type: TOKEN_TYPE_COMMA,
		Value:     NO_VALUE,
		Start:     start,
	}
}

func OpenObjToken(start Location) Token {
	return Token{
		Type: TOKEN_TYPE_OPEN_OBJ,
		Value:     NO_VALUE,
		Start:     start,
	}
}

func CloseObjToken(start Location) Token {
	return Token{
		Type: TOKEN_TYPE_CLOSE_OBJ,
		Value:     NO_VALUE,
		Start:     start,
	}
}

type Tokeniser struct {
	contents  []rune
	currIndex int
}

func NewTokeniser(contents string) Tokeniser {
	return Tokeniser{
		contents:  []rune(contents),
		currIndex: 0,
	}
}

func (t *Tokeniser) PeekAhead(i int) rune {
	if t.currIndex+i >= len(t.contents) {
		return 0
	}

	return t.contents[t.currIndex+i]
}

func (t *Tokeniser) Peek() rune {
	return t.PeekAhead(0)
}

func (t *Tokeniser) Increment() {
	t.currIndex++
}

func (t *Tokeniser) Consume() rune {
	c := t.Peek()

	t.Increment()

	return c
}

func (t *Tokeniser) GoBack() {
	t.currIndex--
}

func (t *Tokeniser) ReadString() string {
	s := ""

	loc := t.GetCurrLineAndCol()
	initial := t.Consume()

	if initial != '"' && initial != '\'' {
		fmt.Println(t.FormatErrorAt("Expected `\"` or `'` to start string", loc))
		os.Exit(1)
	}

	for {
		c := t.Consume()

		if c == '\\' {
			nextloc := t.GetCurrLineAndCol()
			next := t.Consume()

			switch next {
			case '"':
				s += "\""
			case '\'':
				s += "'"
			case '\\':
				s += "\\"
			case 'n':
				s += "\n"
			case 'r':
				s += "\r"
			case 't':
				s += "\t"
			case 'f':
				s += "\f"
			default:
				fmt.Println(t.FormatErrorAt(fmt.Sprintf("Unknown escape sequence: `\\%c`", next), nextloc))
				os.Exit(1)
			}
		} else if c == initial {
			break
		} else {
			s += string(c)
		}
	}

	return s
}

func (t *Tokeniser) ReadWord() string {
	loc := t.GetCurrLineAndCol()
	initial := t.Consume()

	if !unicode.IsLetter(initial) {
		fmt.Println(t.FormatErrorAt("Expected letter to start a word", loc))
	}

	word := string(initial)

	for {
		next := t.Peek()

		if unicode.IsLetter(next) || (unicode.IsNumber(next) && len(word) > 0) {
			word += string(next)
			t.Increment()
		} else {
			break
		}
	}

	return word
}

func (t *Tokeniser) ReadNumber() string {
	loc := t.GetCurrLineAndCol()
	initial := t.Consume()

	if !unicode.IsDigit(initial) {
		fmt.Println(t.FormatErrorAt("Expected digit to start a number", loc))
	}

	number := string(initial)

	for {
		next := t.Peek()

		if unicode.IsDigit(next) || next == '.' {
			number += string(next)
			t.Increment()
		} else if next == '_' {
			t.Increment()
		} else if unicode.IsLetter(next) {
			fmt.Println(t.FormatErrorAt(fmt.Sprintf("Unexpected character in number: `%c`", next), t.GetCurrLineAndCol()))
			os.Exit(1)
		} else {
			break
		}
	}

	return number
}

func (t *Tokeniser) IgnoreComment() {
	loc := t.GetCurrLineAndCol()
	initial := t.Consume()

	if initial != '#' {
		fmt.Println(t.FormatErrorAt("Expected `#` to start a comment, this is a bug, please report it", loc))
		os.Exit(1)
	}

	for {
		next := t.Peek()

		if next == '\n' || next == 0 {
			break
		}

		t.Increment()
	}
}

func (t *Tokeniser) GetLineAndCol(index int) Location {
	line := 1
	col := 1

	for i := 0; i < index; i++ {
		if t.contents[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}

	return Location{line, col}
}

func (t *Tokeniser) GetCurrLineAndCol() Location {
	return t.GetLineAndCol(t.currIndex)
}

func (t *Tokeniser) FormatErrorAt(message string, loc Location) string {
	return fmt.Sprintf("Tokeniser error at line %d, col %d: %s", loc.Line, loc.Col, message)
}

func (t *Tokeniser) FormatError(message string) string {
	loc := t.GetCurrLineAndCol()

	return fmt.Sprintf("Tokeniser error at line %d, col %d: %s", loc.Line, loc.Col, message)
}

func (t *Tokeniser) Tokenise() []Token {
	tokens := []Token{}

	if len(t.contents) == 0 {
		return tokens
	}

	for {
		loc := t.GetCurrLineAndCol()
		c := t.Peek()

		if c == 0 {
			break
		}

		if unicode.IsLetter(c) {
			word := t.ReadWord()

			if word == "true" || word == "false" {
				tokens = append(tokens, BoolToken(word, loc))
			} else {
				tokens = append(tokens, KeyToken(word, loc))
			}
		} else if unicode.IsDigit(c) {
			number := t.ReadNumber()

			tokens = append(tokens, NumberToken(number, loc))
		} else if c == '"' || c == '\'' {
			parsed := t.ReadString()

			tokens = append(tokens, StringToken(parsed, loc))
		} else if c == '#' {
			t.IgnoreComment()
		} else {
			t.Increment()

			if c == '=' {
				tokens = append(tokens, AssignToken(loc))
			} else if c == '$' {
				word := t.ReadWord()

				tokens = append(tokens, ConstantToken(word, loc))
			} else if c == '[' {
				tokens = append(tokens, OpenListToken(loc))
			} else if c == ']' {
				tokens = append(tokens, CloseListToken(loc))
			} else if c == ',' {
				tokens = append(tokens, CommaToken(loc))
			} else if c == '{' {
				tokens = append(tokens, OpenObjToken(loc))
			} else if c == '}' {
				tokens = append(tokens, CloseObjToken(loc))
			} else if unicode.IsSpace(c) {
				continue
			} else {
				fmt.Println(t.FormatErrorAt(fmt.Sprintf("Unexpected character: `%c`", c), loc))
				os.Exit(1)
			}
		}
	}

	return tokens
}
