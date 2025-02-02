package main

import (
	"fmt"
)

type Interpreter struct {
	environment *Envionment
	globals     *Envionment
	locals      map[Expr]int
}

func NewInterpreter() Interpreter {
	globals := NewEnvironment(nil)
	globals.define("clock", clock{})
	return Interpreter{
		globals:     globals,
		environment: globals,
		locals:      make(map[Expr]int),
	}
}

// VisitLiteralExpression will evaluate the literal expression
// which is just the value of the literal
func (i *Interpreter) VisitLiteralExpr(expr *Literal) (interface{}, error) {
	return expr.value, nil
}

// VisitLogicalExpr will evaluate the logical expression
// which is either the left or right expression
// depending on the operator and the truthiness of the left and right
func (i *Interpreter) VisitLogicalExpr(expr *Logical) (interface{}, error) {
	left, err := i.evaluate(expr.left)
	if err != nil {
		return nil, err
	}

	if expr.operator.tokenType == OR {
		if i.IsTruthy(left) {
			return left, nil
		}
	} else {
		// we are in an AND condition
		if !i.IsTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(expr.right)
}

// VisitSetExpr will evaluate the object whos property is being set and check
// to see if its a LoxInstance. If not, thats a runtime error. Otherwise,
// we evaluate the value being set and store it on the instance.
func (i *Interpreter) VisitSetExpr(expr *Set) (interface{}, error) {
	object, err := i.evaluate(expr.object)
	if err != nil {
		return nil, err
	}

	instance, ok := object.(*LoxInstance)
	if !ok {
		return nil, &RuntimeError{expr.name, "Only instances have fields."}
	}

	value, err := i.evaluate(expr.value)
	if err != nil {
		return nil, err
	}
	instance.set(expr.name, value)
	return value, nil
}

func (i *Interpreter) VisitThisExpr(expr *This) (interface{}, error) {
	return i.lookUpVariable(expr.keyword, expr)
}

func (i *Interpreter) VisitSuperExpr(expr *Super) (interface{}, error) {
	// look up the surrounding class's superclass up by looking up super
	// in the correct envionment
	distance := i.locals[expr]
	superclass := i.environment.getAt(distance, "super").(*LoxClass)

	// retrieve the current instance of ("this") by looking it up in the environment.
	// since "super" is stored one level higher in the environment chain,
	// we offset the lookup by one to access "this" from the inner environment.
	object := i.environment.getAt(distance-1, "this").(*LoxInstance)

	// find the method on teh superclass
	method, prs := superclass.findMethod(expr.method.lexeme)
	if !prs {
		return nil, &RuntimeError{expr.method, "Undefined property '" + expr.method.lexeme + "'."}
	}

	return method.bind(object), nil
}

// VisitGroupingExpr will evaluate the expression inside the grouping
func (i *Interpreter) VisitGroupingExpr(expr *Grouping) (interface{}, error) {
	return i.evaluate(expr.expression)
}

// VisitUnaryExpr will evaluate the unary expression
func (i *Interpreter) VisitUnaryExpr(expr *Unary) (interface{}, error) {
	right, err := i.evaluate(expr.right)
	if err != nil {
		return nil, err
	}

	switch expr.operator.tokenType {
	case BANG:
		return !i.IsTruthy(right), nil
	case MINUS:
		err := i.checkNumberOperand(expr.operator, right)
		if err != nil {
			return nil, err
		}
		return -right.(float64), err
	}

	return nil, nil
}

// VisitBinaryExpr will evaluate the binary expression
func (i *Interpreter) VisitBinaryExpr(expr *Binary) (interface{}, error) {
	left, err := i.evaluate(expr.left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.right)
	if err != nil {
		return nil, err
	}

	switch expr.operator.tokenType {
	case MINUS:
		err := i.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case SLASH:
		err := i.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case STAR:
		err := i.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	case PLUS:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l + r, nil
			}
		}
		if l, ok := left.(string); ok {
			if r, ok := right.(string); ok {
				return l + r, nil
			}
		}
		return nil, &RuntimeError{expr.operator, "Operands must be two numbers or two strings."}
	case GREATER:
		err := i.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case GREATER_EQUAL:
		err := i.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case LESS:
		err := i.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case LESS_EQUAL:
		err := i.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case BANG_EQUAL:
		return !i.isEqual(left, right), nil
	case EQUAL_EQUAL:
		return i.isEqual(left, right), nil
	}

	return nil, nil
}

func (i *Interpreter) VisitCallExpr(expr *Call) (interface{}, error) {
	callee, err := i.evaluate(expr.callee)
	if err != nil {
		return nil, err
	}

	var arguments []interface{}
	for _, argument := range expr.arguments {
		arg, err := i.evaluate(argument)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, arg)
	}

	if _, ok := callee.(Callable); !ok {
		return nil, &RuntimeError{expr.paren, "Can only call functions and classes."}
	}
	function := callee.(Callable)
	if len(arguments) != function.arity() {
		return nil, &RuntimeError{expr.paren, "Expected " +
			fmt.Sprintf("%d", function.arity()) +
			" arguments but got " +
			fmt.Sprintf("%d", len(arguments)) + "."}
	}
	return function.call(i, arguments)
}

// VisitGetExpr will evaluate the expression whos property is being accessed
// In Lox, only instances of classes have properties. If the object is some
// other type like a number, inboking a getter is a runtime error
func (i *Interpreter) VisitGetExpr(expr *Get) (interface{}, error) {
	object, err := i.evaluate(expr.object)
	if err != nil {
		return nil, err
	}

	if instance, ok := object.(*LoxInstance); ok {
		return instance.get(expr.name)
	}

	return nil, &RuntimeError{expr.name, "Only instances have properties."}
}

// VisitVarStmt will evaluate the variable statement
// and define the variable in the current environment
// if there is an initializer, it will evaluate the initializer
func (i *Interpreter) VisitVarStmt(stmt *Var) (interface{}, error) {
	var value interface{}
	var err error
	if stmt.initializer != nil {
		value, err = i.evaluate(stmt.initializer)
		if err != nil {
			return err, nil
		}
	}

	i.environment.define(stmt.name.lexeme, value)
	return nil, nil
}

// VisitWhileStmt will execute the while loop
// executing the statement body until the condition is no longer true
func (i *Interpreter) VisitWhileStmt(stmt *While) (interface{}, error) {
	value, err := i.evaluate(stmt.condition)
	if err != nil {
		return err, nil
	}
	for i.IsTruthy(value) {
		err := i.execute(stmt.body)
		if err != nil {
			return err, nil
		}
		value, err = i.evaluate(stmt.condition)
		if err != nil {
			return err, nil
		}
	}
	return nil, nil
}

// VisitAssignExpr will evaluate the assignment expression
// and assign the value to the variable in the current environment
// if the variable is not defined in the current environment
// it will throw a runtime error
func (i *Interpreter) VisitAssignExpr(expr *Assign) (interface{}, error) {
	value, err := i.evaluate(expr.value)
	if err != nil {
		return nil, err
	}

	distance, found := i.locals[expr]
	if found {
		i.environment.assignAt(distance, expr.name, value)
	} else {
		err := i.globals.assign(expr.name, value)
		if err != nil {
			return nil, err
		}
	}

	return value, nil
}

// VisitVariableExpr will evaluate the variable expression
// and return the value of the variable
func (i *Interpreter) VisitVariableExpr(expr *Variable) (interface{}, error) {
	return i.lookUpVariable(expr.name, expr)
}

func (i *Interpreter) lookUpVariable(name Token, expr Expr) (interface{}, error) {
	distance, found := i.locals[expr]
	if found {
		return i.environment.getAt(distance, name.lexeme), nil
	} else {
		return i.globals.get(name)
	}
}

// VisitExpressionStmt will evaluate the expression statement
// and return the value of the expression
func (i *Interpreter) VisitExpressionStmt(stmt *Expression) (interface{}, error) {
	_, err := i.evaluate(stmt.expression)
	return err, nil
}

// VisitFunctionStmt will define the function in the current environment
func (i *Interpreter) VisitFunctionStmt(stmt *Function) (interface{}, error) {
	function := &LoxFunction{*stmt, i.environment, false}
	i.environment.define(stmt.name.lexeme, function)
	return nil, nil
}

// VisitIfStmt will evaluate the if statement
// if the condition is truthy it will execute the then branch
// if there is an else branch and the condition is falsey
// it will execute the else branch
func (i *Interpreter) VisitIfStmt(stmt *If) (interface{}, error) {
	condition, err := i.evaluate(stmt.condition)
	if err != nil {
		return nil, err
	}

	if i.IsTruthy(condition) {
		err = i.execute(stmt.thenBranch)
		if err != nil {
			return nil, err
		}
	} else if stmt.elseBranch != nil {
		err = i.execute(stmt.elseBranch)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// VisitPrintStmt will evaluate the print statement
// and print the value of the expression
func (i *Interpreter) VisitPrintStmt(stmt *Print) (interface{}, error) {
	value, err := i.evaluate(stmt.expression)
	if err != nil {
		return err, nil
	}
	fmt.Println(i.stringify(value))
	return nil, nil
}

// VisitReturnStatement will evaluate the return statement
// and return the value of the expression as an error to be caught
// by the most recent function call
func (i *Interpreter) VisitReturnStmt(stmt *Return) (interface{}, error) {
	var value interface{}
	var err error

	if stmt.value != nil {
		value, err = i.evaluate(stmt.value)
		if err != nil {
			return err, nil
		}
	}
	return nil, &ReturnException{value}
}

// VisitBlockStmt will evaluate the block statement
func (i *Interpreter) VisitBlockStmt(stmt *Block) (interface{}, error) {
	return i.executeBlock(stmt.statements, NewEnvironment(i.environment)), nil
}

func (i *Interpreter) VisitClassStmt(stmt *Class) (interface{}, error) {
	// if the class has a superclass expression, we evaluate it.
	// since that could potentially evaluate to some other kind of object,
	// we have to check at runtime that the thing we want to be a superclass
	// is actually a class
	var superclass *LoxClass = nil
	if stmt.superclass != nil {
		superclassCandidate, err := i.evaluate(stmt.superclass)
		if err != nil {
			return err, nil
		}
		superclassValue, ok := superclassCandidate.(LoxClass)
		if !ok {
			return &RuntimeError{stmt.superclass.name, "Superclass must be a class"}, nil
		}
		superclass = &superclassValue
	}

	i.environment.define(stmt.name.lexeme, nil)

	// when we evaluate a subclass definition, we create a new envionment
	// and store a reference to the superclass. Then when we later create the
	// LoxFunctions for each method, those will capture the environent that
	// defines 'super' as thier closure - thus holding on to the superclass -
	if stmt.superclass != nil {
		i.environment = NewEnvironment(i.environment)
		i.environment.define("super", superclass)
	}

	var methods map[string]LoxFunction = make(map[string]LoxFunction)
	for _, method := range stmt.methods {
		function := LoxFunction{method, i.environment, method.name.lexeme == "init"}
		methods[method.name.lexeme] = function
	}

	var class LoxClass = LoxClass{stmt.name.lexeme, superclass, methods}

	// now that the methods have been created we pop the envionment defining
	// the superclass
	if superclass != nil {
		i.environment = i.environment.enclosing
	}

	i.environment.assign(stmt.name, class)
	return nil, nil
}

// executeBlock will execute the block of statements
// in a new environment. Before returning it will set the
// environment back to the previous environment
func (i *Interpreter) executeBlock(statements []Stmt, envionment *Envionment) error {
	previous := i.environment
	defer func() { i.environment = previous }()
	i.environment = envionment
	for _, statement := range statements {
		err := i.execute(statement)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) interpret(statements []Stmt) {
	for _, statement := range statements {
		err := i.execute(statement)
		if err != nil {
			reportRuntimeError(*(err.(*RuntimeError)))
			return
		}
	}
}

func (i *Interpreter) execute(statement Stmt) error {
	_, err := statement.Accept(i)
	return err
}

func (i *Interpreter) resolve(expr Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) stringify(object interface{}) string {
	if object == nil {
		return "nil"
	}

	if fnum, ok := object.(float64); ok {
		if fnum == float64(int(fnum)) {
			return fmt.Sprintf("%.0f", fnum)
		} else {
			return fmt.Sprintf("%g", fnum)
		}
	}
	return fmt.Sprintf("%v", object)
}

// IsTruthy will return true if the value is not nil or false
func (i *Interpreter) IsTruthy(value interface{}) bool {
	switch value := value.(type) {
	case nil:
		return false
	case bool:
		return value
	}
	return true
}

// isEqual will compare two values and return true if they are equal
// this currently uses golangs == operator which can be adjusted to
// match the behaviour of lox
func (i *Interpreter) isEqual(a interface{}, b interface{}) bool {
	return a == b
}

// Evaluate will evaluate the expression
func (i *Interpreter) evaluate(expr Expr) (interface{}, error) {
	return expr.Accept(i)
}

// checkNumberOperand will check if the operand is a number
func (i *Interpreter) checkNumberOperand(operator Token, operand interface{}) error {
	if _, ok := operand.(float64); ok {
		return nil
	}
	return &RuntimeError{operator, "Operand must be a number,"}
}

// checkNumberOperands will check if the operands are numbers
func (i *Interpreter) checkNumberOperands(operator Token, left interface{}, right interface{}) error {
	if _, ok := left.(float64); ok {
		if _, ok := right.(float64); ok {
			return nil
		}
	}
	return &RuntimeError{operator, "Operands must be numbers."}
}
