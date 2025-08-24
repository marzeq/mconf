package mconf_values

import (
  "fmt"
  "math/big"
  "strings"

  "github.com/marzeq/mconf/v2/mconf_tokeniser"
)

type MconfValue interface {
  ValueToString(indentAndDepth ...int) string
  ToJSONString() string
}

type MconfInt struct {
  Value *big.Int
}

func (v *MconfInt) ValueToString(indentAndDepth ...int) string {
  return v.Value.String()
}

func (v *MconfInt) ToJSONString() string {
  return v.Value.String()
}

func IntToMconfValue(i int64) *MconfInt {
  return &MconfInt{Value: big.NewInt(i)}
}

type MconfFloat struct {
  Value *big.Float
}

func (v *MconfFloat) ValueToString(indentAndDepth ...int) string {
  return v.Value.String()
}

func (v *MconfFloat) ToJSONString() string {
  return v.Value.String()
}

func FloatToMconfValue(f float64) *MconfFloat {
  return &MconfFloat{Value: big.NewFloat(f)}
}

type MconfList struct {
  Value []MconfValue
}

func (v *MconfList) OneLineStringValue() string {
  if len(v.Value) == 0 {
    return "[]"
  }

  s := "["

  for i, val := range v.Value {
    s += val.ValueToString()

    if i < len(v.Value)-1 {
      s += ", "
    }
  }

  s += "]"

  return s
}

func (v *MconfList) ToJSONString() string {
  if len(v.Value) == 0 {
    return "[]"
  }

  s := "["

  for i, val := range v.Value {
    s += val.ToJSONString()

    if i < len(v.Value)-1 {
      s += ","
    }
  }

  s += "]"

  return s
}

func (v *MconfList) ValueToString(indentAndDepth ...int) string {
  var indentSize int
  var depth int
  if len(indentAndDepth) == 0 {
    indentSize = 0
    depth = 1
  } else if len(indentAndDepth) == 1 {
    indentSize = indentAndDepth[0]
    depth = 1
  } else {
    indentSize = indentAndDepth[0]
    depth = indentAndDepth[1]
  }

  if len(v.Value) == 0 {
    return "[]"
  }

  noIndent := v.OneLineStringValue()

  if len(noIndent) < 16 || indentSize == 0 {
    return noIndent
  }

  s := "[\n"

  indent := strings.Repeat(" ", indentSize)
  currindent := strings.Repeat(indent, depth)

  for i, val := range v.Value {
    s += fmt.Sprintf("%s%s", currindent, val.ValueToString(indentSize, depth+1))

    if i < len(v.Value)-1 {
      s += ","
    }

    s += "\n"
  }

  s += fmt.Sprintf("%s]", strings.Repeat(indent, depth-1))

  return s
}

type MconfBool struct {
  Value bool
}

func (v *MconfBool) ValueToString(indentAndDepth ...int) string {
  if v.Value {
    return "true"
  }
  return "false"
}

func (v *MconfBool) ToJSONString() string {
  return v.ValueToString()
}

type MconfNull struct{}

func (v *MconfNull) ValueToString(indentAndDepth ...int) string {
  return "null"
}

func (v *MconfNull) ToJSONString() string {
  return v.ValueToString()
}

type MconfString struct {
  Value string
}

func (v *MconfString) ValueToString(indentAndDepth ...int) string {
  replaced := applyEscapes(v.Value)

  return "\"" + replaced + "\""
}

func (v *MconfString) ToJSONString() string {
  return v.ValueToString()
}

type MconfObject struct {
  Value map[string]MconfValue
  KeysOrder []string
}

func applyEscapes(s string) string {
  s = strings.ReplaceAll(s, "\\", "\\\\")
  s = strings.ReplaceAll(s, "\"", "\\\"")
  s = strings.ReplaceAll(s, "\a", "\\a")
  s = strings.ReplaceAll(s, "\b", "\\b")
  s = strings.ReplaceAll(s, "\t", "\\t")
  s = strings.ReplaceAll(s, "\n", "\\n")
  s = strings.ReplaceAll(s, "\v", "\\v")
  s = strings.ReplaceAll(s, "\f", "\\f")
  s = strings.ReplaceAll(s, "\r", "\\r")

  return s
}

func prepareKey(s string, inJson ...bool) string {
  if s == "" {
    return "\"\""
  }

  if len(inJson) > 0 && inJson[0] || !mconf_tokeniser.IsLegalWord([]rune(s)) {
    return fmt.Sprintf("\"%s\"", applyEscapes(s))
  } else {
    return s
  }
}

func (v *MconfObject) OneLineStringValue() string {
  if len(v.Value) == 0 {
    return "{}"
  }

  s := "{ "

  keycount := 0

  for _, k := range v.KeysOrder {
    val := v.Value[k]
    s += fmt.Sprintf("%s = %s", prepareKey(k), val.ValueToString())

    if keycount < len(v.Value)-1 {
      s += ", "
    }

    keycount++
  }

  s += " }"

  return s
}

func (v *MconfObject) ToJSONString() string {
  if len(v.Value) == 0 {
    return "{}"
  }

  s := "{"

  keycount := 0

  for _, k := range v.KeysOrder{
    val := v.Value[k]
    s += fmt.Sprintf("%s:%s", prepareKey(k, true), val.ToJSONString())

    if keycount < len(v.Value)-1 {
      s += ","
    }

    keycount++
  }

  s += "}"

  return s
}

func (v *MconfObject) ValueToString(indentAndDepth ...int) string {
  var indentSize int
  var depth int
  if len(indentAndDepth) == 0 {
    indentSize = 0
    depth = 1
  } else if len(indentAndDepth) == 1 {
    indentSize = indentAndDepth[0]
    depth = 1
  } else {
    indentSize = indentAndDepth[0]
    depth = indentAndDepth[1]
  }

  if len(v.Value) == 0 {
    return "{}"
  }

  noIndent := v.OneLineStringValue()

  if len(noIndent) < 16 || indentSize == 0 {
    return noIndent
  }

  s := "{\n"

  keycount := 0

  indent := strings.Repeat(" ", indentSize)
  currindent := strings.Repeat(indent, depth)

  for _, k := range v.KeysOrder {
    val := v.Value[k]
    s += fmt.Sprintf("%s%s = %s\n", currindent, prepareKey(k), val.ValueToString(indentSize, depth+1))

    keycount++
  }

  s += fmt.Sprintf("%s}", strings.Repeat(indent, depth-1))

  return s
}
