package main

type Interpreter struct{}

// VisitLiteralExpression will evaluate the literal expression
// which is just the value of the literal
func (i *Interpreter) VisitLiteralExpr(expr *Literal) interface{} {
	return expr.value
}

// VisitGroupingExpr will evaluate the expression inside the grouping
func (i *Interpreter) VisitGroupingExpr(expr *Grouping) interface{} {
	return i.evaluate(expr.expression)
}

// VisitUnaryExpr will evaluate the unary expression
func (i *Interpreter) VisitUnaryExpr(expr *Unary) interface{} {
	right := i.evaluate(expr.right)

	switch expr.operator.tokenType {
	case BANG:
		return !i.IsTruthy(right)
	case MINUS:
		return -right.(float64)
	}

	return nil
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

// VisitBinaryExpr will evaluate the binary expression
func (i *Interpreter) VisitBinaryExpr(expr *Binary) interface{} {
	left := i.evaluate(expr.left)
	right := i.evaluate(expr.right)

	switch expr.operator.tokenType {
	case MINUS:
		return left.(float64) - right.(float64)
	case SLASH:
		return left.(float64) / right.(float64)
	case STAR:
		return left.(float64) * right.(float64)
	case PLUS:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l + r
			}
		}
		if l, ok := left.(string); ok {
			if r, ok := right.(string); ok {
				return l + r
			}
		}
	case GREATER:
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case LESS:
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		return left.(float64) <= right.(float64)
	case BANG_EQUAL:
		return !i.isEqual(left, right)
	case EQUAL_EQUAL:
		return i.isEqual(left, right)
	}

	return nil
}

// isEqual will compare two values and return true if they are equal
// this currently uses golangs == operator which can be adjusted to
// match the behaviour of lox
func (i *Interpreter) isEqual(a interface{}, b interface{}) bool {
	return a == b
}

// Evaluate will evaluate the expression
func (i *Interpreter) evaluate(expr Expr) interface{} {
	return expr.Accept(i)
}
