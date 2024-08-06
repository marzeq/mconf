package parser

import "strconv"

type ParserValueUInt struct {
  Value uint64 
}

func (v *ParserValueUInt) GetType() string {
  return PARSER_VALUE_TYPE_UINT
}

func (v *ParserValueUInt) ValueToString(indentAndDepth ...int) string {
  return strconv.FormatUint(v.Value, 10)
}

func (v *ParserValueUInt) GetUInt() (uint64, error) {
  return v.Value, nil
}

func (v *ParserValueUInt) GetFloat() (float64, error) {
  return float64(v.Value), nil
}

func (v* ParserValueUInt) GetInt() (int64, error) {
  return int64(v.Value), nil
}

func (v *ParserValueUInt) GetString() (string, error) {
  return "", WrongTypeError(PARSER_VALUE_TYPE_STRING, PARSER_VALUE_TYPE_UINT)
}

func (v *ParserValueUInt) GetBool() (bool, error) {
  return false, WrongTypeError(PARSER_VALUE_TYPE_BOOL, PARSER_VALUE_TYPE_UINT)
}

func (v *ParserValueUInt) GetList() ([]ParserValue, error) {
  return nil, WrongTypeError(PARSER_VALUE_TYPE_LIST, PARSER_VALUE_TYPE_UINT)
}

func (v *ParserValueUInt) GetObject() (map[string]ParserValue, error) {
  return nil, WrongTypeError(PARSER_VALUE_TYPE_OBJECT, PARSER_VALUE_TYPE_UINT)
}
