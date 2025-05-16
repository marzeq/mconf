package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/marzeq/mconf"
	"github.com/marzeq/mconf/mconf_values"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type options struct {
	Filename          string
	AcessedProperties []string
	ToJson            bool
	ShowConstants     bool
	EnvFile           string
}

func usage() string {
	progname := filepath.Base(os.Args[0])

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
  %s config.mconf -- property1 1 
  cat config.mconf | %s - -- property1 1`, progname, progname, progname)
}

func version() string {
	return fmt.Sprintf("%s version %s", mconf.PROGNAME, mconf.VERSION)
}

func parseOptions() (options, string, uint) {
	opts := options{}
	opts.ToJson = false
	opts.ShowConstants = false
	opts.EnvFile = ""

	args := os.Args[1:]

	if len(args) == 0 {
		return opts, usage(), 1
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
					return opts, usage(), 0
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
						return opts, usage(), 0
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

	var globalObj map[string]mconf_values.MconfValue
	var constants map[string]mconf_values.MconfValue
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

			if parts[1][0] == '"' && parts[1][len(parts[1])-1] == '"' ||
				parts[1][0] == '\'' && parts[1][len(parts[1])-1] == '\'' {

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
		globalObj, constants, parsingErr = mconf.ParseFromStdin()
	} else {
		globalObj, constants, parsingErr = mconf.ParseFromFile(opts.Filename)
		if parsingErr != nil {
			parsingErr = fmt.Errorf("%s - Error reading file,%s", opts.Filename, strings.Split(parsingErr.Error(), ":")[1])
		}
	}

	check(parsingErr)

	var indexedValue mconf_values.MconfValue = &mconf_values.MconfObject{Value: globalObj}

	indexedString := ""

	for _, p := range opts.AcessedProperties {
		if indexedString == "" {
			indexedString = p
		} else {
			indexedString += "." + p
		}

		switch indexedValue.(type) {
		case *mconf_values.MconfObject:
			obj := indexedValue.(*mconf_values.MconfObject).Value

			next := obj[p]

			if next == nil {
				fmt.Printf("Property %s not found\n", indexedString)
				os.Exit(1)
			}

			indexedValue = next
		case *mconf_values.MconfList:
			list := indexedValue.(*mconf_values.MconfList).Value

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
		default:
			fmt.Printf("Property %s not found, indexed value is not an object or list\n", indexedString)
			os.Exit(1)
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

	switch indexedValue.(type) {
	case *mconf_values.MconfString:
		fmt.Println(indexedValue.(*mconf_values.MconfString).Value)
	default:
		fmt.Println(indexedValue.ValueToString(2))

		if len(opts.AcessedProperties) == 0 && opts.ShowConstants {
			for k, v := range constants {
				fmt.Printf("$%s = %s\n", k, v.ValueToString(2))
			}
		}
	}
}
