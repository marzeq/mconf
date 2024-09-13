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
		filePath:  filePath,
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

func IsHexDigit(c rune) bool {
	return IsAsciiDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func (t *Tokeniser) ReadString() ([]string, []string, error) {
	strings := []string{""}

	loc := t.GetCurrLineAndCol()
	initial := t.Consume()

	constantSubs := []string{}

	if initial != '"' {
		return nil, nil, t.FormatErrorAt("Expected `\"` to start string", loc)
	}

	for {
		c := t.Consume()

		if c == 0 {
			return nil, nil, t.FormatErrorAt("Unexpected end of file in string", loc)
		}

		if c == '$' {
			openbrack := t.Consume()

			if openbrack != '{' {
				return nil, nil, t.FormatErrorAt("Expected `{` after `$` in formatted string", loc)
			}

			constantName, error := t.ReadWord()

			if error != nil {
				return nil, nil, error
			}

			closebrack := t.Consume()

			if closebrack != '}' {
				return nil, nil, t.FormatErrorAt("Expected `}` after constant name in formatted string", loc)
			}

			constantSubs = append(constantSubs, constantName)
			strings = append(strings, "")
			continue
		}

		if c == '\\' {
			nextloc := t.GetCurrLineAndCol()
			next := t.Consume()

			switch next {
			case '"':
				strings[len(strings)-1] += "\""
			case '\'':
				strings[len(strings)-1] += "'"
			case '\\':
				strings[len(strings)-1] += "\\"
			// source https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797#general-ascii-codes
			case 'a':
				strings[len(strings)-1] += "\a"
			case 'b':
				strings[len(strings)-1] += "\b"
			case 't':
				strings[len(strings)-1] += "\t"
			case 'n':
				strings[len(strings)-1] += "\n"
			case 'v':
				strings[len(strings)-1] += "\v"
			case 'f':
				strings[len(strings)-1] += "\f"
			case 'r':
				strings[len(strings)-1] += "\r"
			case 'e':
				strings[len(strings)-1] += "\x1b"
			// end source
			case 'x':
				// read two hex digits (one byte)
				hc1 := t.Consume()
				if !IsHexDigit(hc1) {
					return nil, nil, t.FormatErrorAt(fmt.Sprintf("Expected hex digit after `\\x`, got `%c`", hc1), nextloc)
				}
				h1 := hc1 - '0'
				hc2 := t.Consume()
				if !IsHexDigit(hc2) {
					return nil, nil, t.FormatErrorAt(fmt.Sprintf("Expected hex digit after `\\x%c`, got `%c`", hc1, hc2), nextloc)
				}
				h2 := hc2 - '0'
				ch := h1*16 + h2
				strings[len(strings)-1] += string(ch)
			case 'u':
				fallthrough
			case 'U':
				// read four hex digits (two bytes aka one unicode code point)
				hc1 := t.Consume()
				if !IsHexDigit(hc1) {
					return nil, nil, t.FormatErrorAt(fmt.Sprintf("Expected hex digit after `\\u`, got `%c`", hc1), nextloc)
				}
				h1 := hc1 - '0'
				hc2 := t.Consume()
				if !IsHexDigit(hc2) {
					return nil, nil, t.FormatErrorAt(fmt.Sprintf("Expected hex digit after `\\u%c`, got `%c`", hc1, hc2), nextloc)
				}
				h2 := hc2 - '0'
				hc3 := t.Consume()
				if !IsHexDigit(hc3) {
					return nil, nil, t.FormatErrorAt(fmt.Sprintf("Expected hex digit after `\\u%c%c`, got `%c`", hc1, hc2, hc3), nextloc)
				}
				h3 := hc3 - '0'
				hc4 := t.Consume()
				if !IsHexDigit(hc4) {
					return nil, nil, t.FormatErrorAt(fmt.Sprintf("Expected hex digit after `\\u%c%c%c`, got `%c`", hc1, hc2, hc3, hc4), nextloc)
				}
				h4 := hc4 - '0'
				ch := h1*16*16*16 + h2*16*16 + h3*16 + h4
				strings[len(strings)-1] += string(ch)
			case '$':
				strings[len(strings)-1] += "$"
			default:
				return nil, nil, t.FormatErrorAt(fmt.Sprintf("Unknown escape sequence: `\\%c`", next), nextloc)
			}
		} else if c == initial {
			break
		} else {
			strings[len(strings)-1] += string(c)
		}
	}

	return strings, constantSubs, nil
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

func (t *Tokeniser) ReadNumber() (string, string, error) {
	loc := t.GetCurrLineAndCol()
	initial := t.Consume()

	if !IsAsciiDigit(initial) && initial != '-' && initial != '.' {
		return "", "", t.FormatErrorAt("Expected digit to start a number", loc)
	}

	number := string(initial)

	if initial == '.' {
		number = "0."
	}

	var mode string

	if initial == '-' {
		number = "-"
		mode = TOKEN_TYPE_NUMBER_DECIMAL

		if t.Peek() == '0' && t.PeekAhead(1) == 'x' {
			mode = TOKEN_TYPE_NUMBER_HEX
			t.Increment()
			t.Increment()
		} else if t.Peek() == '0' && t.PeekAhead(1) == 'b' {
			mode = TOKEN_TYPE_NUMBER_BINARY
			t.Increment()
			t.Increment()
		}
	} else if initial == '0' && t.Peek() == 'x' {
		mode = TOKEN_TYPE_NUMBER_HEX
		t.Increment()
	} else if initial == '0' && t.Peek() == 'b' {
		mode = TOKEN_TYPE_NUMBER_BINARY
		t.Increment()
	} else {
		mode = TOKEN_TYPE_NUMBER_DECIMAL
	}

	for {
		next := t.Peek()
		if mode == TOKEN_TYPE_NUMBER_DECIMAL {
			if unicode.IsDigit(next) || next == '.' {
				number += string(next)
				t.Increment()
			} else if next == '_' {
				t.Increment()
			} else if unicode.IsLetter(next) {
				return "", "", t.FormatError(fmt.Sprintf("Unexpected character in decimal number: `%c`", next))
			} else {
				break
			}
		} else if mode == TOKEN_TYPE_NUMBER_HEX {
			if unicode.IsDigit(next) || (next >= 'a' && next <= 'f') || (next >= 'A' && next <= 'F') {
				number += string(next)
				t.Increment()
			} else if next == '_' {
				t.Increment()
			} else if unicode.IsLetter(next) {
				return "", "", t.FormatError(fmt.Sprintf("Unexpected character in hex number: `%c`", next))
			} else {
				break
			}
		} else if mode == TOKEN_TYPE_NUMBER_BINARY {
			if next == '0' || next == '1' {
				number += string(next)
				t.Increment()
			} else if next == '_' {
				t.Increment()
			} else if unicode.IsLetter(next) {
				return "", "", t.FormatError(fmt.Sprintf("Unexpected character in binary number: `%c`", next))
			} else {
				break
			}
		} else {
			return "", "", t.FormatError(fmt.Sprintf("Unreachable code reached, this is a bug, please report it"))
		}
	}

	return number, mode, nil
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
	return fmt.Errorf(fmt.Sprintf("%s:%d:%d - Tokeniser error: %s", t.filePath, loc.Line, loc.Col, message))
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

				for {
					next := t.Peek()
					if unicode.IsSpace(next) {
						t.Increment()
					} else {
						break
					}
				}

				if t.Peek() == '.' {
					t.Increment()
					tokens = append(tokens, DotToken(loc))
				}
			}
		} else if IsAsciiDigit(c) || c == '-' || c == '.' {
			number, mode, error := t.ReadNumber()

			if error != nil {
				return nil, error
			}

			tokens = append(tokens, NumberToken(number, mode, loc))
		} else if c == '"' {
			parsed, constantSubs, error := t.ReadString()

			if error != nil {
				return nil, error
			}

			tokens = append(tokens, StringToken(parsed, constantSubs, loc))
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

				tokens = append(tokens, DirectiveToken(word, loc))
			} else if c == '[' {
				tokens = append(tokens, OpenListToken(loc))
			} else if c == ']' {
				tokens = append(tokens, CloseListToken(loc))
			} else if c == '?' {
				tokens = append(tokens, QuestionMarkToken(loc))
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
