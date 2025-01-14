package main

type Stmt interface {
	Accept(visitor StmtVisitor) error
}

type StmtVisitor interface {
	VisitBlockStmt(stmt *Block) error
	VisitExpressionStmt(stmt *Expression) error
	VisitFunctionStmt(stmt *Function) error
	VisitIfStmt(stmt *If) error
	VisitPrintStmt(stmt *Print) error
	VisitReturnStmt(stmt *Return) error
	VisitVarStmt(stmt *Var) error
	VisitWhileStmt(stmt *While) error
}

type Block struct {
	statements []Stmt
}

func (b *Block) Accept(visitor StmtVisitor) error {
	return visitor.VisitBlockStmt(b)
}

type Expression struct {
	expression Expr
}

func (e *Expression) Accept(visitor StmtVisitor) error {
	return visitor.VisitExpressionStmt(e)
}

type Function struct {
	name   Token
	params []Token
	body   []Stmt
}

func (f *Function) Accept(visitor StmtVisitor) error {
	return visitor.VisitFunctionStmt(f)
}

type If struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (i *If) Accept(visitor StmtVisitor) error {
	return visitor.VisitIfStmt(i)
}

type Print struct {
	expression Expr
}

func (p *Print) Accept(visitor StmtVisitor) error {
	return visitor.VisitPrintStmt(p)
}

type Return struct {
	keyword Token
	value   Expr
}

func (r *Return) Accept(visitor StmtVisitor) error {
	return visitor.VisitReturnStmt(r)
}

type Var struct {
	name        Token
	initializer Expr
}

func (v *Var) Accept(visitor StmtVisitor) error {
	return visitor.VisitVarStmt(v)
}

type While struct {
	condition Expr
	body      Stmt
}

func (w *While) Accept(visitor StmtVisitor) error {
	return visitor.VisitWhileStmt(w)
}
