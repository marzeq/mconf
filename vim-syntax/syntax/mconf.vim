" Syntax highlighting for .mconf files

" Define the comment syntax
syn match mconfComment "#.*$"

" Define the constant syntax
syn match mconfConstant "\$[a-zA-Z_]\w*"

" Define the boolean values
syn match mconfBoolean "\<\(true\|false\)\>"

" Define the numbers (integer, float, negative)
syn match mconfNumber "\<\(-\=\d\+\|-\=\d*\.\d\+\|-\=\.\d\+\)\>"

" Define the strings and handle escape sequences
syn region mconfString start=+"+ skip=+\\."+ end=+"+ contains=mconfEscape
syn match mconfEscape "\\."

" Define the key-value pairs (key without the equal sign)
syn match mconfKey "[a-zA-Z_]\w*\s*\ze="

" Define the object and list syntax
syn match mconfObject "[{}]"
syn match mconfList "[\[\]]"

" Define the highlighter groups
hi def link mconfComment Comment
hi def link mconfConstant Constant
hi def link mconfBoolean Boolean
hi def link mconfNumber Number
hi def link mconfString String
hi def link mconfEscape SpecialChar
hi def link mconfKey Identifier
hi def link mconfObject Structure
hi def link mconfList Structure

