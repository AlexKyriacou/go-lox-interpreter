package main

type Stmt interface {
	Accept(visitor StmtVisitor) error
}

type StmtVisitor interface {
	VisitExpressionStmt(stmt *Expression) error
	VisitPrintStmt(stmt *Print) error
}

type Expression struct {
	expression Expr
}

func (e *Expression) Accept(visitor StmtVisitor) error {
	return visitor.VisitExpressionStmt(e)
}

type Print struct {
	expression Expr
}

func (p *Print) Accept(visitor StmtVisitor) error {
	return visitor.VisitPrintStmt(p)
}
