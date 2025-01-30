package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

// VisitSetExpr implements ExprVisitor.
func (p *AstPrinter) VisitSetExpr(expr *Set) (interface{}, error) {
	return p.parenthesize("set", expr.object, expr.name, expr.value)
}

func (p *AstPrinter) VisitSuperExpr(expr *Super) (interface{}, error) {
	return p.parenthesize("super", expr.keyword, expr.method)
}

func (p *AstPrinter) VisitThisExpr(expr *This) (interface{}, error) {
	return expr.keyword.lexeme, nil
}

// VisitGetExpr implements ExprVisitor.
func (p *AstPrinter) VisitGetExpr(expr *Get) (interface{}, error) {
	return p.parenthesize("get", expr.object, expr.name)
}

func (p *AstPrinter) VisitBinaryExpr(expr *Binary) (interface{}, error) {
	return p.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (p *AstPrinter) VisitGroupingExpr(expr *Grouping) (interface{}, error) {
	return p.parenthesize("group", expr.expression)
}

func (p *AstPrinter) VisitLiteralExpr(expr *Literal) (interface{}, error) {
	switch expr.value.(type) {
	case float64:
		// This is here as go prints a 1.0 float as 1 and we want a minimum
		// of one decimal place to pass the tests
		// if this is no longer a requirement, we can remove this check
		if expr.value.(float64) == float64(int(expr.value.(float64))) {
			return fmt.Sprintf("%.1f", expr.value), nil
		} else {
			return fmt.Sprintf("%g", expr.value), nil
		}
	case nil:
		return "nil", nil
	default:
		return fmt.Sprint(expr.value), nil
	}
}

func (p *AstPrinter) VisitUnaryExpr(expr *Unary) (interface{}, error) {
	return p.parenthesize(expr.operator.lexeme, expr.right)
}

func (p *AstPrinter) VisitVariableExpr(expr *Variable) (interface{}, error) {
	return expr.name.lexeme, nil
}

func (p *AstPrinter) VisitAssignExpr(expr *Assign) (interface{}, error) {
	return p.parenthesize(expr.name.lexeme+" = ", expr.value)
}

func (p *AstPrinter) VisitLogicalExpr(expr *Logical) (interface{}, error) {
	return p.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (p *AstPrinter) VisitCallExpr(expr *Call) (interface{}, error) {
	return p.parenthesize("call", append([]Expr{expr.callee}, expr.arguments...)...)
}

func (p *AstPrinter) parenthesize(name string, exprs ...Expr) (interface{}, error) {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteString(" ")
		result, err := expr.Accept(p)
		if err != nil {
			return "", err
		}
		sb.WriteString(result.(string))
	}
	sb.WriteString(")")
	return sb.String(), nil
}

func (p *AstPrinter) print(expr Expr) string {
	result, err := expr.Accept(p)
	if err != nil {
		return ""
	}
	return result.(string)
}
