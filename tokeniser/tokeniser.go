package tokeniser

import (
	"fmt"
	"os"
	"path"
	"unicode"
)

type Tokeniser struct {
	contents    []rune
	currIndex   int
	filePath    string
	relativeDir string
}

func NewTokeniser(contents string, filePath string, relativeDir string) Tokeniser {
	return Tokeniser{
		contents:    []rune(contents),
		currIndex:   0,
		filePath:    filePath,
		relativeDir: relativeDir,
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

func (t *Tokeniser) ReadUnicodeEscape(backslashULoc Location) (rune, error) {
	hc1 := t.Consume()
	if !IsHexDigit(hc1) {
		return 0, t.FormatErrorAt(fmt.Sprintf("Expected hex digit after `\\u`, got `%c`", hc1), backslashULoc)
	}
	h1 := hexToDec(hc1)

	hc2 := t.Consume()
	if !IsHexDigit(hc2) {
		return 0, t.FormatErrorAt(fmt.Sprintf("Expected hex digit after `\\u%c`, got `%c`", hc1, hc2), backslashULoc)
	}
	h2 := hexToDec(hc2)

	hc3 := t.Consume()
	if !IsHexDigit(hc3) {
		return 0, t.FormatErrorAt(fmt.Sprintf("Expected hex digit after `\\u%c%c`, got `%c`", hc1, hc2, hc3), backslashULoc)
	}
	h3 := hexToDec(hc3)

	hc4 := t.Consume()
	if !IsHexDigit(hc4) {
		return 0, t.FormatErrorAt(fmt.Sprintf("Expected hex digit after `\\u%c%c%c`, got `%c`", hc1, hc2, hc3, hc4), backslashULoc)
	}
	h4 := hexToDec(hc4)

	// may be a standalone code point or part of a surrogate pair
	ch := h1*16*16*16 + h2*16*16 + h3*16 + h4

	// check if a high surrogate [0xD800..0xDBFF]
	if ch >= 0xD800 && ch <= 0xDBFF {
		if t.Peek() == '\\' && t.PeekAhead(1) == 'u' {
			t.Consume()
			t.Consume()
			lowSurrogate, err := t.ReadUnicodeEscape(backslashULoc)
			if err != nil {
				return 0, err
			}

			// ensure a valid low surrogate [0xDC00..0xDFFF]
			if lowSurrogate < 0xDC00 || lowSurrogate > 0xDFFF {
				// return null character if not a valid low surrogate rather than an error (it's their fault)
				return 0, nil
			}

			// combine the surrogate pair into a single code point using documented formula
			fullCodePoint := (ch-0xD800)*0x400 + (int(lowSurrogate) - 0xDC00) + 0x10000
			return rune(fullCodePoint), nil
		} else {
			return rune(ch), nil // it's their fault they didn't provide a low surrogate
		}
	}

	return rune(ch), nil
}

func hexToDec(hc rune) int {
	if hc >= '0' && hc <= '9' {
		return int(hc - '0')
	}
	if hc >= 'A' && hc <= 'F' {
		return int(hc - 'A' + 10)
	}
	if hc >= 'a' && hc <= 'f' {
		return int(hc - 'a' + 10)
	}
	return 0
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
				fallthrough
			case 'X':
				hc1 := t.Consume()
				if !IsHexDigit(hc1) {
					return nil, nil, t.FormatErrorAt(fmt.Sprintf("Expected hex digit after `\\x`, got `%c`", hc1), nextloc)
				}
				h1 := hexToDec(hc1)
				hc2 := t.Consume()
				if !IsHexDigit(hc2) {
					return nil, nil, t.FormatErrorAt(fmt.Sprintf("Expected hex digit after `\\x%c`, got `%c`", hc1, hc2), nextloc)
				}
				h2 := hexToDec(hc2)
				strings[len(strings)-1] += string(rune(h1*16 + h2))
			case 'u':
				fallthrough
			case 'U':
				unicodeChar, err := t.ReadUnicodeEscape(nextloc)
				if err != nil {
					return nil, nil, err
				}
				strings[len(strings)-1] += string(unicodeChar)
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
			} else if next == 'e' || next == 'E' {
				number += string(next)
				t.Increment()

				next = t.Peek()
				if next == '+' || next == '-' {
					number += string(next)
					t.Increment()
				}

				next = t.Peek()
				if !unicode.IsDigit(next) {
					return "", "", t.FormatError(fmt.Sprintf("Expected digit after exponent in decimal number, got `%c`", next))
				}

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
	var prettyFile string

	if t.filePath == "" {
		prettyFile = "(stdin)"
	} else {
		prettyFile = path.Join(t.relativeDir, t.filePath)
	}

	return fmt.Errorf(fmt.Sprintf("%s:%d:%d - Tokeniser error: %s", prettyFile, loc.Line, loc.Col, message))
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

			if word == "true" || word == "yes" || word == "on" {
				tokens = append(tokens, BoolToken("true", loc))
			} else if word == "false" || word == "no" || word == "off" {
				tokens = append(tokens, BoolToken("false", loc))
			} else if word == "null" {
				tokens = append(tokens, NullToken(loc))
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

			if c == '=' || c == ':' {
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
			} else if c == '~' {
				tokens = append(tokens, TildeToken(loc))
			} else if c == '|' {
				tokens = append(tokens, PipeToken(loc))
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
