package main

type LoxFunction struct {
	declaration Function
}

func (l *LoxFunction) call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	var environment = NewEnvironment(interpreter.globals)
	for i, param := range l.declaration.params {
		environment.define(param.lexeme, arguments[i])
	}

	err := interpreter.executeBlock(l.declaration.body, environment)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (l *LoxFunction) arity() int {
	return len(l.declaration.params)
}

func (l *LoxFunction) String() string {
	return "<fn " + l.declaration.name.lexeme + ">"
}
