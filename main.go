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
	ToJson            bool
}

func usage(progname string) string {
	return fmt.Sprintf(`Usage:
  %s <filename> [-- property1 property2 ...]

Arguments:
  <filename>        Path to the configuration file. Use '-' to read from stdin.
  [-- property1 property2 ...]  List of properties to access. Multiple properties are used to access nested objects or lists. If no properties are provided, the global object is printed. -- is simply there for readability.

Options:
  -h, --help        Show this message
  -v, --version     Show version
  -j, --json        Output as JSON (in a compact format, prettyfication is up to the user)

Examples:
  %s config.mconf -- property1 property2
  cat config.mconf | %s - -- property1 property2`, progname, progname, progname)
}

func version() string {
	return "mconf (development version)"
}

func parseOptions() (options, string) {
	opts := options{}
	opts.ToJson = false

	args := os.Args[1:]

	binname := filepath.Base(os.Args[0])

	if len(args) == 0 {
		return opts, usage(binname)
	}

	i := 0
	providedFilename := false
	for {
		if i >= len(args) {
			break
		}

		arg := args[i]

		if arg[0] == '-' {
			if arg[1] == '-' {
				if arg == "--" {
					if !providedFilename {
						return opts, "No filename provided"
					}

					opts.AcessedProperties = args[i+1:]
					break
				} else if arg == "--help" {
					return opts, usage(binname)
				} else if arg == "--version" {
					return opts, version()
				} else if arg == "--json" {
					opts.ToJson = true
				}
			} else {
				for _, c := range arg[1:] {
					switch c {
					case 'h':
						return opts, usage(binname)
					case 'v':
						return opts, version()
					case 'j':
						opts.ToJson = true
					}
				}
			}
		} else {
			if providedFilename {
				return opts, "Provided multiple filenames, only one is allowed"
			}

			opts.Filename = arg
			providedFilename = true
		}

		i++
	}

	if !providedFilename {
		return opts, "No filename provided"
	}

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

	if opts.ToJson {
		fmt.Println(indexedValue.ToJSONString())
		return
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
