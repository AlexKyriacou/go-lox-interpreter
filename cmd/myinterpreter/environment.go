package main

type Envionment struct {
	values map[string]interface{}
}

func (e *Envionment) define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Envionment) get(name Token) (interface{}, error) {
	value, exists := e.values[name.lexeme]
	if exists {
		return value, nil
	}

	return nil, &RuntimeError{name, "Undefined variable '" + name.lexeme + "'."}
}

func (e *Envionment) assign(name Token, value interface{}) error {
	_, exists := e.values[name.lexeme]
	if exists {
		e.values[name.lexeme] = value
		return nil
	}
	return &RuntimeError{name, "Undefined variable '" + name.lexeme + "'."}
}
