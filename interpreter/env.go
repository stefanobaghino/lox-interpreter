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

func (e *Env) Get(name string) interface{} {
	if value, ok := e.values[name]; ok {
		return value
	} else if e.parent != nil {
		return e.parent.Get(name)
	} else {
		panic(&RuntimeError{message: "variable not defined"})
	}
}
