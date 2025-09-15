(integer_literal) @number
(string_literal) @string
(boolean) @boolean
(identifier) @variable
(call_expression) @function.call

"fnc" @keyword.function
"return" @keyword.return
"if" @keyword.conditional
"else" @keyword.conditional
"var" @keyword

[
  ","
  ";"
  ] @punctuation.delimiter

[
  "{"
  "}"
  "("
  ")"
  "["
  "]"
  ] @punctuation.bracket

[
  "!"
  "-"
  "=="
  "!="
  "<"
  ">"
  "+"
  "-"
  "*"
  "/"
  ] @operator