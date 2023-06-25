package compiler

var Universe *Scope = NewScope(nil)

var builtinObjects = []*Object{
	{Name: "println", MangledName: "@tiny_go_builtin_println"},
	{Name: "exit", MangledName: "@tiny_go_builtin_exit"},
}

func init() {
	for _, obj := range builtinObjects {
		Universe.Insert(obj)
	}
}
