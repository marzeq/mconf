package parser

import (
  "strings"
  "fmt"
)

type ParserValueList struct {
  Value []ParserValue
}

func (v *ParserValueList) GetType() string {
  return PARSER_VALUE_TYPE_LIST
}

func (v *ParserValueList) OneLineStringValue() string {
  if len(v.Value) == 0 {
    return "[]"
  }

  s := "["

  for i, val := range v.Value {
    s += val.ValueToString()

    if i < len(v.Value) - 1 {
      s += ", "
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
    return "{}"
  }

  noIndent := v.OneLineStringValue()

  if len(noIndent) < 16 || indentSize == 0 {
    return noIndent
  }

  s := "[\n"
  
  indent := strings.Repeat(" ", indentSize)
  currindent := strings.Repeat(indent, depth)

  for i, val := range v.Value {
    s += fmt.Sprintf("%s%s", currindent, val.ValueToString(indentSize, depth + 1))

    if i < len(v.Value) - 1 {
      s += ","
    }

    s += "\n"
  }

  s += fmt.Sprintf("%s]", strings.Repeat(indent, depth - 1))

  return s
}

func (v *ParserValueList) GetList() ([]ParserValue, error) {
  return v.Value, nil
}

func (v *ParserValueList) GetString() (string, error) {
  return "", WrongTypeError(PARSER_VALUE_TYPE_STRING, PARSER_VALUE_TYPE_LIST)
}

func (v *ParserValueList) GetNumber() (float64, error) {
  return 0, WrongTypeError(PARSER_VALUE_TYPE_NUMBER, PARSER_VALUE_TYPE_LIST)
}

func (v *ParserValueList) GetBool() (bool, error) {
  return false, WrongTypeError(PARSER_VALUE_TYPE_BOOL, PARSER_VALUE_TYPE_LIST)
}

func (v *ParserValueList) GetObject() (map[string]ParserValue, error) {
  return nil, WrongTypeError(PARSER_VALUE_TYPE_OBJECT, PARSER_VALUE_TYPE_LIST)
}

