package main

import (
	"errors"
	"fmt"
	"strings"
)

type Interpreter struct{}

// VisitLiteralExpression will evaluate the literal expression
// which is just the value of the literal
func (i *Interpreter) VisitLiteralExpr(expr *Literal) (interface{}, error) {
	return expr.value, nil
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

func (i *Interpreter) checkNumberOperand(operator Token, operand interface{}) error {
	if _, ok := operand.(float64); ok {
		return nil
	}
	return &RuntimeError{operator, "Operand must be a number,"}
}

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

func (i *Interpreter) interpret(expression Expr) {
	value, err := i.evaluate(expression)
	if errors.Is(err, &RuntimeError{}){
		reportRuntimeError(*(err.(*RuntimeError)))
		return
	}
	fmt.Println(i.stringify(value))
}

func (i *Interpreter) stringify(object interface{}) string{
	if object == nil{
		return "nil"
	}

	if _, ok := object.(float64); ok {
		text := fmt.Sprintf("%v", object)
		if strings.HasSuffix(text, ".0"){
			text = text[0:]
		}
		return text
	}
	return fmt.Sprintf("%v", object)
}