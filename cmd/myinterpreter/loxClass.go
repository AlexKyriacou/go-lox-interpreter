package main

type LoxClass struct {
	name string
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
