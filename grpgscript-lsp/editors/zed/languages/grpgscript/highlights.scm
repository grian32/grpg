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

((identifier) @function.builtin
  (#match? @function.builtin "^(println|len|push|unshift|concat|onInteract|onTalkNpc|spawnNpc|getObjState|setObjState|playerInvAdd|playerAddXp|timer|talkPlayer|talkNpc|clearDialogueQueue|startDialogue)$"))
