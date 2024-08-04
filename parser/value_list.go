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

func (v *ParserValueList) GetList() []ParserValue {
  return v.Value
}

func (v *ParserValueList) GetString() string {
  WrongTypeError(PARSER_VALUE_TYPE_STRING, PARSER_VALUE_TYPE_LIST)
  return ""
}

func (v *ParserValueList) GetNumber() float64 {
  WrongTypeError(PARSER_VALUE_TYPE_NUMBER, PARSER_VALUE_TYPE_LIST)
  return 0
}

func (v *ParserValueList) GetBool() bool {
  WrongTypeError(PARSER_VALUE_TYPE_BOOL, PARSER_VALUE_TYPE_LIST)
  return false
}

func (v *ParserValueList) GetObject() map[string]ParserValue {
  WrongTypeError(PARSER_VALUE_TYPE_OBJECT, PARSER_VALUE_TYPE_LIST)
  return nil
}

