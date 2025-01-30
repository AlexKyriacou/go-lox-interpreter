package main

type Stmt interface {
	Accept(visitor StmtVisitor) (interface{}, error)
}

type StmtVisitor interface {
	VisitBlockStmt(stmt *Block) (interface{}, error)
	VisitClassStmt(stmt *Class) (interface{}, error)
	VisitExpressionStmt(stmt *Expression) (interface{}, error)
	VisitFunctionStmt(stmt *Function) (interface{}, error)
	VisitIfStmt(stmt *If) (interface{}, error)
	VisitPrintStmt(stmt *Print) (interface{}, error)
	VisitReturnStmt(stmt *Return) (interface{}, error)
	VisitVarStmt(stmt *Var) (interface{}, error)
	VisitWhileStmt(stmt *While) (interface{}, error)
}

type Block struct {
	statements []Stmt
}

func (b *Block) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitBlockStmt(b)
}

type Class struct {
	name       Token
	superclass *Variable
	methods    []Function
}

func (c *Class) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitClassStmt(c)
}

type Expression struct {
	expression Expr
}

func (e *Expression) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitExpressionStmt(e)
}

type Function struct {
	name   Token
	params []Token
	body   []Stmt
}

func (f *Function) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitFunctionStmt(f)
}

type If struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (i *If) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitIfStmt(i)
}

type Print struct {
	expression Expr
}

func (p *Print) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitPrintStmt(p)
}

type Return struct {
	keyword Token
	value   Expr
}

func (r *Return) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitReturnStmt(r)
}

type Var struct {
	name        Token
	initializer Expr
}

func (v *Var) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitVarStmt(v)
}

type While struct {
	condition Expr
	body      Stmt
}

func (w *While) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitWhileStmt(w)
}
