package tokeniser

import (
	"fmt"
)

const (
	TOKEN_TYPE_KEY            = "KEY"
	TOKEN_TYPE_CONSTANT       = "CONSTANT"
	TOKEN_TYPE_ASSIGN         = "ASSIGN"
	TOKEN_TYPE_NUMBER_DECIMAL = "NUMBER_DECIMAL"
	TOKEN_TYPE_NUMBER_HEX     = "NUMBER_HEX"
	TOKEN_TYPE_NUMBER_BINARY  = "NUMBER_BINARY"
	TOKEN_TYPE_STRING         = "STRING"
	TOKEN_TYPE_BOOL           = "BOOL"
	TOKEN_TYPE_OPEN_LIST      = "OPEN_LIST"
	TOKEN_TYPE_CLOSE_LIST     = "CLOSE_LIST"
	TOKEN_TYPE_COMMA          = "COMMA"
	TOKEN_TYPE_DOT            = "DOT"
	TOKEN_TYPE_QUESTION_MARK  = "QUESTION_MARK"
	TOKEN_TYPE_OPEN_OBJ       = "OPEN_OBJ"
	TOKEN_TYPE_CLOSE_OBJ      = "CLOSE_OBJ"
	TOKEN_TYPE_DIRECTIVE      = "DIRECTIVE"

	TOKEN_TYPE_EOF = "EOF"
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
	Type  string
	Value string
	Start Location
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
		Type:  TOKEN_TYPE_KEY,
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

func NumberToken(value string, numtype string, start Location) Token {
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

func StringToken(value string, start Location) Token {
	return Token{
		Type:  TOKEN_TYPE_STRING,
		Value: value,
		Start: start,
	}
}

func BoolToken(value string, start Location) Token {
	return Token{
		Type:  TOKEN_TYPE_BOOL,
		Value: value,
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
