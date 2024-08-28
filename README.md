# mconf

this is my own configuration language, made mostly for fun and to learn more about tokenisation and parsing

if you have any suggestions, feedback or questions, feel free to contact me in any way you want (open an issue here, message me on any platform, etc.)

## building & running

this repo comes with a justfile, so you can use [just](https://github.com/casey/just) to run the commands

```sh
just run # runs the project, equivalent to `go run .`
just build # builds the project in build/current-arch, equivalent to `go build .`
just build-all # builds the project for every combination of: windows, darwin (macos) and linux & amd64 and arm64
```

if you don't want to use just, you can always use the go build and go run commands directly

## spec

note that mconf fully suppports unicode, so a letter means any unicode latin letter and not just ascii letters, and string values can contain any unicode character

### comments

```mconf
# this is a comment
```

### keys

```mconf
key = "value"
"strings as keys" = "are allowed"
óóóó_unicode = true
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
```

### boolean values

```mconf
a_bool = true
also_a_bool = false
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

- [ ] support for formatted strings that allow for inserting constants into the string like `foo = f"bar ${baz}"` (where baz is a previously defined constant)
- [x] merge current env vars with the constants
- [x] allow for specifying what exactly to import from a file
- [ ] a `@template` directive that allows for defining a template that can be used in the file, like `@template !my_template(foo) { foo = $foo }` and then calling it like `foo = !my_template(123)`
- [ ] allow specyfing default values for enviorment variables if they are not set
- [ ] hexadecimal and binary numbers

## license

[do whatever with this, i don't care](./LICENSE)
