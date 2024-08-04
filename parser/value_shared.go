package parser

import (
  "fmt"
  "os"
)

func WrongTypeError(attemptedType string, actualType string) {
  fmt.Printf("Tried to get value of type %s, but the underlying value is of type %s\n", attemptedType, actualType)
  os.Exit(1)
}
