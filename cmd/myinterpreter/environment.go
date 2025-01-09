package main

type Envionment struct {
	values map[string]interface{}
}

func (e *Envionment) define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Envionment) get(name Token) (interface{}, error) {
	val, exists := e.values[name.lexeme]
	if exists {
		return val, nil
	}

	return nil, &RuntimeError{name, "Undefined variable '" + name.lexeme + "'."}
}
