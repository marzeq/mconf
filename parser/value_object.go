package parser

import "fmt"

type ParserValueObject struct {
  Value map[string]ParserValue
}

func (v *ParserValueObject) GetType() string {
  return PARSER_VALUE_TYPE_LIST
}

func (v *ParserValueObject) GetObject() map[string]ParserValue {
  return v.Value
}

func (v *ParserValueObject) ValueToString() string {
  s := "{"
  
  keycount := 0

  for k, val := range v.Value {
    s += fmt.Sprintf("%s = %s", k, val.ValueToString())

    if keycount < len(v.Value) - 1 {
      s += ", "
    }

    keycount++
  }

  s += "}"

  return s
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

func (v *ParserValueObject) GetList() []ParserValue {
  WrongTypeError(PARSER_VALUE_TYPE_LIST, PARSER_VALUE_TYPE_LIST)
  return nil
}

