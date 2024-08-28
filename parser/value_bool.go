package parser

import (
	"math/big"
	"strconv"
)

type ParserValueBool struct {
	Value bool
}

func (v *ParserValueBool) GetType() string {
	return PARSER_VALUE_TYPE_BOOL
}

func (v *ParserValueBool) ValueToString(indentAndDepth ...int) string {
	return strconv.FormatBool(v.Value)
}

func (v *ParserValueBool) GetBool() (bool, error) {
	return v.Value, nil
}

func (v *ParserValueBool) GetString() (string, error) {
	return "", WrongTypeError(PARSER_VALUE_TYPE_STRING, PARSER_VALUE_TYPE_BOOL)
}

func (v *ParserValueBool) GetFloat() (*big.Float, error) {
	return nil, WrongTypeError(PARSER_VALUE_TYPE_FLOAT, PARSER_VALUE_TYPE_BOOL)
}

func (v *ParserValueBool) GetInt() (*big.Int, error) {
	return nil, WrongTypeError(PARSER_VALUE_TYPE_INT, PARSER_VALUE_TYPE_BOOL)
}

func (v *ParserValueBool) GetList() ([]ParserValue, error) {
	return nil, WrongTypeError(PARSER_VALUE_TYPE_LIST, PARSER_VALUE_TYPE_BOOL)
}

func (v *ParserValueBool) GetObject() (map[string]ParserValue, error) {
	return nil, WrongTypeError(PARSER_VALUE_TYPE_OBJECT, PARSER_VALUE_TYPE_BOOL)
}
