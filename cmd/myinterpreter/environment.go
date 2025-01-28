package main

type Envionment struct {
	values    map[string]interface{}
	enclosing *Envionment
}

func NewEnvironment(enclosing *Envionment) *Envionment {
	return &Envionment{values: make(map[string]interface{}), enclosing: enclosing}
}

func (e *Envionment) define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Envionment) get(name Token) (interface{}, error) {
	value, exists := e.values[name.lexeme]
	if exists {
		return value, nil
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	return nil, &RuntimeError{name, "Undefined variable '" + name.lexeme + "'."}
}

func (e *Envionment) getAt(distance int, name string) interface{} {
	return e.ancestor(distance).values[name]
}

func (e *Envionment) assignAt(distance int, name Token, value interface{}) {
	e.ancestor(distance).values[name.lexeme] = value
}

func (e *Envionment) ancestor(distance int) *Envionment {
	var environment *Envionment = e
	for i := 0; i < distance; i++ {
		environment = environment.enclosing
	}
	return environment
}

func (e *Envionment) assign(name Token, value interface{}) error {
	_, exists := e.values[name.lexeme]
	if exists {
		e.values[name.lexeme] = value
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.assign(name, value)
	}
	return &RuntimeError{name, "Undefined variable '" + name.lexeme + "'."}
}
