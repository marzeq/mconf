package parser

type ParserValueNumber struct {
  Value float64
}

func (v *ParserValueNumber) GetType() string {
  return PARSER_VALUE_TYPE_NUMBER
}

func (v *ParserValueNumber) GetNumber() float64 {
  return v.Value
}

func (v *ParserValueNumber) GetString() string {
  WrongTypeError(PARSER_VALUE_TYPE_STRING, PARSER_VALUE_TYPE_NUMBER)
  return ""
}

func (v *ParserValueNumber) GetBool() bool {
  WrongTypeError(PARSER_VALUE_TYPE_BOOL, PARSER_VALUE_TYPE_NUMBER)
  return false
}

func (v *ParserValueNumber) GetList() []ParserValue {
  WrongTypeError(PARSER_VALUE_TYPE_LIST, PARSER_VALUE_TYPE_NUMBER)
  return nil
}

func (v *ParserValueNumber) GetObject() map[string]ParserValue {
  WrongTypeError(PARSER_VALUE_TYPE_OBJECT, PARSER_VALUE_TYPE_NUMBER)
  return nil
}
