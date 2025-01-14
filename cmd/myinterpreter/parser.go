package main

import (
	"fmt"
)

type Parser struct {
	tokens  []Token
	current int
}

// parse is the entry point for the parser
func (p *Parser) parse() []Stmt {
	var statements []Stmt
	for !p.isAtEnd() {
		declaration := p.declaration()
		statements = append(statements, declaration)
	}
	return statements
}

// declaration represents the declaration rule of the grammar
// declaration -> varDecl | statement | funDecl
func (p *Parser) declaration() Stmt {
	var stmt Stmt
	var err error

	if p.match(FUN) {
		stmt, err = p.function("function")
	} else if p.match(VAR) {
		stmt, err = p.varDeclaration()
	} else {
		stmt, err = p.statement()
	}

	if err != nil {
		p.synchonize()
		return nil
	}
	return stmt
}

// function represents the function rule of the grammar
// function -> IDENTIFIER "(" parameters? ")" block;
func (p *Parser) function(kind string) (Stmt, error) {
	name, err := p.consume(IDENTIFIER, "Expect "+kind+"name.")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(LEFT_PAREN, "Expect '(' after "+kind+" name.")
	if err != nil {
		return nil, err
	}
	var parameters []Token
	if !p.check(RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				p.error(p.peek(), "Can't have more than 255 parameters.")
			}

			parameter, err := p.consume(IDENTIFIER, "Expect parameter name.")
			if err != nil {
				return nil, err
			}
			parameters = append(parameters, parameter)

			if !p.match(COMMA) {
				break
			}
		}
	}
	_, err = p.consume(RIGHT_PAREN, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LEFT_BRACE, "Expect '{' before"+kind+" body")
	if err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}
	return &Function{name, parameters, body}, nil
}

// varDeclaration represents the var declaration rule of the grammar
// varDecl -> "var" IDENTIFIER ( "=" expression )? ";";
func (p *Parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}
	var initializer Expr
	if p.match(EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(SEMICOLON, "Expect ';' after variable declaration")
	if err != nil {
		return nil, err
	}
	return &Var{name, initializer}, nil
}

// represents the statment rule of the grammar
// statement -> exprStmt | ifStmt | printStmt | while | block
func (p *Parser) statement() (Stmt, error) {
	if p.match(FOR) {
		return p.forStatement()
	}
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(RETURN) {
		return p.returnStatement()
	}
	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(LEFT_BRACE) {
		statements, err := p.block()
		if err != nil {
			return nil, err
		}
		return &Block{statements}, nil
	}
	return p.expressionStatement()
}

// represents the return statement rule of the grammar
// returnStmt -> "return" expression? ";";
func (p *Parser) returnStatement() (Stmt, error) {
	var err error
	keyword := p.previous()
	var value Expr = nil

	if !p.check(SEMICOLON) {
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(SEMICOLON, "Expect ';' after return value.")
	if err != nil {
		return nil, err
	}
	return &Return{keyword, value}, nil
}

// represents the for statement rule of the grammar
// forStmt -> "for" "(" ( varDecl | exprStmt | ";" )
//
//	expression? ";"
//	expression? ")" statement ;
func (p *Parser) forStatement() (Stmt, error) {
	_, err := p.consume(LEFT_PAREN, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	var initializer Stmt
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	var condition Expr
	if !p.check(SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(SEMICOLON, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	var increment Expr
	if !p.check(RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(RIGHT_PAREN, "Expect ')' after for clauses.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	// desuguaring the for loop into a while loop structure
	if increment != nil {
		body = &Block{[]Stmt{body, &Expression{increment}}}
	}
	if condition == nil {
		condition = &Literal{true}
	}
	body = &While{condition: condition, body: body}
	if initializer != nil {
		body = &Block{[]Stmt{initializer, body}}
	}
	return body, nil
}

// represents the while statement rule of the grammar
// whileStmt -> "while" "(" expression ")" statement ;
func (p *Parser) whileStatement() (Stmt, error) {
	_, err := p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RIGHT_PAREN, "Expect ')' after condition.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &While{condition: condition, body: body}, nil
}

// represents the if statement rule of the grammar
// ifStmt -> "if" "(" expression ")" statement
//
//	( "else" statement )? ;
func (p *Parser) ifStatement() (Stmt, error) {
	_, err := p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(RIGHT_PAREN, "Expect ')' after if condition.")
	if err != nil {
		return nil, err
	}
	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}
	var elseBranch Stmt = nil
	if p.match(ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}
	return &If{condition: condition, thenBranch: thenBranch, elseBranch: elseBranch}, nil
}

// represents the block rule of the grammar
// block -> "{" declaration* "{";
func (p *Parser) block() ([]Stmt, error) {
	var statements []Stmt
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	_, err := p.consume(RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}
	return statements, nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(SEMICOLON, "Expect ';' after expression.")
	if err != nil {
		return nil, err
	}
	return &Expression{value}, nil
}

// represents the assignment rule of the grammar
// assignment -> IDENTIFIER "=" assignment | logic_or ;
func (p *Parser) assignment() (Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}
		variable, ok := expr.(*Variable)
		if ok {
			name := variable.name
			return &Assign{name, value}, nil
		}

		p.error(equals, "Invalid assignment target.")
	}
	return expr, nil
}

func (p *Parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}
	for p.match(OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = &Logical{expr, operator, right}
	}

	return expr, nil
}

func (p *Parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = &Logical{expr, operator, right}
	}
	return expr, nil
}

func (p *Parser) printStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}
	return &Print{value}, nil
}

// NewParser creates a new Parser with the provided tokens
func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, current: 0}
}

// represents the expression rule of the grammar
// expression -> assignment
func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

// represents the equality rule of the grammar
// equality -> comparison ( ( "!=" | "==" ) comparison )*
func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &Binary{expr, operator, right}
	}
	return expr, nil
}

// represents the comparison rule of the grammar
// comparison -> term ( ( ">" | ">=" | "<" | "<=" ) term )*
func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &Binary{expr, operator, right}
	}
	return expr, nil
}

// represents the term rule of the grammar
// term -> factor ( ( "-" | "+" ) factor )*
func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &Binary{expr, operator, right}
	}
	return expr, nil
}

// represents the factor rule of the grammar
// factor -> unary ( ( "/" | "*" ) unary )*
func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &Binary{expr, operator, right}
	}
	return expr, nil
}

// represents the unary rule of the grammar
// unary -> ( "!" | "-" ) unary | primary
func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &Unary{operator: operator, right: right}, nil
	}
	return p.call()
}

func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(LEFT_PAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(callee Expr) (Expr, error) {
	var arguments []Expr
	if !p.check(RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				p.error(p.peek(), "Can't have more than 255 arguements.")
			}
			value, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, value)
			if !p.match(COMMA) {
				break
			}
		}
	}

	paren, err := p.consume(RIGHT_PAREN, "Expect ')' after arguements.")
	if err != nil {
		return nil, err
	}

	return &Call{callee: callee, paren: paren, arguments: arguments}, nil
}

// represents the primary rule of the grammar
// primary -> NUMBER | STRING | "false" | "true" | "nil" | "(" expression ")"
func (p *Parser) primary() (Expr, error) {
	if p.match(FALSE) {
		return &Literal{false}, nil
	} else if p.match(TRUE) {
		return &Literal{true}, nil
	} else if p.match(NIL) {
		return &Literal{nil}, nil
	} else if p.match(NUMBER, STRING) {
		return &Literal{p.previous().literal}, nil
	} else if p.match(IDENTIFIER) {
		return &Variable{p.previous()}, nil
	} else if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return &Grouping{expr}, nil
	}
	err := p.error(p.peek(), "Expect expression.")
	return nil, err
}

// consume consumes the current token if it is of the provided type
// otherwise it will throw an error
func (p *Parser) consume(tokenType TokenType, message string) (Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	err := p.error(p.peek(), message)
	return Token{}, err
}

// report prints an error message to the console
func (p *Parser) error(token Token, message string) ParseError {
	if token.tokenType == EOF {
		report(token.line, " at end", message)
	} else {
		report(token.line, " at '"+token.lexeme+"'", message)
	}
	return ParseError{token, message}
}

// synchonize will skip tokens until it finds a statement boundary
func (p *Parser) synchonize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().tokenType == SEMICOLON {
			return
		}

		switch p.peek().tokenType {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}
		p.advance()
	}
}

// ParseError represents an error that occurred during parsing
type ParseError struct {
	token   Token
	message string
}

// Error returns a formatted error message
func (e ParseError) Error() string {
	return fmt.Sprintf("[line %d] Error%s: %s", e.token.line, e.token.lexeme, e.message)
}

// match checks to see if the current token is one of the provided types
// if it is, then it will consume the token and return true otherwise
// the token will be left alone and will return false
func (p *Parser) match(types ...TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

// advance consumes the current token and returns it
func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// check returns true if the current token is of the provided type
// this will never consume the token it only looks at it
func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().tokenType == tokenType
}

// checks if we have run out of tokens to parse
func (p *Parser) isAtEnd() bool {
	return p.peek().tokenType == EOF
}

// peek returns the current token
func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

// previous returns the most recently consumed token
func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}
