package main

import (
	"fmt"
	"strings"
)

type Interpreter struct {
	environment *Envionment
}

func NewInterpreter() Interpreter {
	return Interpreter{NewEnvironment(nil)}
}

// VisitLiteralExpression will evaluate the literal expression
// which is just the value of the literal
func (i *Interpreter) VisitLiteralExpr(expr *Literal) (interface{}, error) {
	return expr.value, nil
}

//  VisitLogicalExpr will evaluate the logical expression
//  which is either the left or right expression
//  depending on the operator and the truthiness of the left and right
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

// VisitAssignExpr will evaluate the assignment expression
// and assign the value to the variable in the current environment
// if the variable is not defined in the current environment
// it will throw a runtime error
func (i *Interpreter) VisitAssignExpr(expr *Assign) (interface{}, error) {
	value, err := i.evaluate(expr.value)
	if err != nil {
		return nil, err
	}
	err = i.environment.assign(expr.name, value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// VisitVariableExpr will evaluate the variable expression
// and return the value of the variable
func (i *Interpreter) VisitVariableExpr(expr *Variable) (interface{}, error) {
	return i.environment.get(expr.name)
}

// VisitExpressionStmt will evaluate the expression statement
// and return the value of the expression
func (i *Interpreter) VisitExpressionStmt(stmt *Expression) error {
	_, err := i.evaluate(stmt.expression)
	return err
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

// VisitBlockStmt will evaluate the block statement
func (i *Interpreter) VisitBlockStmt(stmt *Block) error {
	return i.executeBlock(stmt.statements, NewEnvironment(i.environment))
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

func (i *Interpreter) stringify(object interface{}) string {
	if object == nil {
		return "nil"
	}

	if _, ok := object.(float64); ok {
		text := fmt.Sprintf("%v", object)
		if strings.HasSuffix(text, ".0") {
			text = text[0:]
		}
		return text
	}
	return fmt.Sprintf("%v", object)
}
