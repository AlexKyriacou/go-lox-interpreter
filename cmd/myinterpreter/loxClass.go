package main

type LoxClass struct {
	name    string
	methods map[string]LoxFunction
}

func (l LoxClass) String() string {
	return l.name
}

func (l LoxClass) call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	var instance *LoxInstance = NewLoxInstance(l)
	return instance, nil
}

func (l LoxClass) arity() int {
	return 0
}

func (l LoxClass) findMethod(name string) (*LoxFunction, bool) {
	if value, prs := l.methods[name]; prs {
		return &value, true
	}
	return nil, false
}
