package main

type LoxClass struct {
	name       string
	superclass *LoxClass
	methods    map[string]LoxFunction
}

func (l LoxClass) String() string {
	return l.name
}

func (l LoxClass) call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	var instance *LoxInstance = NewLoxInstance(l)
	initializer, prs := l.findMethod("init")
	if prs {
		initializer.bind(instance).call(interpreter, arguments)
	}
	return instance, nil
}

func (l LoxClass) arity() int {
	initializer, prs := l.findMethod("init")
	if !prs {
		return 0
	}
	return initializer.arity()
}

func (l LoxClass) findMethod(name string) (*LoxFunction, bool) {
	if value, prs := l.methods[name]; prs {
		return &value, true
	}

	// if we don't find the method on the instance we recurse up the
	// superclass chain and look there
	if l.superclass != nil {
		return l.superclass.findMethod(name)
	}

	return nil, false
}
