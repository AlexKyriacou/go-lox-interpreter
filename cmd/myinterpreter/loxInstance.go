package main

type LoxInstance struct {
	class  LoxClass
	fields map[string]interface{}
}

func NewLoxInstance(class LoxClass) *LoxInstance {
	return &LoxInstance{class, make(map[string]interface{})}
}

func (l *LoxInstance) get(name Token) (interface{}, error) {
	field, prs := l.fields[name.lexeme]
	if prs {
		return field, nil
	}

	method, prs := l.class.findMethod(name.lexeme)
	if prs {
		return method, nil
	}

	return nil, &RuntimeError{name, "Undefiend property '" + name.lexeme + "'."}
}

func (l *LoxInstance) set(name Token, value interface{}) {
	l.fields[name.lexeme] = value
}

func (l LoxInstance) String() string {
	return l.class.name + " instance"
}
