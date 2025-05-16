package parser

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/marzeq/mconf/tokeniser"
)

type ParserValueInt struct {
	Value *big.Int
}

func (v *ParserValueInt) ValueToString(indentAndDepth ...int) string {
	return v.Value.String()
}

func (v *ParserValueInt) ToJSONString() string {
	return v.Value.String()
}

type ParserValueFloat struct {
	Value *big.Float
}

func (v *ParserValueFloat) ValueToString(indentAndDepth ...int) string {
	return v.Value.String()
}

func (v *ParserValueFloat) ToJSONString() string {
	return v.Value.String()
}

type ParserValueList struct {
	Value []ParserValue
}

func (v *ParserValueList) OneLineStringValue() string {
	if len(v.Value) == 0 {
		return "[]"
	}

	s := "["

	for i, val := range v.Value {
		s += val.ValueToString()

		if i < len(v.Value)-1 {
			s += ", "
		}
	}

	s += "]"

	return s
}

func (v *ParserValueList) ToJSONString() string {
	if len(v.Value) == 0 {
		return "[]"
	}

	s := "["

	for i, val := range v.Value {
		s += val.ToJSONString()

		if i < len(v.Value)-1 {
			s += ","
		}
	}

	s += "]"

	return s
}

func (v *ParserValueList) ValueToString(indentAndDepth ...int) string {
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
		return "[]"
	}

	noIndent := v.OneLineStringValue()

	if len(noIndent) < 16 || indentSize == 0 {
		return noIndent
	}

	s := "[\n"

	indent := strings.Repeat(" ", indentSize)
	currindent := strings.Repeat(indent, depth)

	for i, val := range v.Value {
		s += fmt.Sprintf("%s%s", currindent, val.ValueToString(indentSize, depth+1))

		if i < len(v.Value)-1 {
			s += ","
		}

		s += "\n"
	}

	s += fmt.Sprintf("%s]", strings.Repeat(indent, depth-1))

	return s
}

type ParserValueBool struct {
	Value bool
}

func (v *ParserValueBool) ValueToString(indentAndDepth ...int) string {
	if v.Value {
		return "true"
	}
	return "false"
}

func (v *ParserValueBool) ToJSONString() string {
	return v.ValueToString()
}

type ParserValueNull struct{}

func (v *ParserValueNull) ValueToString(indentAndDepth ...int) string {
	return "null"
}

func (v *ParserValueNull) ToJSONString() string {
	return v.ValueToString()
}

type ParserValueString struct {
	Value string
}

func (v *ParserValueString) ValueToString(indentAndDepth ...int) string {
	replaced := applyEscapes(v.Value)

	return "\"" + replaced + "\""
}

func (v *ParserValueString) ToJSONString() string {
	return v.ValueToString()
}

type ParserValueObject struct {
	Value map[string]ParserValue
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
