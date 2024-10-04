package parser

import (
	"math/big"
)

type ParserValueNull struct {
	Value bool
}

func (v *ParserValueNull) IsNull() bool {
	return true
}

func (v *ParserValueNull) GetType() string {
	return PARSER_VALUE_TYPE_NULL
}

func (v *ParserValueNull) ValueToString(indentAndDepth ...int) string {
	return "null"
}

func (v *ParserValueNull) ToJSONString() string {
	return v.ValueToString()
}

func (v *ParserValueNull) GetBool() (bool, error) {
	return false, WrongTypeError(PARSER_VALUE_TYPE_BOOL, PARSER_VALUE_TYPE_NULL)
}

func (v *ParserValueNull) GetString() (string, error) {
	return "", WrongTypeError(PARSER_VALUE_TYPE_STRING, PARSER_VALUE_TYPE_BOOL)
}

func (v *ParserValueNull) GetFloat() (*big.Float, error) {
	return nil, WrongTypeError(PARSER_VALUE_TYPE_FLOAT, PARSER_VALUE_TYPE_BOOL)
}

func (v *ParserValueNull) GetInt() (*big.Int, error) {
	return nil, WrongTypeError(PARSER_VALUE_TYPE_INT, PARSER_VALUE_TYPE_BOOL)
}

func (v *ParserValueNull) GetList() ([]ParserValue, error) {
	return nil, WrongTypeError(PARSER_VALUE_TYPE_LIST, PARSER_VALUE_TYPE_BOOL)
}

func (v *ParserValueNull) GetObject() (map[string]ParserValue, error) {
	return nil, WrongTypeError(PARSER_VALUE_TYPE_OBJECT, PARSER_VALUE_TYPE_BOOL)
}
