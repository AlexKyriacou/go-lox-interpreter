package main

type Resolver struct {
	scopes          stack[map[string]bool]
	currentFunction FunctionType
	interpreter     *Interpreter
}

type FunctionType int

const (
	NONE FunctionType = iota
	FUNCTION
)

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{interpreter: interpreter, currentFunction: NONE, scopes: *NewStack[map[string]bool]()}
}

func (r *Resolver) endScope() {
	r.scopes.Pop()
}

func (r *Resolver) beginScope() {
	r.scopes.Push(make(map[string]bool))
}

func (r *Resolver) define(name Token) {
	if r.scopes.isEmpty() {
		return
	}
	r.scopes.Peek()[name.lexeme] = true
}

func (r *Resolver) declare(name Token) {
	if r.scopes.isEmpty() {
		return
	}
	scope := r.scopes.Peek()
	if _, ok := scope[name.lexeme]; ok {
		r.error(name, "Already a variable with this name in this scope.")
	}
	scope[name.lexeme] = false
}

// report prints an error message to the console
func (r *Resolver) error(token Token, message string) {
	if token.tokenType == EOF {
		report(token.line, " at end", message)
	} else {
		report(token.line, " at '"+token.lexeme+"'", message)
	}
}

func (r *Resolver) resolveStatements(statements []Stmt) {
	for _, statement := range statements {
		r.resolveStatement(statement)
	}
}

func (r *Resolver) resolveStatement(stmt Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveExpression(expr Expr) {
	expr.Accept(r)
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := r.scopes.Len() - 1; i >= 0; i-- {
		_, present := r.scopes.get(i)[name.lexeme]
		if present {
			r.interpreter.resolve(expr, r.scopes.Len()-1-i)
			return
		}
	}
}

func (r *Resolver) VisitExpressionStmt(stmt *Expression) error {
	r.resolveExpression(stmt.expression)
	return nil
}

func (r *Resolver) VisitFunctionStmt(stmt *Function) error {
	r.declare(stmt.name)
	r.define(stmt.name)

	r.resolveFunction(stmt, FUNCTION)
	return nil
}

func (r *Resolver) resolveFunction(function *Function, funcType FunctionType) {
	var enclosingFunction FunctionType = r.currentFunction
	r.currentFunction = funcType

	r.beginScope()
	for _, param := range function.params {
		r.declare(param)
		r.define(param)
	}
	r.resolveStatements(function.body)
	r.endScope()
	r.currentFunction = enclosingFunction
}

func (r *Resolver) VisitIfStmt(stmt *If) error {
	r.resolveExpression(stmt.condition)
	r.resolveStatement(stmt.thenBranch)
	if stmt.elseBranch != nil {
		r.resolveStatement(stmt.elseBranch)
	}
	return nil
}

func (r *Resolver) VisitPrintStmt(stmt *Print) error {
	r.resolveExpression(stmt.expression)
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt *Return) error {
	if r.currentFunction == NONE{
		r.error(stmt.keyword, "Can't return from top-level code.")
	}
	
	if stmt.value != nil {
		r.resolveExpression(stmt.value)
	}
	return nil
}

func (r *Resolver) VisitVarStmt(stmt *Var) error {
	r.declare(stmt.name)
	if stmt.initializer != nil {
		r.resolveExpression(stmt.initializer)
	}
	r.define(stmt.name)
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt *While) error {
	r.resolveExpression(stmt.condition)
	r.resolveStatement(stmt.body)
	return nil
}

func (r *Resolver) VisitBlockStmt(stmt *Block) error {
	r.beginScope()
	r.resolveStatements(stmt.statements)
	r.endScope()
	return nil
}

func (r *Resolver) VisitAssignExpr(expr *Assign) (interface{}, error) {
	r.resolveExpression(expr.value)
	r.resolveLocal(expr, expr.name)
	return nil, nil
}

func (r *Resolver) VisitBinaryExpr(expr *Binary) (interface{}, error) {
	r.resolveExpression(expr.left)
	r.resolveExpression(expr.right)
	return nil, nil
}

func (r *Resolver) VisitCallExpr(expr *Call) (interface{}, error) {
	r.resolveExpression(expr.callee)
	for _, argument := range expr.arguments {
		r.resolveExpression(argument)
	}
	return nil, nil
}

func (r *Resolver) VisitGroupingExpr(expr *Grouping) (interface{}, error) {
	r.resolveExpression(expr.expression)
	return nil, nil
}

func (r *Resolver) VisitLiteralExpr(expr *Literal) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) VisitLogicalExpr(expr *Logical) (interface{}, error) {
	r.resolveExpression(expr.left)
	r.resolveExpression(expr.right)
	return nil, nil
}

func (r *Resolver) VisitUnaryExpr(expr *Unary) (interface{}, error) {
	r.resolveExpression(expr.right)
	return nil, nil
}

func (r *Resolver) VisitVariableExpr(expr *Variable) (interface{}, error) {
	isDefined := r.scopes.Peek()[expr.name.lexeme]
	if r.scopes.isEmpty() && !isDefined {
		report(expr.name.line, "", "Can't read local variable in its own initializer.")
	}
	r.resolveLocal(expr, expr.name)
	return nil, nil
}
