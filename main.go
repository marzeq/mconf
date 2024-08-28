package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func ParseFromString(s string, rootDir string, rootFile string) (map[string]parser.ParserValue, error) {
	t := tokeniser.NewTokeniser(s, rootFile)
	tokens, err := t.Tokenise()
	if err != nil {
		return nil, err
	}

	p := parser.NewParser(tokens, rootDir, rootFile)

	return p.Parse()
}

func ParseFromFile(filename string) (map[string]parser.ParserValue, error) {
	f, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	s := string(f)

	fileDir, fdirErr := filepath.Abs(filepath.Dir(filename))
	if fdirErr != nil {
		return nil, fdirErr
	}

	baseFile := filepath.Base(filename)

	return ParseFromString(s, fileDir, baseFile)
}

func ParseFromStdin() (map[string]parser.ParserValue, error) {
	b, err := readStdin()
	if err != nil {
		return nil, err
	}

	s := string(b)

	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return nil, cwdErr
	}

	absCwd, absCwdErr := filepath.Abs(cwd)
	if absCwdErr != nil {
		return nil, absCwdErr
	}

	return ParseFromString(s, absCwd, "")
}

func readStdin() ([]byte, error) {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}

	return b, nil
}

type options struct {
	Filename          string
	AcessedProperties []string
}

func usage(progname string) string {
	return fmt.Sprintf(`Usage:
  %s <filename> [properties ...]

Arguments:
  <filename>        Path to the configuration file. Use '-' to read from stdin.
  [properties ...]  List of properties to access. Multiple properties are used to access nested objects or lists. If no properties are provided, the global object is printed.

Options:
  -h, --help        Show this message
  -v, --version     Show version

Examples:
  %s config.mconf property1 property2
  cat config.mconf | %s - property1 property2`, progname, progname, progname)
}

func parseOptions() (options, string) {
	opts := options{}

	if len(os.Args) == 1 {
		return opts, usage(os.Args[0])
	}

	for _, arg := range os.Args[1:] {
		if arg[0] == '-' && arg != "-" {
			if arg == "-h" || arg == "--help" {
				return opts, usage(os.Args[0])
			} else if arg == "-v" || arg == "--version" {
				return opts, "mconf (development version)"
			} else {
				return opts, fmt.Sprintf("Unknown option: %s", arg)
			}
		}
	}

	opts.Filename = os.Args[1]
	opts.AcessedProperties = os.Args[2:]

	return opts, ""
}

func main() {
	opts, usage := parseOptions()

	if usage != "" {
		fmt.Println(usage)
		os.Exit(1)
	}

	var m map[string]parser.ParserValue
	var parsingErr error

	if opts.Filename == "-" {
		m, parsingErr = ParseFromStdin()
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
			os.Exit(1)
		}

		if indexedValue.GetType() == parser.PARSER_VALUE_TYPE_OBJECT {
			obj, err := indexedValue.GetObject()
			if err != nil {
				fmt.Printf("Unexpected error, indexed value has type object but cannot be converted to object, please report this bug\n")
				os.Exit(1)
			}

			next := obj[p]

			if next == nil {
				fmt.Printf("Property %s not found\n", indexedString)
				os.Exit(1)
			}

			indexedValue = next
		} else {
			list, err := indexedValue.GetList()
			if err != nil {
				fmt.Printf("Unexpected error, indexed value has type list but cannot be converted to list, please report this bug\n")
				os.Exit(1)
			}

			index, err := strconv.Atoi(p)
			if err != nil {
				fmt.Printf("Property %s not found, index is not an integer\n", indexedString)
				os.Exit(1)
			}

			if index < 0 || index >= len(list) {
				fmt.Printf("Property %s not found, index out of bounds\n", indexedString)
				os.Exit(1)
			}

			indexedValue = list[index]
		}
	}

	if indexedValue.GetType() == parser.PARSER_VALUE_TYPE_STRING {
		cast, ok := indexedValue.(*parser.ParserValueString)

		if !ok {
			fmt.Printf("Unexpected error, indexed value has type string but cannot be cast to string, please report this bug\n")
			os.Exit(1)
		}

		fmt.Println(cast.Value)
	} else {
		fmt.Println(indexedValue.ValueToString(2))
	}
}
