package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func (p *AstPrinter) VisitBinaryExpr(expr *Binary) interface{} {
	return p.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (p *AstPrinter) VisitGroupingExpr(expr *Grouping) interface{} {
	return p.parenthesize("group", expr.expression)
}

func (p *AstPrinter) VisitLiteralExpr(expr *Literal) interface{} {
	return fmt.Sprint(expr.value)
}

func (p *AstPrinter) VisitUnaryExpr(expr *Unary) interface{} {
	return p.parenthesize(expr.operator.lexeme, expr.right)
}

func (p *AstPrinter) parenthesize(name string, exprs ...Expr) interface{} {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(name)
	for _, expr := range exprs{
		sb.WriteString(" ")
		sb.WriteString(expr.Accept(p).(string))
	}
	sb.WriteString(")")
	return sb.String()
}

func (p *AstPrinter) print(expr Expr) string {
	return expr.Accept(p).(string)
}