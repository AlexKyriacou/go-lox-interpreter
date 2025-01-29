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
func (i *Interpreter) VisitVarStmt(stmt *Var) error {
	var value interface{}
	var err error
	if stmt.initializer != nil {
		value, err = i.evaluate(stmt.initializer)
		if err != nil {
			return err
		}
	}

	i.environment.define(stmt.name.lexeme, value)
	return nil
}

// VisitWhileStmt will execute the while loop
// executing the statement body until the condition is no longer true
func (i *Interpreter) VisitWhileStmt(stmt *While) error {
	value, err := i.evaluate(stmt.condition)
	if err != nil {
		return err
	}
	for i.IsTruthy(value) {
		err := i.execute(stmt.body)
		if err != nil {
			return err
		}
		value, err = i.evaluate(stmt.condition)
		if err != nil {
			return err
		}
	}
	return nil
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

func (i *Interpreter) lookUpVariable(name Token, expr *Variable) (interface{}, error) {
	distance, found := i.locals[expr]
	if found {
		return i.environment.getAt(distance, name.lexeme), nil
	} else {
		return i.globals.get(name)
	}
}

// VisitExpressionStmt will evaluate the expression statement
// and return the value of the expression
func (i *Interpreter) VisitExpressionStmt(stmt *Expression) error {
	_, err := i.evaluate(stmt.expression)
	return err
}

// VisitFunctionStmt will define the function in the current environment
func (i *Interpreter) VisitFunctionStmt(stmt *Function) error {
	function := &LoxFunction{declaration: *stmt, closure: i.environment}
	i.environment.define(stmt.name.lexeme, function)
	return nil
}

// VisitIfStmt will evaluate the if statement
// if the condition is truthy it will execute the then branch
// if there is an else branch and the condition is falsey
// it will execute the else branch
func (i *Interpreter) VisitIfStmt(stmt *If) error {
	condition, err := i.evaluate(stmt.condition)
	if err != nil {
		return err
	}

	if i.IsTruthy(condition) {
		err = i.execute(stmt.thenBranch)
		if err != nil {
			return err
		}
	} else if stmt.elseBranch != nil {
		err = i.execute(stmt.elseBranch)
		if err != nil {
			return err
		}
	}
	return nil
}

// VisitPrintStmt will evaluate the print statement
// and print the value of the expression
func (i *Interpreter) VisitPrintStmt(stmt *Print) error {
	value, err := i.evaluate(stmt.expression)
	if err != nil {
		return err
	}
	fmt.Println(i.stringify(value))
	return nil
}

// VisitReturnStatement will evaluate the return statement
// and return the value of the expression as an error to be caught
// by the most recent function call
func (i *Interpreter) VisitReturnStmt(stmt *Return) error {
	var value interface{}
	var err error

	if stmt.value != nil {
		value, err = i.evaluate(stmt.value)
		if err != nil {
			return err
		}
	}
	return &ReturnException{value}
}

// VisitBlockStmt will evaluate the block statement
func (i *Interpreter) VisitBlockStmt(stmt *Block) error {
	return i.executeBlock(stmt.statements, NewEnvironment(i.environment))
}

func (i *Interpreter) VisitClassStmt(stmt *Class) error {
	i.environment.define(stmt.name.lexeme, nil)

	var methods map[string]LoxFunction = make(map[string]LoxFunction)
	for _, method := range stmt.methods {
		function := LoxFunction{method, i.environment}
		methods[method.name.lexeme] = function
	}

	var class LoxClass = LoxClass{stmt.name.lexeme, methods}
	i.environment.assign(stmt.name, class)
	return nil
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
	return statement.Accept(i)
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
