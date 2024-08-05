package parser

import (
	"strings"
)

type ParserValueString struct {
  Value string
}

func (v *ParserValueString) GetType() string {
  return PARSER_VALUE_TYPE_STRING
}

func (v *ParserValueString) ValueToString(indentAndDepth ...int) string {
  replaced := strings.ReplaceAll(v.Value, "\"", "\\\"")
  replaced = strings.ReplaceAll(replaced, "\n", "\\n")
  replaced = strings.ReplaceAll(replaced, "\r", "\\r")
  replaced = strings.ReplaceAll(replaced, "\t", "\\t")
  replaced = strings.ReplaceAll(replaced, "\f", "\\f")

  return "\"" + replaced + "\""
}

func (v *ParserValueString) GetString() (string, error) {
  return v.Value, nil
}

func (v *ParserValueString) GetNumber() (float64, error) {
  return 0, WrongTypeError(PARSER_VALUE_TYPE_NUMBER, PARSER_VALUE_TYPE_STRING)
}

func (v *ParserValueString) GetBool() (bool, error) {
  return false, WrongTypeError(PARSER_VALUE_TYPE_BOOL, PARSER_VALUE_TYPE_STRING)
}

func (v *ParserValueString) GetList() ([]ParserValue, error) {
  return nil, WrongTypeError(PARSER_VALUE_TYPE_LIST, PARSER_VALUE_TYPE_STRING)
}

func (v *ParserValueString) GetObject() (map[string]ParserValue, error) {
  return nil, WrongTypeError(PARSER_VALUE_TYPE_OBJECT, PARSER_VALUE_TYPE_STRING)
}
