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
# float value (tokenised the same way but may be parsed as a different value depending on the implementation)
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
  baz = false
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

### constants

```mconf
$some_constant = 123 
abc = $some_constant
```

### includes

files can include other files, and the included file will be parsed and merged with the current file

```mconf
@include "other_file.mconf"
```

inclusion is in an early stage in development, so you cannot access the included file's keys or constants in the current file, it's final parsed global object will be merged with the current file's global object

## todo:

- [ ] support for formatted strings that allow for inserting constants into the string like `foo = f"bar ${baz}"` (where baz is a previously defined constant)
- [ ] merge current env vars with the constants
- [ ] support for specifying which keys to include with syntax like `@include { key1, key2 } "other_file.mconf"`
- [ ] when using `@include`, merge the included file's constants with the current file's constants, and as a logical result, also allow for specyfing which constants to include in the syntax above
- [ ] allow the use of the `@include` directive as a value with syntax like `foo = @include bar "other_file.mconf"`
- [ ] a `@template` directive that allows for defining a template that can be used in the file, like `@template !my_template(foo) { foo = $foo }` and then calling it like `foo = !my_template(123)`

## quirks of this particular parser

### numbers

- if a number does not contain either a `.` or a `-`, it will be parsed as a `uint64`, and if it overflows, it will be reset back to the max `uint64` value
- if a number does not contain a `.` but contains a `-`, it will be parsed as a `int64`, and if it overflows or underflows, it will be reset back to the max `int64` value
- otherwise, the number will be parsed as a float, and no bounds checking will be done

## license

[do whatever with this, i don't care](./LICENSE)
