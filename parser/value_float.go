package parser

import "strconv"

type ParserValueFloat struct {
  Value float64
}

func (v *ParserValueFloat) GetType() string {
  return PARSER_VALUE_TYPE_FLOAT
}

func (v *ParserValueFloat) ValueToString(indentAndDepth ...int) string {
  return strconv.FormatFloat(v.Value, 'f', -1, 64)
}

func (v *ParserValueFloat) GetFloat() (float64, error) {
  return v.Value, nil
}

func (v* ParserValueFloat) GetInt() (int64, error) {
  return int64(v.Value), nil
}

func (v *ParserValueFloat) GetUInt() (uint64, error) {
  return uint64(v.Value), nil
}

func (v *ParserValueFloat) GetString() (string, error) {
  return "", WrongTypeError(PARSER_VALUE_TYPE_STRING, PARSER_VALUE_TYPE_FLOAT)
}

func (v *ParserValueFloat) GetBool() (bool, error) {
  return false, WrongTypeError(PARSER_VALUE_TYPE_BOOL, PARSER_VALUE_TYPE_FLOAT)
}

func (v *ParserValueFloat) GetList() ([]ParserValue, error) {
  return nil, WrongTypeError(PARSER_VALUE_TYPE_LIST, PARSER_VALUE_TYPE_FLOAT)
}

func (v *ParserValueFloat) GetObject() (map[string]ParserValue, error) {
  return nil, WrongTypeError(PARSER_VALUE_TYPE_OBJECT, PARSER_VALUE_TYPE_FLOAT)
}
