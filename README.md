# mconf

this is my own configuration language, made mostly for fun and to learn more about tokenisation and parsing

if you have any suggestions, feedback or questions, feel free to contact me in any way you want (open an issue here, message me on any platform, etc.)

## spec

note that mconf fully suppports unicode, so a letter means any unicode latin letter and not just ascii letters, and string values can contain any unicode character

### comments

```conf
# this is a comment
```

### keys

```conf
key = "value"
"strings as keys" = "are allowed"
23abc = false # illegal, keys must start with a letter or underscore
óóóó_unicode = true
```

### string values

```conf
a_str = "bar"
multiline_str = "123
456"
escapes = "\"escaped quotes\""
unicode = "😊"
```

### numerical values

```conf
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

```conf
a_bool = true
also_a_bool = false
```

### list values

```conf
list = [1, 2, 3, "abc", true, false]
two_dimensional_list = [[1, 2, 3], [4, 5, 6], [7, 8, 9]]
```

commas in lists are required

### object values

```conf
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

```conf
object = {
  foo = "bar",
  bar = 123,
  baz = false
}
```

#### top level objects

this is a a neat way to organise your file, where if you put an object at the top level, it's keys will be merged with the top level keys

you can think of this as a way to split the file into sections of multiple related keys

```conf
{
  foo = 123
  bar = 123
}

{
  baz = 123
}
```

is equivalent to

```conf
foo = 123
bar = 123
baz = 123
```

### constants

```conf
$some_constant = 123 
abc = $some_constant
```

## quirks of this particular parser

### numbers

- if a number does not contain either a `.` or a `-`, it will be parsed as a `uint64`, and if it overflows, it will be reset back to the max `uint64` value
- if a number does not contain a `.` but contains a `-`, it will be parsed as a `int64`, and if it overflows or underflows, it will be reset back to the max `int64` value
- otherwise, the number will be parsed as a float, and no bounds checking will be done

## license

[do whatever with this, i don't care](./LICENSE)
