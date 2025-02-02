package main

import "errors"

type LoxFunction struct {
	declaration   Function
	closure       *Envionment
	isInitializer bool
}

func (l *LoxFunction) bind(instance *LoxInstance) LoxFunction {
	var environment *Envionment = NewEnvironment(l.closure)
	environment.define("this", instance)
	return LoxFunction{l.declaration, environment, l.isInitializer}
}

func (l LoxFunction) call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	var environment = NewEnvironment(l.closure)
	for i, param := range l.declaration.params {
		environment.define(param.lexeme, arguments[i])
	}

	err := interpreter.executeBlock(l.declaration.body, environment)
	if err != nil {
		if errors.Is(err, &ReturnException{}) {
			if l.isInitializer {
				return l.closure.getAt(0, "this"), nil
			}
			// if a return exception is caught we want to return its value
			return err.(*ReturnException).value, nil
		}
		return nil, err
	}

	if l.isInitializer {
		return l.closure.getAt(0, "this"), nil
	}
	return nil, nil
}

func (l LoxFunction) arity() int {
	return len(l.declaration.params)
}

func (l LoxFunction) String() string {
	return "<fn " + l.declaration.name.lexeme + ">"
}
