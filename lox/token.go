package main

import (
	"fmt"
)

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   interface{}
	line      int
}

func (t Token) String() string {
	return fmt.Sprintf("%v %s %s", t.tokenType, t.lexeme, t.formatLiteral())
}

func (t Token) formatLiteral() string {
	var literalStr string
	if t.literal == nil {
		literalStr = "null"
	} else if t.tokenType == NUMBER {
		// This is here as go prints a 1.0 float as 1 and we want a minimum
		// of one decimal place to pass the tests
		// if this is no longer a requirement, we can remove this check
		if t.literal.(float64) == float64(int(t.literal.(float64))) {
			literalStr = fmt.Sprintf("%.1f", t.literal)
		} else {
			literalStr = fmt.Sprintf("%g", t.literal)
		}
	} else {
		literalStr = fmt.Sprintf("%v", t.literal)
	}
	return literalStr
}
