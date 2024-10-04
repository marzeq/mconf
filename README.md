# mconf

## why?

this was originally made as a recreational programming exercise and a way to learn more about tokenisation and parsing

as i was working on it, it grew on me that i actually prefer the syntax of it over any other 

## getting the binary

go onto the [releases page](https://github.com/marzeq/mconf/releases) and download the binary for your platform

## editor support

for now, there is only a tree-sitter parser for mconf, which you can find [here](https://github.com/marzeq/tree-sitter-mconf). if you know how to make a syntax highlighter for your editor, you can use this to make one.

for neovim, i suggest you look at my [dotfiles](https://github.com/marzeq/dotfiles) to see how i setup a custom tree-sitter parser

## building & running

this repo comes with a justfile, so you can use [just](https://github.com/casey/just) to run the commands

```sh
just run        # runs the project, equivalent to `go run .`
just build      # builds the project in build/current-arch, equivalent to `go build .`
just build-all  # builds the project for every combination of: windows, darwin (macos) and linux & amd64 and arm64
```

if you don't want to use just, you can always use the go build and go run commands directly

## command line usage

```
Usage:
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
  cat config.mconf | %s - -- property1 property2
```

## spec

mconf fully suppports unicode, so a letter means any unicode latin letter and not just ascii letters, and string values can contain any unicode character

note: if you just one want a one-file example, look at the [examples/example.mconf](./examples/example.mconf) file

mconf is a superset of JSON, so any valid JSON file (base spec) is also a valid mconf file

### comments

```mconf
# this is a comment
a = 1 # this is a comment as well
```

### keys

```mconf
key = "value"
"strings as keys" = "are allowed"
óóóó_unicode = true
test: 1 # colon is also valid for JSON compatibility reasons

23abc = false # illegal, keys must start with a letter or underscore
```

if a key is defined many times, the last one will shadow the previous ones

### string values

```mconf
a_str = "bar"
multiline_str = "123
456"
escapes = "\"escaped quotes\""
unicode = "😊"
```

#### formatted strings

```mconf
user_and_port = "${USER}:${PORT}"
```

keep in mind, the {} **are required**

if the value substituted is not a string, it will be converted to one

### numerical values

```mconf
# integer value
an_int = 123
# signed integer value
a_uint = -123
# float value
a_float = 123.456
# fancy floats
fancy_float = .5

# hexadecimal value
hex = 0x123
# binary value
bin = 0b1010

# scientific notation
sci = 1.23e3
scineg = 1.23e-3
```

### boolean values

```mconf
a_bool = true
also_a_bool = false

yes_are_bools_too = yes
and_nos_as_well = no
and_on = on
and_off = off
```

### null values

only really added for JSON compatibility

```mconf
null_value = null
```

### list values

```mconf
list = [1, 2, 3, "abc", true, false]
two_dimensional_list = [
  [1, 2, 3],
  [4, 5, 6],
  [7, 8, 9]
]
```

commas in lists are required

### object values

```mconf
object = {
  foo = "bar"
  bar = 123
  baz = false
}

nested_object_and_list = {
  foo = {
    bar = "baz"
  }
  list = [1, 2, 3]
}
```

commas in objects are optional, but you can use them if you want

```mconf
object = {
  foo = "bar",
  bar = 123,
  baz = false,
}
```

#### top level objects

this is a a neat way to organise your file, where if you put an object at the top level, it's keys will be merged with the top level keys

you can think of this as a way to split the file into sections of multiple related keys

```mconf
{
  foo = 123
  bar = 123
}

{
  baz = 123
}
```

is equivalent to

```mconf
foo = 123
bar = 123
baz = 123
```

if a value is defined many times, the last one will shadow the previous ones, just like if they were all defined at the top level

### constants

```mconf
$some_constant = 123 
abc = $some_constant
```

if a constant is redefined, the previous references will not be affected, but the following ones will be

```mconf
$a = 123
foo = $a
$a = 456
bar = $a
```

will result in

```mconf
foo = 123
bar = 456
```

#### environment variables

enviorment variables are automatically loaded as any other constant

```mconf
user = $USER
```

#### default values

if a constant is not defined, you can put a `?` after it and then another constant or value that will be used as the default value

```mconf
$default_user = "some_user"
user = $USER?$default_user

# OR

user = $USER?"some_user"
```

### import

files can import other files, and the imported file will be parsed and merged with the current file (constants are shared between the files as well)

```mconf
@import "other_file.mconf"
```

in the case of an import cycle, the file that is second in the chain will only have access to the properties of the first file that were defined before the import

#### specific import

you can specify exactly what you want to import from a file

`a.mconf`:
```mconf
@import { foo, $bar, baz.bar } "b.mconf"
a = $bar
```

`b.mconf`:
```mconf
foo = 123
$bar = 456
baz = {
  bar = 789
}
```

will result in:
```mconf
foo = 123
bar = 789
a = 456
```

## todo:

- [x] support for formatted strings
- [x] merge current env vars with the constants
- [x] allow for specifying what exactly to import from a file
- [x] allow specyfing default values for constants/env vars if they are not set
- [x] hexadecimal and binary numbers
- [x] add a --json flag to convert mconf to json

## license

[do whatever with this, i don't care](./LICENSE)
