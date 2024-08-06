package parser

import "strconv"

type ParserValueInt struct {
  Value int64 
}

func (v *ParserValueInt) GetType() string {
  return PARSER_VALUE_TYPE_INT
}

func (v *ParserValueInt) ValueToString(indentAndDepth ...int) string {
  return strconv.FormatInt(v.Value, 10)
}

func (v* ParserValueInt) GetInt() (int64, error) {
  return v.Value, nil
}

func (v *ParserValueInt) GetFloat() (float64, error) {
  return float64(v.Value), nil
}

func (v *ParserValueInt) GetUInt() (uint64, error) {
  return uint64(v.Value), nil
}

func (v *ParserValueInt) GetString() (string, error) {
  return "", WrongTypeError(PARSER_VALUE_TYPE_STRING, PARSER_VALUE_TYPE_FLOAT)
}

func (v *ParserValueInt) GetBool() (bool, error) {
  return false, WrongTypeError(PARSER_VALUE_TYPE_BOOL, PARSER_VALUE_TYPE_FLOAT)
}

func (v *ParserValueInt) GetList() ([]ParserValue, error) {
  return nil, WrongTypeError(PARSER_VALUE_TYPE_LIST, PARSER_VALUE_TYPE_FLOAT)
}

func (v *ParserValueInt) GetObject() (map[string]ParserValue, error) {
  return nil, WrongTypeError(PARSER_VALUE_TYPE_OBJECT, PARSER_VALUE_TYPE_FLOAT)
}
