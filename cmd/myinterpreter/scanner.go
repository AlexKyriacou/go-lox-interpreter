package main

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func (s *Scanner) scanTokens() []Token {
	s.start = 0
	s.current = 0
	s.line = 1
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{EOF, "", nil, s.line})
	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

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
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
	case ' ':
		fallthrough
	case '\r':
		fallthrough
	case '\t':
		// ignore whitespace
	case '\n':
		s.line++
	default:
		error(s.line, "Unexpected character: "+string(c))
	}
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\000'
	}
	return s.source[s.current]
}

func (s *Scanner) advance() byte {
	current := s.source[s.current]
	s.current++
	return current
}

func (s *Scanner) addToken(tokenType TokenType) {
	s.addTokenLiteral(tokenType, nil)
}

func (s *Scanner) addTokenLiteral(tokenType TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{tokenType, text, literal, s.line})
}

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
