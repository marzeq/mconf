package parser

type ParserValueBool struct {
  Value bool
}

func (v *ParserValueBool) GetType() string {
  return PARSER_VALUE_TYPE_BOOL
}

func (v *ParserValueBool) GetBool() bool {
  return v.Value
}

func (v *ParserValueBool) GetString() string {
  WrongTypeError(PARSER_VALUE_TYPE_STRING, PARSER_VALUE_TYPE_BOOL)
  return ""
}

func (v *ParserValueBool) GetNumber() float64 {
  WrongTypeError(PARSER_VALUE_TYPE_NUMBER, PARSER_VALUE_TYPE_BOOL)
  return 0
}

func (v *ParserValueBool) GetList() []ParserValue {
  WrongTypeError(PARSER_VALUE_TYPE_LIST, PARSER_VALUE_TYPE_BOOL)
  return nil
}

func (v *ParserValueBool) GetObject() map[string]ParserValue {
  WrongTypeError(PARSER_VALUE_TYPE_OBJECT, PARSER_VALUE_TYPE_BOOL)
  return nil
}
