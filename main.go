package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/marzeq/mconf/parser"
	"github.com/marzeq/mconf/tokeniser"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ParseFromString(s string, rootDir string, rootFile string, relativeDir string) (map[string]parser.ParserValue, map[string]parser.ParserValue, error) {
	t := tokeniser.NewTokeniser(s, rootFile, relativeDir)
	tokens, err := t.Tokenise()
	if err != nil {
		return nil, nil, err
	}

	p := parser.NewParser(tokens, rootDir, rootFile, relativeDir)
	parsed, err := p.Parse()
	if err != nil {
		return nil, nil, err
	}
	constants := p.GetConstants()

	return parsed, constants, nil
}

func ParseFromFile(filename string) (map[string]parser.ParserValue, map[string]parser.ParserValue, error) {
	f, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}

	s := string(f)

	relativeDir := filepath.Dir(filename)

	fileDir, fdirErr := filepath.Abs(relativeDir)
	if fdirErr != nil {
		return nil, nil, fdirErr
	}

	baseFile := filepath.Base(filename)

	return ParseFromString(s, fileDir, baseFile, relativeDir)
}

func ParseFromStdin() (map[string]parser.ParserValue, map[string]parser.ParserValue, error) {
	b, err := readStdin()
	if err != nil {
		return nil, nil, err
	}

	s := string(b)

	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return nil, nil, cwdErr
	}

	absCwd, absCwdErr := filepath.Abs(cwd)
	if absCwdErr != nil {
		return nil, nil, absCwdErr
	}

	return ParseFromString(s, absCwd, "", cwd)
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
	ShowConstants     bool
	EnvFile           string
}

func usage(progname string) string {
	return fmt.Sprintf(`Usage:
  %s <filename> [-- property1 property2 ...]

Arguments:
  <filename>                    Path to the configuration file. Use '-' to read from stdin.
  [-- property1 property2 ...]  List of properties to access. Multiple properties are used to access nested objects or lists. If no properties are provided, the global object is printed. '--' is simply there for readability.

Options:
  -h, --help        Show this message
  -v, --version     Show version
  -j, --json        Output as JSON (in a compact format, prettyfication is up to the user)
  -d, --dotenv      Load .env file in current directory
	--envfile <file>  Load specified enviorment variables file
  -c, --constants   Show constants (only displayed when no properties are provided)

Examples:
  %s config.mconf -- property1 property2
  cat config.mconf | %s - -- property1 property2`, progname, progname, progname)
}

func version() string {
	return "mconf (version 0.1.1)"
}

func parseOptions() (options, string, uint) {
	opts := options{}
	opts.ToJson = false
	opts.ShowConstants = false
	opts.EnvFile = ""

	args := os.Args[1:]

	binname := filepath.Base(os.Args[0])

	if len(args) == 0 {
		return opts, usage(binname), 1
	}

	i := 0
	providedFilename := false
	for {
		if i >= len(args) {
			break
		}

		arg := args[i]

		if arg[0] == '-' {
			if arg == "-" {
				if providedFilename {
					return opts, "Provided multiple filenames, only one is allowed", 1
				}

				opts.Filename = "-"
				providedFilename = true
				i++
				continue
			}

			if arg[1] == '-' {
				if arg == "--" {
					if !providedFilename {
						return opts, "No filename provided", 0
					}

					opts.AcessedProperties = args[i+1:]
					break
				} else if arg == "--help" {
					return opts, usage(binname), 0
				} else if arg == "--version" {
					return opts, version(), 0
				} else if arg == "--json" {
					opts.ToJson = true
				} else if arg == "--constants" {
					opts.ShowConstants = true
				} else if arg == "--dotenv" {
					opts.EnvFile = ".env"
				} else if arg == "--envfile" {
					if i+1 >= len(args) {
						return opts, "No argument provided for --envfile", 1
					} else {
						opts.EnvFile = args[i+1]
						i++
					}
				}
			} else {
				for _, c := range arg[1:] {
					switch c {
					case 'h':
						return opts, usage(binname), 0
					case 'v':
						return opts, version(), 0
					case 'j':
						opts.ToJson = true
					case 'c':
						opts.ShowConstants = true
					case 'd':
						opts.EnvFile = ".env"
					}
				}
			}
		} else {
			if providedFilename {
				return opts, "Provided multiple filenames, only one is allowed", 1
			}

			opts.Filename = arg
			providedFilename = true
		}

		i++
	}

	if !providedFilename {
		return opts, "No filename provided", 1
	}

	return opts, "", 0
}

func main() {
	opts, usage, exitcode := parseOptions()

	if usage != "" {
		fmt.Println(usage)
		os.Exit(int(exitcode))
	}

	var globalObj map[string]parser.ParserValue
	var constants map[string]parser.ParserValue
	var parsingErr error

	if opts.EnvFile != "" {
		err := os.Setenv("MCONF_ENV_FILE", opts.EnvFile)
		if err != nil {
			fmt.Println("Error setting environment variable MCONF_ENV_FILE")
			os.Exit(1)
		}

		envFile, err := os.ReadFile(opts.EnvFile)
		if err != nil {
			fmt.Printf("Error reading environment file %s\n", opts.EnvFile)
			os.Exit(1)
		}

		envFileStr := string(envFile)
		envLines := strings.Split(envFileStr, "\n")

		for _, line := range envLines {
			if line == "" {
				continue
			}

			parts := strings.SplitN(line, "=", 2)

			if len(parts) != 2 {
				fmt.Printf("Error parsing environment file %s\n", opts.EnvFile)
				os.Exit(1)
			}

			if parts[0] == "" {
				fmt.Printf("Error setting environment variable (%s)\n", line)
				os.Exit(1)
			}

			if parts[1][0] == '"' && parts[1][len(parts[1])-1] == '"' {
				parts[1] = parts[1][1 : len(parts[1])-1]
				parts[1] = strings.ReplaceAll(parts[1], "\\n", "\n")
				parts[1] = strings.ReplaceAll(parts[1], "\\r", "\r")
				parts[1] = strings.ReplaceAll(parts[1], "\\t", "\t")
				parts[1] = strings.ReplaceAll(parts[1], "\\\"", "\"")
				parts[1] = strings.ReplaceAll(parts[1], "\\\\", "\\")
			}

			err := os.Setenv(parts[0], parts[1])
			if err != nil {
				fmt.Printf("Error setting environment variable (%s)\n", line)
				os.Exit(1)
			}
		}
	}

	if opts.Filename == "-" {
		globalObj, constants, parsingErr = ParseFromStdin()
	} else {
		globalObj, constants, parsingErr = ParseFromFile(opts.Filename)
	}

	check(parsingErr)

	var indexedValue parser.ParserValue = &parser.ParserValueObject{Value: globalObj}

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
		if opts.ShowConstants {
			fmt.Printf("Displaying constants is not supported when outputting as JSON\n")
			os.Exit(1)
		}
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

		if len(opts.AcessedProperties) == 0 && opts.ShowConstants {
			for k, v := range constants {
				fmt.Printf("$%s = %s\n", k, v.ValueToString(2))
			}
		}
	}
}
