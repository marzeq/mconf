package main

import (
	"fmt"
  "io"
	"os"
  "strconv"

	"github.com/marzeq/mconf/parser"
	"github.com/marzeq/mconf/tokeniser"
)

func check(err error) {
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func ParseFromString(s string) (map[string]parser.ParserValue, error) {
  t := tokeniser.NewTokeniser(s)
  tokens, err := t.Tokenise()
  
  if err != nil { return nil, err }

  p := parser.NewParser(tokens)

  return p.Parse()
}

func ParseFromFile(filename string) (map[string]parser.ParserValue, error) {
  f, err := os.ReadFile(filename)
  
  if err != nil {
    return nil, err
  }

  s := string(f)

  return ParseFromString(s)
}

func ParseFromBytes(b []byte) (map[string]parser.ParserValue, error) {
  s := string(b)

  return ParseFromString(s)
}

func readStdin() ([]byte, error) {
  b, err := io.ReadAll(os.Stdin)

	if err != nil { return nil, err }

  return b, nil
}

type options struct {
  Filename string
  AcessedProperties []string
}

func usage(progname string) string {
  return fmt.Sprintf(`Usage:
  %s <filename> [properties ...]
  %s - [properties ...]

Arguments:
  <filename>        Path to the configuration file. Use '-' to read from stdin.
  [properties ...]  List of properties to access. Multiple properties are used to access nested objects or lists. If no properties are provided, the global object is printed.

Examples:
  %s config.mconf property1 property2
  cat config.mconf | %s - property1 property2`, progname, progname, progname, progname)
}

func parseOptions() (options, string) {
  if len(os.Args) == 1 {
    return options{}, usage(os.Args[0])
  }

  return options{
    Filename: os.Args[1],
    AcessedProperties: os.Args[2:],
  }, ""
}

func main() {
  opts, usage := parseOptions()
  
  if usage != "" {
    fmt.Println(usage)
    return
  }

  var m map[string]parser.ParserValue
  var parsingErr error

  if opts.Filename == "-" {
    b, err := readStdin()
    check(err)

    m, parsingErr = ParseFromBytes(b)
  } else {
    m, parsingErr = ParseFromFile(opts.Filename)
  }

  check(parsingErr)

  var indexedValue parser.ParserValue = &parser.ParserValueObject{Value: m}

  indexedString := ""

  for _, p := range opts.AcessedProperties {
    if indexedString == "" {
      indexedString = p
    } else {
      indexedString += "." + p
    }

    if indexedValue.GetType() != parser.PARSER_VALUE_TYPE_OBJECT && indexedValue.GetType() != parser.PARSER_VALUE_TYPE_LIST {
      fmt.Printf("Property %s not found, indexed value is not an object or list\n", indexedString)
      return
    }

    if indexedValue.GetType() == parser.PARSER_VALUE_TYPE_OBJECT {
      obj, err := indexedValue.GetObject()

      if err != nil {
        fmt.Printf("Unexpected error, indexed value has type object but cannot be converted to object, please report this bug\n")
        return
      }

      next := obj[p]

      if next == nil {
        fmt.Printf("Property %s not found\n", indexedString)
        return
      }

      indexedValue = next
    } else {
      list, err := indexedValue.GetList()

      if err != nil {
        fmt.Printf("Unexpected error, indexed value has type list but cannot be converted to list, please report this bug\n")
        return
      }

      index, err := strconv.Atoi(p)

      if err != nil {
        fmt.Printf("Property %s not found, index is not an integer\n", indexedString)
        return
      }

      if index < 0 || index >= len(list) {
        fmt.Printf("Property %s not found, index out of bounds\n", indexedString)
        return
      }

      indexedValue = list[index]
    }
  }

  if indexedValue.GetType() == parser.PARSER_VALUE_TYPE_STRING {
    cast, ok := indexedValue.(*parser.ParserValueString)

    if !ok {
      fmt.Printf("Unexpected error, indexed value has type string but cannot be cast to string, please report this bug\n")
      return
    }

    fmt.Println(cast.Value)
  } else {
    fmt.Println(indexedValue.ValueToString(2))
  }
}
