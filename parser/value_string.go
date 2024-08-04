package parser

type ParserValueString struct {
  Value string
}

func (v *ParserValueString) GetType() string {
  return PARSER_VALUE_TYPE_STRING
}

func (v *ParserValueString) GetString() string {
  return v.Value
}

func (v *ParserValueString) GetNumber() float64 {
  WrongTypeError(PARSER_VALUE_TYPE_NUMBER, PARSER_VALUE_TYPE_STRING)
  return 0
}

func (v *ParserValueString) GetBool() bool {
  WrongTypeError(PARSER_VALUE_TYPE_BOOL, PARSER_VALUE_TYPE_STRING)
  return false
}

func (v *ParserValueString) GetList() []ParserValue {
  WrongTypeError(PARSER_VALUE_TYPE_LIST, PARSER_VALUE_TYPE_STRING)
  return nil
}

func (v *ParserValueString) GetObject() map[string]ParserValue {
  WrongTypeError(PARSER_VALUE_TYPE_OBJECT, PARSER_VALUE_TYPE_STRING)
  return nil
}
