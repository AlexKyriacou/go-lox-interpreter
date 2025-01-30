package main

import "fmt"

type ReturnException struct {
	value interface{}
}

func (r *ReturnException) Error() string {
	return fmt.Sprint(r.value)
}

func (r ReturnException) Is(target error) bool {
	_, ok := target.(*ReturnException)
	return ok
}
