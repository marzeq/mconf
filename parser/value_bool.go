package parser

import "strconv"

type ParserValueBool struct {
  Value bool
}

func (v *ParserValueBool) GetType() string {
  return PARSER_VALUE_TYPE_BOOL
}

func (v *ParserValueBool) ValueToString(indentAndDepth ...int) string {
  return strconv.FormatBool(v.Value)
}

func (v *ParserValueBool) GetBool() (bool, error) {
  return v.Value, nil
}

func (v *ParserValueBool) GetString() (string, error) {
  return "", WrongTypeError(PARSER_VALUE_TYPE_STRING, PARSER_VALUE_TYPE_BOOL)
}

func (v *ParserValueBool) GetNumber() (float64, error) {
  return 0, WrongTypeError(PARSER_VALUE_TYPE_NUMBER, PARSER_VALUE_TYPE_BOOL)
}

func (v *ParserValueBool) GetList() ([]ParserValue, error) {
  return nil, WrongTypeError(PARSER_VALUE_TYPE_LIST, PARSER_VALUE_TYPE_BOOL)
}

func (v *ParserValueBool) GetObject() (map[string]ParserValue, error) {
  return nil, WrongTypeError(PARSER_VALUE_TYPE_OBJECT, PARSER_VALUE_TYPE_BOOL)
}
