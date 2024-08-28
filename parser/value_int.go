package parser

import (
	"math/big"
)

type ParserValueInt struct {
	Value *big.Int
}

func (v *ParserValueInt) GetType() string {
	return PARSER_VALUE_TYPE_INT
}

func (v *ParserValueInt) ValueToString(indentAndDepth ...int) string {
	return v.Value.String()
}

func (v *ParserValueInt) GetInt() (*big.Int, error) {
	return v.Value, nil
}

func (v *ParserValueInt) GetFloat() (*big.Float, error) {
	return new(big.Float).SetInt(v.Value), nil
}

func (v *ParserValueInt) GetString() (string, error) {
	return "", WrongTypeError(PARSER_VALUE_TYPE_STRING, PARSER_VALUE_TYPE_FLOAT)
}

func (v *ParserValueInt) GetBool() (bool, error) {
	return false, WrongTypeError(PARSER_VALUE_TYPE_BOOL, PARSER_VALUE_TYPE_FLOAT)
}

func (v *ParserValueInt) GetList() ([]ParserValue, error) {
	return nil, WrongTypeError(PARSER_VALUE_TYPE_LIST, PARSER_VALUE_TYPE_FLOAT)
}

func (v *ParserValueInt) GetObject() (map[string]ParserValue, error) {
	return nil, WrongTypeError(PARSER_VALUE_TYPE_OBJECT, PARSER_VALUE_TYPE_FLOAT)
}
