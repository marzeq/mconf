package parser

import (
  "strings"
  "fmt"
)

type ParserValueList struct {
  Value []ParserValue
}

func (v *ParserValueList) GetType() string {
  return PARSER_VALUE_TYPE_LIST
}

func (v *ParserValueList) ValueToString() string {
  if len(v.Value) == 0 {
    return "[]"
  }

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

func (v *ParserValueList) IndentedString(indent string, depth int) string {
  if len(v.Value) == 0 {
    return "[]"
  }

  nonIndented := v.ValueToString()

  if len(nonIndented) < 16 {
    return nonIndented
  }

  s := "[\n"
  
  currindent := strings.Repeat(indent, depth)

  for i, val := range v.Value {
    if val.GetType() == PARSER_VALUE_TYPE_OBJECT {
      objval := val.(*ParserValueObject)

      s += fmt.Sprintf("%s%s", currindent, objval.IndentedString(indent, depth + 1))
    } else if val.GetType() == PARSER_VALUE_TYPE_LIST {
      listval := val.(*ParserValueList)

      s += fmt.Sprintf("%s%s", currindent, listval.IndentedString(indent, depth + 1))
    } else {
      s += fmt.Sprintf("%s%s", currindent, val.ValueToString())
    }

    if i < len(v.Value) - 1 {
      s += ","
    }

    s += "\n"
  }

  s += fmt.Sprintf("%s]", strings.Repeat(indent, depth - 1))

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

