package parser

import "fmt"

type ParserValueObject struct {
  Value map[string]ParserValue
}

func (v *ParserValueObject) GetType() string {
  return PARSER_VALUE_TYPE_LIST
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

func (v *ParserValueObject) GetObject() (map[string]ParserValue, error) {
  return v.Value, nil
}

func (v *ParserValueObject) GetString() (string, error) {
  return "", WrongTypeError(PARSER_VALUE_TYPE_STRING, PARSER_VALUE_TYPE_LIST)
}

func (v *ParserValueObject) GetNumber() (float64, error) {
  return 0, WrongTypeError(PARSER_VALUE_TYPE_NUMBER, PARSER_VALUE_TYPE_LIST)
}

func (v *ParserValueObject) GetBool() (bool, error) {
  return false, WrongTypeError(PARSER_VALUE_TYPE_BOOL, PARSER_VALUE_TYPE_LIST)
}

func (v *ParserValueObject) GetList() ([]ParserValue, error) {
  return nil, WrongTypeError(PARSER_VALUE_TYPE_LIST, PARSER_VALUE_TYPE_LIST)
}
