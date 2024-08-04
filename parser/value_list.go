package parser

type ParserValueList struct {
  Value []ParserValue
}

func (v *ParserValueList) GetType() string {
  return PARSER_VALUE_TYPE_LIST
}

func (v *ParserValueList) ValueToString() string {
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

