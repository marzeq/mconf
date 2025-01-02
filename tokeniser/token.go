package tokeniser

import (
	"fmt"
)

type TokenType int

const (
	TOKEN_TYPE_WORD TokenType = iota
	TOKEN_TYPE_CONSTANT
	TOKEN_TYPE_ASSIGN
	TOKEN_TYPE_NUMBER_DECIMAL
	TOKEN_TYPE_NUMBER_HEX
	TOKEN_TYPE_NUMBER_BINARY
	TOKEN_TYPE_STRING
	TOKEN_TYPE_BOOL
	TOKEN_TYPE_NULL
	TOKEN_TYPE_OPEN_LIST
	TOKEN_TYPE_CLOSE_LIST
	TOKEN_TYPE_COMMA
	TOKEN_TYPE_DOT
	TOKEN_TYPE_QUESTION_MARK
	TOKEN_TYPE_TILDE
	TOKEN_TYPE_PIPE
	TOKEN_TYPE_OPEN_OBJ
	TOKEN_TYPE_CLOSE_OBJ
	TOKEN_TYPE_DIRECTIVE
	TOKEN_TYPE_EOF
)

func (tt TokenType) String() string {
	switch tt {
	case TOKEN_TYPE_WORD:
		return "WORD"
	case TOKEN_TYPE_CONSTANT:
		return "CONSTANT"
	case TOKEN_TYPE_ASSIGN:
		return "ASSIGN"
	case TOKEN_TYPE_NUMBER_DECIMAL:
		return "NUMBER_DECIMAL"
	case TOKEN_TYPE_NUMBER_HEX:
		return "NUMBER_HEX"
	case TOKEN_TYPE_NUMBER_BINARY:
		return "NUMBER_BINARY"
	case TOKEN_TYPE_STRING:
		return "STRING"
	case TOKEN_TYPE_BOOL:
		return "BOOL"
	case TOKEN_TYPE_NULL:
		return "NULL"
	case TOKEN_TYPE_OPEN_LIST:
		return "OPEN_LIST"
	case TOKEN_TYPE_CLOSE_LIST:
		return "CLOSE_LIST"
	case TOKEN_TYPE_COMMA:
		return "COMMA"
	case TOKEN_TYPE_DOT:
		return "DOT"
	case TOKEN_TYPE_QUESTION_MARK:
		return "QUESTION_MARK"
	case TOKEN_TYPE_TILDE:
		return "TILDE"
	case TOKEN_TYPE_PIPE:
		return "PIPE"
	case TOKEN_TYPE_OPEN_OBJ:
		return "OPEN_OBJ"
	case TOKEN_TYPE_CLOSE_OBJ:
		return "CLOSE_OBJ"
	case TOKEN_TYPE_DIRECTIVE:
		return "DIRECTIVE"
	case TOKEN_TYPE_EOF:
		return "EOF"
	}

	panic("token type with no correspoding string value")
}

const (
	NO_VALUE = "NO_VALUE"
)

type Location struct {
	Line int
	Col  int
}

func (l Location) String() string {
	return fmt.Sprintf("Location{Line: %d, Col: %d}", l.Line, l.Col)
}

type Token struct {
	Type       TokenType
	Value      string
	Values     []string
	StringSubs []string
	Start      Location
}

func (t Token) String() string {
	if t.Value == NO_VALUE {
		return fmt.Sprintf(`Token{
  Type: %s,
  Location: %s
}`, t.Type, t.Start)
	}
	return fmt.Sprintf(`Token{
  Type: %s,
  Value: %s,
  Location: %s
}`, t.Type, t.Value, t.Start)
}

func WordToken(value string, start Location) Token {
	return Token{
		Type:  TOKEN_TYPE_WORD,
		Value: value,
		Start: start,
	}
}

func ConstantToken(value string, start Location) Token {
	return Token{
		Type:  TOKEN_TYPE_CONSTANT,
		Value: value,
		Start: start,
	}
}

func AssignToken(start Location) Token {
	return Token{
		Type:  TOKEN_TYPE_ASSIGN,
		Value: NO_VALUE,
		Start: start,
	}
}

func NumberToken(value string, numtype TokenType, start Location) Token {
	switch numtype {
	case TOKEN_TYPE_NUMBER_DECIMAL:
		return Token{
			Type:  TOKEN_TYPE_NUMBER_DECIMAL,
			Value: value,
			Start: start,
		}
	case TOKEN_TYPE_NUMBER_HEX:
		return Token{
			Type:  TOKEN_TYPE_NUMBER_HEX,
			Value: value,
			Start: start,
		}
	case TOKEN_TYPE_NUMBER_BINARY:
		return Token{
			Type:  TOKEN_TYPE_NUMBER_BINARY,
			Value: value,
			Start: start,
		}
	default:
		return Token{
			Type:  TOKEN_TYPE_NUMBER_DECIMAL,
			Value: value,
			Start: start,
		}
	}
}

func StringToken(values []string, stringsubs []string, start Location) Token {
	return Token{
		Type:       TOKEN_TYPE_STRING,
		Values:     values,
		StringSubs: stringsubs,
		Start:      start,
	}
}

func BoolToken(value string, start Location) Token {
	return Token{
		Type:  TOKEN_TYPE_BOOL,
		Value: value,
		Start: start,
	}
}

func NullToken(start Location) Token {
	return Token{
		Type:  TOKEN_TYPE_NULL,
		Value: NO_VALUE,
		Start: start,
	}
}

func OpenListToken(start Location) Token {
	return Token{
		Type:  TOKEN_TYPE_OPEN_LIST,
		Value: NO_VALUE,
		Start: start,
	}
}

func CloseListToken(start Location) Token {
	return Token{
		Type:  TOKEN_TYPE_CLOSE_LIST,
		Value: NO_VALUE,
		Start: start,
	}
}

func CommaToken(start Location) Token {
	return Token{
		Type:  TOKEN_TYPE_COMMA,
		Value: NO_VALUE,
		Start: start,
	}
}

func DotToken(start Location) Token {
	return Token{
		Type:  TOKEN_TYPE_DOT,
		Value: NO_VALUE,
		Start: start,
	}
}

func QuestionMarkToken(start Location) Token {
	return Token{
		Type:  TOKEN_TYPE_QUESTION_MARK,
		Value: NO_VALUE,
		Start: start,
	}
}

func TildeToken(start Location) Token {
	return Token{
		Type:  TOKEN_TYPE_TILDE,
		Value: NO_VALUE,
		Start: start,
	}
}

func PipeToken(start Location) Token {
	return Token{
		Type:  TOKEN_TYPE_PIPE,
		Value: NO_VALUE,
		Start: start,
	}
}

func OpenObjToken(start Location) Token {
	return Token{
		Type:  TOKEN_TYPE_OPEN_OBJ,
		Value: NO_VALUE,
		Start: start,
	}
}

func CloseObjToken(start Location) Token {
	return Token{
		Type:  TOKEN_TYPE_CLOSE_OBJ,
		Value: NO_VALUE,
		Start: start,
	}
}

func DirectiveToken(value string, start Location) Token {
	return Token{
		Type:  TOKEN_TYPE_DIRECTIVE,
		Value: value,
		Start: start,
	}
}

func EOFToken() Token {
	return Token{
		Type:  TOKEN_TYPE_EOF,
		Value: NO_VALUE,
		Start: Location{0, 0},
	}
}
