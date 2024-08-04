# mconf

this is my own configuration language, made mostly for fun and to learn more about tokenisation and parsing

for this reason, i wouldn't really use this implementation for anything serious, but you're welcome to use it if you want

the syntax and design is a bit better, designing a config language is obviously easier than implementing it, and i think i've done a good job with it, but there may be a couple things you might find weird/not like

if you have any suggestions, feedback or questions, feel free to contact me in any way you want (open an issue here, message me on any platform, etc.)

## syntax

note that mconf fully suppports unicode, so a letter means any unicode letter and not just ascii letters, and string values can contain any unicode character

### comments

```conf
# this is a comment
```

### keys

```conf
key = "value"
"strings as keys" = "are allowed"
23abc = false # illegal, keys must start with a letter
```

### string values

```conf
foo = "bar"
baz = 'abc'
def = "123
456"
ghi = "\"escaped quotes\""
```

### numerical values

```conf
# integer value
foo = 123
# float value (tokenised the same but may be parsed differently depending on the target language)
bar = 123.456
```

### boolean values

```conf
deez = true
nuts = false
```

### lists

```conf
list = [1, 2, 3, "abc", 'def', true, false]
two_dimensional_list = [[1, 2, 3], [4, 5, 6], [7, 8, 9]]
```

### objects/dictionaries/maps/(whatever you want to call them)

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

## license

do whatever with this, i don't care
