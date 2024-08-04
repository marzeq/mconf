package parser

type ParserValueObject struct {
  Value map[string]ParserValue
}

func (v *ParserValueObject) GetType() string {
  return PARSER_VALUE_TYPE_LIST
}

func (v *ParserValueObject) GetObject() map[string]ParserValue {
  return v.Value
}

func (v *ParserValueObject) GetString() string {
  WrongTypeError(PARSER_VALUE_TYPE_STRING, PARSER_VALUE_TYPE_LIST)
  return ""
}

func (v *ParserValueObject) GetNumber() float64 {
  WrongTypeError(PARSER_VALUE_TYPE_NUMBER, PARSER_VALUE_TYPE_LIST)
  return 0
}

func (v *ParserValueObject) GetBool() bool {
  WrongTypeError(PARSER_VALUE_TYPE_BOOL, PARSER_VALUE_TYPE_LIST)
  return false
}

func (v *ParserValueObject) GetList() map[string]ParserValue {
  WrongTypeError(PARSER_VALUE_TYPE_LIST, PARSER_VALUE_TYPE_LIST)
  return nil
}

