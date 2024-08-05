" Indentation settings for .mconf files

" Load the default indentation rules for now
setlocal indentexpr=

" Set the basic indent and shift width
setlocal shiftwidth=2
setlocal tabstop=2
setlocal expandtab

" Custom indentation function
function! s:indent_mconf() abort
  let lnum = v:lnum

  " Get the current line and the previous line
  let line = getline(lnum)
  let prevline = getline(lnum - 1)

  " Remove leading and trailing whitespace
  let line = substitute(line, '^\s*', '', '')
  let prevline = substitute(prevline, '^\s*', '', '')

  " If the current line is inside an object or list, increase indent
  if line =~ '^\s*[a-zA-Z_]\w*\s*=' && prevline =~ '{\|[\['
    return indent(lnum - 1) + &shiftwidth
  endif

  " If the previous line ends with an opening brace or bracket, increase indent
  if prevline =~ '{\|[\['
    return indent(lnum - 1) + &shiftwidth
  endif

  " If the current line closes an object or list, decrease indent
  if line =~ '^\s*[\}\]]'
    return indent(lnum - 1) - &shiftwidth
  endif

  " Otherwise, use the same indent as the previous line
  return indent(lnum - 1)
endfunction

" Set the custom indentation function for .mconf files
setlocal indentexpr=s:indent_mconf()
