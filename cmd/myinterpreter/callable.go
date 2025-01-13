package main

type Callable interface {
	call(interpreter *Interpreter, arguements []interface{}) (interface{}, error)
	arity() int
}
