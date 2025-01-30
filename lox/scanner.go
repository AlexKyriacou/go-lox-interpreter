package main

import (
	"strconv"
)

type Scanner struct {
	source string
	tokens []Token

	// points to the first character in the lexeme being scanned
	start int

	// points to the character currently being considered
	current int

	// what source line currernt is on
	line int

	// map of all reesrved identifiers
	keywords map[string]TokenType
}

// NewScanner creates a new Scanner with the provided source code
// initializes the Scanner's properties including the reserved keywords
// returning a pointer to the Scanner
func NewScanner(source string) *Scanner {
	s := &Scanner{source: source, tokens: []Token{}, start: 0, current: 0, line: 1}
	s.keywords = map[string]TokenType{
		"and":    AND,
		"class":  CLASS,
		"else":   ELSE,
		"false":  FALSE,
		"for":    FOR,
		"fun":    FUN,
		"if":     IF,
		"nil":    NIL,
		"or":     OR,
		"print":  PRINT,
		"return": RETURN,
		"super":  SUPER,
		"this":   THIS,
		"true":   TRUE,
		"var":    VAR,
		"while":  WHILE,
	}
	return s
}

// scanTokens loops through the source code of the Scanner
// and adds tokens to the Scanners tokens slice until it
// runs out of characters.
func (s *Scanner) scanTokens() []Token {
	s.start = 0
	s.current = 0
	s.line = 1
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	// append EOF token
	s.tokens = append(s.tokens, Token{EOF, "", nil, s.line})
	return s.tokens
}

// isAtEnd returns true if the Scanner has consumed all the characters
func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

// scanToken consumes the current token and adds it to the Scanner's tokens slice
// advancing the Scanner's current pointer to the next lexeme
func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '/':
		if s.match('/') {
			// We are at a comment, consume the rest of the line
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
	case ' ':
		// ignore whitespace
	case '\r':
		// ignore whitespace
	case '\t':
		// ignore whitespace
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			report(s.line, "", "Unexpected character: "+string(c))
		}
	}
}

// identifier consumes all the characters within a valid identifier and adds the
// token to the scanner checks against the scanner reserved words list
func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	tokenType, isReserved := s.keywords[text]
	if !isReserved {
		tokenType = IDENTIFIER
	}
	s.addToken(tokenType)
}

// isAlphaNumber returns true if the character is a valid alpha character or
// digit for a lox identifier
func (s *Scanner) isAlphaNumeric(c byte) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

// return true if the character is a valid alpha character for a lox identifier
func (s *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

// isDigit will return a bool based on if a char is within 0-9
// we need our own function as the unicode.IsDigit func
// also includes Devangari digits and other stuff not required
func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// number will consume as many digits as found for the integer part of the literal
// it will look for a decimal point that is followed by at least one digit
// if a fractional part exists it will then consume any lettes found after
func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	// look for a fractional part of the number
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		//consume the decimal point
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	value, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		report(s.line, "", "Invalid number.")
		return
	}
	s.addTokenLiteral(NUMBER, value)
}

// string consumes characters until it reaches the " that ends the string.
// Will also gracefully handle running out of input until the string is
// closed and report the error
func (s *Scanner) string() {
	// consume the contents of the string
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		report(s.line, "", "Unterminated string.")
		return
	}

	// consume the closing "
	s.advance()

	// Trim off the surrounding quotes from the string value
	var value string = s.source[s.start+1 : s.current-1]
	s.addTokenLiteral(STRING, value)
}

// peek returns the current character without consuming it
// will return null terminator if the Scanner is at the end of the source
func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\000'
	}
	return s.source[s.current]
}

// peekNext return the character 2 ahead without consuming it
// will return null terminator if the Scanner is 2 away from the end of the source
func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return '\000'
	}
	return s.source[s.current+1]
}

// advance consumes the current character and returns it
func (s *Scanner) advance() byte {
	current := s.source[s.current]
	s.current++
	return current
}

// addToken grabs the text of the current lexeme and adds a token tot he Scanner's tokens slice
func (s *Scanner) addToken(tokenType TokenType) {
	s.addTokenLiteral(tokenType, nil)
}

// addTokenLiteral grabs the text of the current lexeme and adds a token to the Scanner's tokens slice
func (s *Scanner) addTokenLiteral(tokenType TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{tokenType, text, literal, s.line})
}

// match returns true if the current character matches the expected character
// if it does, it advances the Scanner's current pointer
func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}

	s.current++

	return true
}
