package main

type Callable interface {
	call(interpreter *Interpreter, arguments []interface{}) (interface{}, error)
	arity() int
}
