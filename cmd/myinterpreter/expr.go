package main

type Expr interface {
	Accept(visitor ExprVisitor) (interface{}, error)
}

type ExprVisitor interface {
	VisitBinaryExpr(expr *Binary) (interface{}, error)
	VisitGroupingExpr(expr *Grouping) (interface{}, error)
	VisitLiteralExpr(expr *Literal) (interface{}, error)
	VisitUnaryExpr(expr *Unary) (interface{}, error)
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func (b *Binary) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitBinaryExpr(b)
}

type Grouping struct {
	expression Expr
}

func (g *Grouping) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitGroupingExpr(g)
}

type Literal struct {
	value interface{}
}

func (l *Literal) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitLiteralExpr(l)
}

type Unary struct {
	operator Token
	right    Expr
}

func (u *Unary) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitUnaryExpr(u)
}
