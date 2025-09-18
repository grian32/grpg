package grpgscript_lsp

import "grpgscript/object"

func BuildDefinitions() map[string]BuiltinDefinition {
	defMap := make(map[string]BuiltinDefinition)

	defMap["spawnNpc"] = NewBuiltinDefinition(
		"spawnNpc",
		[]string{"npcId", "x", "y"},
		[]TypeTag{INT, INT, INT},
		NULL,
	)

	defMap["onInteract"] = NewBuiltinDefinition(
		"onInteract",
		[]string{"objId", "callback"},
		[]TypeTag{INT, FUNCTION},
		NULL,
	)

	defMap["onTalkNpc"] = NewBuiltinDefinition(
		"onTalkNpc",
		[]string{"npcId", "callback"},
		[]TypeTag{INT, FUNCTION},
		NULL,
	)

	defMap["getObjState"] = NewBuiltinDefinition(
		"getObjState",
		[]string{},
		[]TypeTag{},
		INT,
	)

	defMap["setObjState"] = NewBuiltinDefinition(
		"setObjState",
		[]string{"newState"},
		[]TypeTag{INT},
		NULL,
	)

	defMap["playerInvAdd"] = NewBuiltinDefinition(
		"playerInvAdd",
		[]string{"itemId"},
		[]TypeTag{INT},
		NULL,
	)

	defMap["timer"] = NewBuiltinDefinition(
		"timer",
		[]string{"tickCount", "callback"},
		[]TypeTag{INT, FUNCTION},
		NULL,
	)

	defMap["talkPlayer"] = NewBuiltinDefinition(
		"talkPlayer",
		[]string{"talk"},
		[]TypeTag{STRING},
		NULL,
	)

	defMap["talkNpc"] = NewBuiltinDefinition(
		"talkNpc",
		[]string{"talk"},
		[]TypeTag{STRING},
		NULL,
	)

	defMap["clearDialogueQueue"] = NewBuiltinDefinition(
		"clearDialogueQueue",
		[]string{},
		[]TypeTag{},
		NULL,
	)

	defMap["startDialogue"] = NewBuiltinDefinition(
		"startDialogue",
		[]string{},
		[]TypeTag{},
		NULL,
	)

	return defMap
}

func MockBuiltins(env *object.Environment, definitions map[string]BuiltinDefinition) {
	env.Set("spawnNpc", MockBuiltin(definitions["spawnNpc"], nil))

	env.Set("onInteract", MockBuiltin(definitions["onInteract"], []NamedBuiltin{
		NewNamedBuiltin("getObjState", MockBuiltin(definitions["getObjState"], nil)),
		NewNamedBuiltin("setObjState", MockBuiltin(definitions["setObjState"], nil)),
		NewNamedBuiltin("playerInvAdd", MockBuiltin(definitions["playerInvAdd"], nil)),
		NewNamedBuiltin("timer", MockBuiltin(definitions["timer"], nil)),
	}))

	env.Set("onTalkNpc", MockBuiltin(definitions["onTalkNpc"], []NamedBuiltin{
		NewNamedBuiltin("talkPlayer", MockBuiltin(definitions["talkPlayer"], nil)),
		NewNamedBuiltin("talkNpc", MockBuiltin(definitions["talkNpc"], nil)),
		NewNamedBuiltin("clearDialogueQueue", MockBuiltin(definitions["clearDialogueQueue"], nil)),
		NewNamedBuiltin("startDialogue", MockBuiltin(definitions["startDialogue"], nil)),
	}))
}
