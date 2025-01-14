package main

import "time"

type clock struct{}

func (c clock) arity() int {
	return 0
}

func (c clock) call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	return float64(time.Now().Unix()), nil
}

func (c clock) String() string {
	return "<native fn>"
}
