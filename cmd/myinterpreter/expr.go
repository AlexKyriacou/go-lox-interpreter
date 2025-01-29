package main

type Expr interface {
	Accept(visitor ExprVisitor) (interface{}, error)
}

type ExprVisitor interface {
	VisitAssignExpr(expr *Assign) (interface{}, error)
	VisitBinaryExpr(expr *Binary) (interface{}, error)
	VisitCallExpr(expr *Call) (interface{}, error)
	VisitGetExpr(expr *Get) (interface{}, error)
	VisitGroupingExpr(expr *Grouping) (interface{}, error)
	VisitLiteralExpr(expr *Literal) (interface{}, error)
	VisitLogicalExpr(expr *Logical) (interface{}, error)
	VisitSetExpr(expr *Set) (interface{}, error)
	VisitUnaryExpr(expr *Unary) (interface{}, error)
	VisitVariableExpr(expr *Variable) (interface{}, error)
}

type Assign struct {
	name  Token
	value Expr
}

func (a *Assign) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitAssignExpr(a)
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func (b *Binary) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitBinaryExpr(b)
}

type Call struct {
	callee    Expr
	paren     Token
	arguments []Expr
}

func (c *Call) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitCallExpr(c)
}

type Get struct {
	object Expr
	name   Token
}

func (g *Get) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitGetExpr(g)
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

type Logical struct {
	left     Expr
	operator Token
	right    Expr
}

func (l *Logical) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitLogicalExpr(l)
}

type Set struct {
	object Expr
	name   Token
	value  Expr
}

func (s *Set) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitSetExpr(s)
}

type Unary struct {
	operator Token
	right    Expr
}

func (u *Unary) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitUnaryExpr(u)
}

type Variable struct {
	name Token
}

func (v *Variable) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitVariableExpr(v)
}
