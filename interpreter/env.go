package interpreter

type Env struct {
	parent *Env
	values map[string]interface{}
}

func NewGlobalEnv() *Env {
	return &Env{values: make(map[string]interface{})}
}

func NewEnv(parent *Env) *Env {
	return &Env{parent: parent, values: make(map[string]interface{})}
}

func (e *Env) Define(name string, initializer func() interface{}) {
	if _, ok := e.values[name]; ok {
		panic(&RuntimeError{message: "variable already declared"})
	}
	e.values[name] = initializer()
}

func (e *Env) Assign(name string, initializer func() interface{}) interface{} {
	if _, ok := e.values[name]; ok {
		value := initializer()
		e.values[name] = value
		return value
	} else if e.parent != nil {
		return e.parent.Assign(name, initializer)
	} else {
		panic(&RuntimeError{message: "variable not declared"})
	}
}

func (e *Env) AssignAt(distance int, name string, initializer func() interface{}) interface{} {
	return e.ancestor(distance).Assign(name, initializer)
}

func (e *Env) Get(name string) interface{} {
	if value, ok := e.values[name]; ok {
		return value
	} else if e.parent != nil {
		return e.parent.Get(name)
	} else {
		panic(&RuntimeError{message: "variable not defined"})
	}
}

func (e *Env) GetAt(distance int, name string) interface{} {
	return e.ancestor(distance).Get(name)
}

func (e *Env) ancestor(distance int) *Env {
	env := e
	for i := 0; i < distance; i++ {
		env = env.parent
	}
	return env
}
