package parser

import (
	"fmt"
	"strings"
)

type ParserValueObject struct {
  Value map[string]ParserValue
}

func (v *ParserValueObject) GetType() string {
  return PARSER_VALUE_TYPE_OBJECT
}

func (v *ParserValueObject) ValueToString() string {
  if len(v.Value) == 0 {
    return "{}"
  }

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

func (v *ParserValueObject) IndentedString(indent string, depth int) string {
  if len(v.Value) == 0 {
    return "{}"
  }

  s := "{\n"

  keycount := 0

  currindent := strings.Repeat(indent, depth)

  for k, val := range v.Value {
    if val.GetType() == PARSER_VALUE_TYPE_OBJECT {
      objval := val.(*ParserValueObject)

      s += fmt.Sprintf("%s%s = %s\n", currindent, k, objval.IndentedString(indent, depth + 1))
    } else if val.GetType() == PARSER_VALUE_TYPE_LIST {
      listval := val.(*ParserValueList)

      s += fmt.Sprintf("%s%s = %s\n", currindent, k, listval.IndentedString(indent, depth + 1))
    } else {
      s += fmt.Sprintf("%s%s = %s\n", currindent, k, val.ValueToString())
    }

    keycount++
  }

  s += fmt.Sprintf("%s}", strings.Repeat(indent, depth - 1))

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
