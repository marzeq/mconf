package parser

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/marzeq/mconf/tokeniser"
)

type ParserValueObject struct {
	Value map[string]ParserValue
}

func (v *ParserValueObject) GetType() string {
	return PARSER_VALUE_TYPE_OBJECT
}

func applyEscapes(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\a", "\\a")
	s = strings.ReplaceAll(s, "\b", "\\b")
	s = strings.ReplaceAll(s, "\t", "\\t")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\v", "\\v")
	s = strings.ReplaceAll(s, "\f", "\\f")
	s = strings.ReplaceAll(s, "\r", "\\r")

	return s
}

func prepareKey(s string, inJson ...bool) string {
	if s == "" {
		return "\"\""
	}

	if len(inJson) > 0 && inJson[0] || !tokeniser.IsLegalWord([]rune(s)) {
		return fmt.Sprintf("\"%s\"", applyEscapes(s))
	} else {
		return s
	}
}

func (v *ParserValueObject) OneLineStringValue() string {
	if len(v.Value) == 0 {
		return "{}"
	}

	s := "{ "

	keycount := 0

	for k, val := range v.Value {
		s += fmt.Sprintf("%s = %s", prepareKey(k), val.ValueToString())

		if keycount < len(v.Value)-1 {
			s += ", "
		}

		keycount++
	}

	s += " }"

	return s
}

func (v *ParserValueObject) ToJSONString() string {
	if len(v.Value) == 0 {
		return "{}"
	}

	s := "{"

	keycount := 0

	for k, val := range v.Value {
		s += fmt.Sprintf("%s:%s", prepareKey(k, true), val.ToJSONString())

		if keycount < len(v.Value)-1 {
			s += ","
		}

		keycount++
	}

	s += "}"

	return s
}

func (v *ParserValueObject) ValueToString(indentAndDepth ...int) string {
	var indentSize int
	var depth int
	if len(indentAndDepth) == 0 {
		indentSize = 0
		depth = 1
	} else if len(indentAndDepth) == 1 {
		indentSize = indentAndDepth[0]
		depth = 1
	} else {
		indentSize = indentAndDepth[0]
		depth = indentAndDepth[1]
	}

	if len(v.Value) == 0 {
		return "{}"
	}

	noIndent := v.OneLineStringValue()

	if len(noIndent) < 16 || indentSize == 0 {
		return noIndent
	}

	s := "{\n"

	keycount := 0

	indent := strings.Repeat(" ", indentSize)
	currindent := strings.Repeat(indent, depth)

	for k, val := range v.Value {
		s += fmt.Sprintf("%s%s = %s\n", currindent, prepareKey(k), val.ValueToString(indentSize, depth+1))

		keycount++
	}

	s += fmt.Sprintf("%s}", strings.Repeat(indent, depth-1))

	return s
}

func (v *ParserValueObject) GetObject() (map[string]ParserValue, error) {
	return v.Value, nil
}

func (v *ParserValueObject) GetString() (string, error) {
	return "", WrongTypeError(PARSER_VALUE_TYPE_STRING, PARSER_VALUE_TYPE_OBJECT)
}

func (v *ParserValueObject) GetFloat() (*big.Float, error) {
	return nil, WrongTypeError(PARSER_VALUE_TYPE_FLOAT, PARSER_VALUE_TYPE_OBJECT)
}

func (v *ParserValueObject) GetInt() (*big.Int, error) {
	return nil, WrongTypeError(PARSER_VALUE_TYPE_INT, PARSER_VALUE_TYPE_OBJECT)
}

func (v *ParserValueObject) GetBool() (bool, error) {
	return false, WrongTypeError(PARSER_VALUE_TYPE_BOOL, PARSER_VALUE_TYPE_OBJECT)
}

func (v *ParserValueObject) GetList() ([]ParserValue, error) {
	return nil, WrongTypeError(PARSER_VALUE_TYPE_LIST, PARSER_VALUE_TYPE_OBJECT)
}
