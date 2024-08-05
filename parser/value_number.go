package parser

import "strconv"

type ParserValueNumber struct {
  Value float64
}

func (v *ParserValueNumber) GetType() string {
  return PARSER_VALUE_TYPE_NUMBER
}

func (v *ParserValueNumber) ValueToString(indentAndDepth ...int) string {
  return strconv.FormatFloat(v.Value, 'f', -1, 64)
}

func (v *ParserValueNumber) GetNumber() (float64, error) {
  return v.Value, nil
}

func (v *ParserValueNumber) GetString() (string, error) {
  return "", WrongTypeError(PARSER_VALUE_TYPE_STRING, PARSER_VALUE_TYPE_NUMBER)
}

func (v *ParserValueNumber) GetBool() (bool, error) {
  return false, WrongTypeError(PARSER_VALUE_TYPE_BOOL, PARSER_VALUE_TYPE_NUMBER)
}

func (v *ParserValueNumber) GetList() ([]ParserValue, error) {
  return nil, WrongTypeError(PARSER_VALUE_TYPE_LIST, PARSER_VALUE_TYPE_NUMBER)
}

func (v *ParserValueNumber) GetObject() (map[string]ParserValue, error) {
  return nil, WrongTypeError(PARSER_VALUE_TYPE_OBJECT, PARSER_VALUE_TYPE_NUMBER)
}
