package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func (p *AstPrinter) VisitBlockStmt(stmt *Block) (interface{}, error) {
	var sb strings.Builder
	sb.WriteString("(block ")

	for _, statement := range stmt.statements {
		result, err := statement.Accept(p)
		if err != nil {
			return "", err
		}
		sb.WriteString(result.(string))
	}

	sb.WriteString(")")
	return sb.String(), nil
}

func (p *AstPrinter) VisitClassStmt(stmt *Class) (interface{}, error) {
	var sb strings.Builder
	sb.WriteString("(class " + stmt.name.lexeme)

	if stmt.superclass != nil {
		sb.WriteString(" < " + p.print(stmt.superclass))
	}

	for _, method := range stmt.methods {
		sb.WriteString(" " + p.printStmt(&method))
	}

	sb.WriteString(")")
	return sb.String(), nil
}

func (p *AstPrinter) VisitExpressionStmt(stmt *Expression) (interface{}, error) {
	return p.parenthesize(";", stmt.expression)
}

func (p *AstPrinter) VisitFunctionStmt(stmt *Function) (interface{}, error) {
	var sb strings.Builder
	sb.WriteString("(fun " + stmt.name.lexeme + "(")

	for i, param := range stmt.params {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(param.lexeme)
	}

	sb.WriteString(") ")

	for _, statement := range stmt.body {
		sb.WriteString(p.printStmt(statement))
	}

	sb.WriteString(")")
	return sb.String(), nil
}

func (p *AstPrinter) VisitIfStmt(stmt *If) (interface{}, error) {
	if stmt.elseBranch == nil {
		return p.parenthesize2("if", stmt.condition, stmt.thenBranch)
	}
	return p.parenthesize2("if-else", stmt.condition, stmt.thenBranch, stmt.elseBranch)
}

func (p *AstPrinter) VisitPrintStmt(stmt *Print) (interface{}, error) {
	return p.parenthesize("print", stmt.expression)
}

func (p *AstPrinter) VisitReturnStmt(stmt *Return) (interface{}, error) {
	if stmt.value == nil {
		return "(return)", nil
	}
	return p.parenthesize("return", stmt.value)
}

func (p *AstPrinter) VisitVarStmt(stmt *Var) (interface{}, error) {
	if stmt.initializer == nil {
		return p.parenthesize2("var", stmt.name)
	}
	return p.parenthesize2("var", stmt.name, "=", stmt.initializer)
}

func (p *AstPrinter) VisitWhileStmt(stmt *While) (interface{}, error) {
	return p.parenthesize2("while", stmt.condition, stmt.body)
}

func (p *AstPrinter) VisitSetExpr(expr *Set) (interface{}, error) {
	return p.parenthesize2("=", expr.object, expr.name.lexeme, expr.value)
}

func (p *AstPrinter) VisitSuperExpr(expr *Super) (interface{}, error) {
	return p.parenthesize2("super", expr.method)
}

func (p *AstPrinter) VisitThisExpr(expr *This) (interface{}, error) {
	return "this", nil
}

func (p *AstPrinter) VisitGetExpr(expr *Get) (interface{}, error) {
	return p.parenthesize2(".", expr.object, expr.name.lexeme)
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
	return p.parenthesize2("=", expr.name.lexeme, expr.value)
}

func (p *AstPrinter) VisitLogicalExpr(expr *Logical) (interface{}, error) {
	return p.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (p *AstPrinter) VisitCallExpr(expr *Call) (interface{}, error) {
	return p.parenthesize2("call", expr.callee, expr.arguments)
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

func (p *AstPrinter) parenthesize2(name string, parts ...interface{}) (string, error) {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(name)
	if err := p.transform(&sb, parts...); err != nil {
		return "", err
	}
	sb.WriteString(")")
	return sb.String(), nil
}

func (p *AstPrinter) transform(sb *strings.Builder, parts ...interface{}) error {
	for _, part := range parts {
		sb.WriteString(" ")
		switch part := part.(type) {
		case Token:
			sb.WriteString(part.lexeme)
		case Expr:
			result, err := part.Accept(p)
			if err != nil {
				return err
			}
			sb.WriteString(result.(string))
		case Stmt:
			result, err := part.Accept(p)
			if err != nil {
				return err
			}
			sb.WriteString(result.(string))
		case []interface{}:
			if err := p.transform(sb, part...); err != nil {
				return err
			}
		default:
			sb.WriteString(fmt.Sprint(part))
		}
	}
	return nil
}

func (p *AstPrinter) print(expr Expr) string {
	result, err := expr.Accept(p)
	if err != nil {
		return ""
	}
	return result.(string)
}

func (p *AstPrinter) printStmt(stmt Stmt) string {
	result, err := stmt.Accept(p)
	if err != nil {
		return ""
	}
	return result.(string)
}
