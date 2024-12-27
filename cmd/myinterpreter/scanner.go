package main

type Scanner struct {
	source string
	tokens []Token
	// points to the first character in the lexeme being scanned
	start int
	// points to the character currently being considered
	current int
	// what source line currernt is on
	line int
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
		error(s.line, "Unexpected character: "+string(c))
	}
}

// string consumes characters until it reaches the " that ends the string.
// Will also gracefully handle running out of input until the string is
// closed and report the error
func (s *Scanner) string() {
	// consume the contents of the string
	for s.peek() != '"' && !s.isAtEnd(){
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		error(s.line, "Unterminated string.")
		return
	}

	// consume the closing "
	s.advance()

	// Trim off the surrounding quotes from the string value
	var value string = s.source[s.start + 1: s.current - 1]
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
