package parser

import (
	"math/big"
)

type ParserValueFloat struct {
	Value *big.Float
}

func (v *ParserValueFloat) GetType() string {
	return PARSER_VALUE_TYPE_FLOAT
}

func (v *ParserValueFloat) IsNull() bool {
	return false
}

func (v *ParserValueFloat) ValueToString(indentAndDepth ...int) string {
	return v.Value.String()
}

func (v *ParserValueFloat) ToJSONString() string {
	return v.Value.String()
}

func (v *ParserValueFloat) GetFloat() (*big.Float, error) {
	return v.Value, nil
}

func (v *ParserValueFloat) GetInt() (*big.Int, error) {
	i, _ := v.Value.Int(nil)
	return i, nil
}

func (v *ParserValueFloat) GetString() (string, error) {
	return "", WrongTypeError(PARSER_VALUE_TYPE_STRING, PARSER_VALUE_TYPE_FLOAT)
}

func (v *ParserValueFloat) GetBool() (bool, error) {
	return false, WrongTypeError(PARSER_VALUE_TYPE_BOOL, PARSER_VALUE_TYPE_FLOAT)
}

func (v *ParserValueFloat) GetList() ([]ParserValue, error) {
	return nil, WrongTypeError(PARSER_VALUE_TYPE_LIST, PARSER_VALUE_TYPE_FLOAT)
}

func (v *ParserValueFloat) GetObject() (map[string]ParserValue, error) {
	return nil, WrongTypeError(PARSER_VALUE_TYPE_OBJECT, PARSER_VALUE_TYPE_FLOAT)
}
