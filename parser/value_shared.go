package parser

import (
	"fmt"
)

func WrongTypeError(attemptedType ValueType, actualType ValueType) error {
	return fmt.Errorf("Tried to get value of type %s, but the underlying value is of type %s\n", attemptedType, actualType)
}
