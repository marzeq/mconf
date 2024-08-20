package tokeniser

import (
	"fmt"
	"os"
	"unicode"
)

type Tokeniser struct {
	contents  []rune
	currIndex int
	filePath  string
}

func NewTokeniser(contents string, filePath string) Tokeniser {
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

func (t *Tokeniser) ReadString() (string, error) {
	s := ""

	loc := t.GetCurrLineAndCol()
	initial := t.Consume()

	if initial != '"' {
		return "", t.FormatErrorAt("Expected `\"` to start string", loc)
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
				return "", t.FormatErrorAt(fmt.Sprintf("Unknown escape sequence: `\\%c`", next), nextloc)
			}
		} else if c == initial {
			break
		} else {
			s += string(c)
		}
	}

	return s, nil
}

func (t *Tokeniser) ReadWord() (string, error) {
	loc := t.GetCurrLineAndCol()
	initial := t.Consume()

	if !IsLegalWordStart(initial) {
		return "", t.FormatErrorAt("Expected letter to start a word", loc)
	}

	word := string(initial)

	for {
		next := t.Peek()

		if IsLegalWordStart(next) || (IsAsciiDigit(next) && len(word) > 0) {
			word += string(next)
			t.Increment()
		} else {
			break
		}
	}

	return word, nil
}

func (t *Tokeniser) ReadNumber() (string, error) {
	loc := t.GetCurrLineAndCol()
	initial := t.Consume()

	if !IsAsciiDigit(initial) && initial != '-' && initial != '.' {
		return "", t.FormatErrorAt("Expected digit to start a number", loc)
	}

	number := string(initial)

	if initial == '.' {
		number = "0."
	}

	for {
		next := t.Peek()

		if unicode.IsDigit(next) || next == '.' {
			number += string(next)
			t.Increment()
		} else if next == '_' {
			t.Increment()
		} else if unicode.IsLetter(next) {
			return "", t.FormatError(fmt.Sprintf("Unexpected character in number: `%c`", next))
		} else {
			break
		}
	}

	return number, nil
}

func (t *Tokeniser) IgnoreComment() error {
	loc := t.GetCurrLineAndCol()
	initial := t.Consume()

	if initial != '#' {
		return t.FormatErrorAt("Expected `#` to start a comment, this is a bug, please report it", loc)
	}

	for {
		next := t.Peek()

		if next == '\n' || next == 0 {
			break
		}

		t.Increment()
	}

	return nil
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

func (t *Tokeniser) FormatErrorAt(message string, loc Location) error {
	return fmt.Errorf(fmt.Sprintf("%s - Tokeniser error at line %d, col %d: %s", t.filePath, loc.Line, loc.Col, message))
}

func (t *Tokeniser) FormatError(message string) error {
	return t.FormatErrorAt(message, t.GetCurrLineAndCol())
}

func IsAsciiDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func IsLatin(c rune) bool {
	return unicode.Is(unicode.Latin, c)
}

func IsLegalWordStart(c rune) bool {
	return IsLatin(c) || c == '_'
}

func IsLegalWord(cs []rune) bool {
	for i, c := range cs {
		if i == 0 && !IsLegalWordStart(c) {
			return false
		}

		if !IsLegalWordStart(c) && !IsAsciiDigit(c) {
			return false
		}
	}

	return true
}

func (t *Tokeniser) Tokenise() ([]Token, error) {
	tokens := []Token{}

	if len(t.contents) == 0 {
		return tokens, nil
	}

	for {
		loc := t.GetCurrLineAndCol()
		c := t.Peek()

		if c == 0 {
			break
		}

		if IsLegalWordStart(c) {
			word, error := t.ReadWord()

			if error != nil {
				return nil, error
			}

			if word == "true" || word == "false" {
				tokens = append(tokens, BoolToken(word, loc))
			} else {
				tokens = append(tokens, KeyToken(word, loc))
			}
		} else if IsAsciiDigit(c) || c == '-' || c == '.' {
			number, error := t.ReadNumber()

			if error != nil {
				return nil, error
			}

			tokens = append(tokens, NumberToken(number, loc))
		} else if c == '"' {
			parsed, error := t.ReadString()

			if error != nil {
				return nil, error
			}

			tokens = append(tokens, StringToken(parsed, loc))
		} else if c == '#' {
			t.IgnoreComment()
		} else {
			t.Increment()

			if c == '=' {
				tokens = append(tokens, AssignToken(loc))
			} else if c == '$' {
				word, error := t.ReadWord()

				if error != nil {
					return nil, error
				}

				tokens = append(tokens, ConstantToken(word, loc))
			} else if c == '@' {
				word, error := t.ReadWord()

				if error != nil {
					return nil, error
				}

				tokens = append(tokens, KeywordToken(word, loc))
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

	return tokens, nil
}
