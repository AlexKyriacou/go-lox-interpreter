package main

type Resolver struct {
	scopes          stack[map[string]bool]
	currentFunction FunctionType
	currentClass    ClassType
	interpreter     *Interpreter
}

type FunctionType int

const (
	FUNCTION_NONE FunctionType = iota
	FUNCTION_FUNCTION
	FUNCTION_INITIALIZER
	FUNCTION_METHOD
)

type ClassType int

const (
	CLASS_NONE ClassType = iota
	CLASS_CLASS
	CLASS_SUBCLASS
)

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{interpreter: interpreter, currentFunction: FUNCTION_NONE, currentClass: CLASS_NONE, scopes: *NewStack[map[string]bool]()}
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

func (r *Resolver) VisitExpressionStmt(stmt *Expression) (interface{}, error) {
	r.resolveExpression(stmt.expression)
	return nil, nil
}

func (r *Resolver) VisitFunctionStmt(stmt *Function) (interface{}, error) {
	r.declare(stmt.name)
	r.define(stmt.name)

	r.resolveFunction(stmt, FUNCTION_FUNCTION)
	return nil, nil
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

func (r *Resolver) VisitIfStmt(stmt *If) (interface{}, error) {
	r.resolveExpression(stmt.condition)
	r.resolveStatement(stmt.thenBranch)
	if stmt.elseBranch != nil {
		r.resolveStatement(stmt.elseBranch)
	}
	return nil, nil
}

func (r *Resolver) VisitPrintStmt(stmt *Print) (interface{}, error) {
	r.resolveExpression(stmt.expression)
	return nil, nil
}

func (r *Resolver) VisitReturnStmt(stmt *Return) (interface{}, error) {
	if r.currentFunction == FUNCTION_NONE {
		r.error(stmt.keyword, "Can't return from top-level code.")
	}

	if stmt.value != nil {
		if r.currentFunction == FUNCTION_INITIALIZER {
			r.error(stmt.keyword, "Can't return a value from an initializer.")
		}
		r.resolveExpression(stmt.value)
	}
	return nil, nil
}

func (r *Resolver) VisitVarStmt(stmt *Var) (interface{}, error) {
	r.declare(stmt.name)
	if stmt.initializer != nil {
		r.resolveExpression(stmt.initializer)
	}
	r.define(stmt.name)
	return nil, nil
}

func (r *Resolver) VisitWhileStmt(stmt *While) (interface{}, error) {
	r.resolveExpression(stmt.condition)
	r.resolveStatement(stmt.body)
	return nil, nil
}

func (r *Resolver) VisitBlockStmt(stmt *Block) (interface{}, error) {
	r.beginScope()
	r.resolveStatements(stmt.statements)
	r.endScope()
	return nil, nil
}

func (r *Resolver) VisitClassStmt(stmt *Class) (interface{}, error) {
	var enclosingClass ClassType = r.currentClass
	r.currentClass = CLASS_CLASS
	defer func() { r.currentClass = enclosingClass }()

	r.declare(stmt.name)
	r.define(stmt.name)

	// prevent the case of self inheritance i.e.
	// class Foo < Foo {}
	if stmt.superclass != nil && stmt.name.lexeme == stmt.superclass.name.lexeme {
		r.error(stmt.superclass.name, "A class can't inherit from itself.")
	}

	// traverse into and resolve the superclass subexpression.
	// since classes are usually declared at the top level, this is unlikely
	// to do anything useful however since Lox allows class declarations even
	// inside blocks, its possible the superclass name refers to a local
	// variable. In this case, we need to make sure its resolved
	if stmt.superclass != nil {
		r.currentClass = CLASS_SUBCLASS
		r.resolveExpression(stmt.superclass)
	}

	// if the class definition has a superclass, then we create a new scope
	// surrounding all of its methods. In that scope, we define the name "super"
	// Once we're done resolving the classes methods, we discard the scope
	if stmt.superclass != nil {
		r.beginScope()
		r.scopes.Peek()["super"] = true
		defer r.endScope()
	}

	r.beginScope()
	r.scopes.Peek()["this"] = true

	for _, method := range stmt.methods {
		declaration := FUNCTION_METHOD
		if method.name.lexeme == "init" {
			declaration = FUNCTION_INITIALIZER
		}
		r.resolveFunction(&method, declaration)
	}

	r.endScope()

	return nil, nil
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

func (r *Resolver) VisitGetExpr(expr *Get) (interface{}, error) {
	r.resolveExpression(expr.object)
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

func (r *Resolver) VisitSetExpr(expr *Set) (interface{}, error) {
	r.resolveExpression(expr.value)
	r.resolveExpression(expr.object)
	return nil, nil
}

// the super token is resolved as if it was a variable. The resolution
// stores the number of envionment 'hops' the interpreter needs to walk
// to find the envionment where the superclass is stored
func (r *Resolver) VisitSuperExpr(expr *Super) (interface{}, error) {
	// check to see if currently inside of a subclass
	if r.currentClass == CLASS_NONE {
		r.error(expr.keyword, "Can't use 'super' outside of a class.")
	} else if r.currentClass != CLASS_SUBCLASS {
		r.error(expr.keyword, "Can't use 'super' in a class with no superclass.")
	}

	r.resolveLocal(expr, expr.keyword)
	return nil, nil
}

func (r *Resolver) VisitThisExpr(expr *This) (interface{}, error) {
	if r.currentClass == CLASS_NONE {
		r.error(expr.keyword, "Can't use 'this' outside of a class.")
		return nil, nil
	}
	r.resolveLocal(expr, expr.keyword)
	return nil, nil
}

func (r *Resolver) VisitUnaryExpr(expr *Unary) (interface{}, error) {
	r.resolveExpression(expr.right)
	return nil, nil
}

func (r *Resolver) VisitVariableExpr(expr *Variable) (interface{}, error) {
	if !r.scopes.isEmpty() {
		if val, ok := r.scopes.Peek()[expr.name.lexeme]; ok && !val {
			report(expr.name.line, "", "Can't read local variable in its own initializer.")
		}
	}
	r.resolveLocal(expr, expr.name)
	return nil, nil
}
