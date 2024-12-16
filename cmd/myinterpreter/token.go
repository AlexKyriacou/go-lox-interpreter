package main

import "fmt"

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   interface{}
	line      int
}

func (t Token) toString() string {
	var literalStr string
	if t.literal == nil {
		literalStr = "null"
	} else {
		literalStr = fmt.Sprintf("%v", t.literal)
	}
	return fmt.Sprintf("%v %s %s", t.tokenType, t.lexeme, literalStr)
}