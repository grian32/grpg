package object

func NewEnclosedEnvinronment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	for s, _ := range outer.Names {
		env.Names[s] = struct{}{}
	}
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, Names: make(map[string]struct{})}
}

type Environment struct {
	store map[string]Object
	outer *Environment
	// Names is a list of all items contained in this env and its outer, used primarly for lsp lookup
	Names map[string]struct{}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	e.Names[name] = struct{}{}
	return val
}
