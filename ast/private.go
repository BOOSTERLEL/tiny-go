package ast

import "tiny-go/token"

func (b *BlockStmt) Pos() token.Pos {
	return token.NoPos
}

func (e *ExprStmt) Pos() token.Pos {
	return token.NoPos
}

func (c *CallExpr) Pos() token.Pos {
	return token.NoPos
}

func (b *BinaryExpr) Pos() token.Pos {
	return token.NoPos
}

func (n *Int) Pos() token.Pos {
	return token.NoPos
}

func (u *UnaryExpr) Pos() token.Pos {
	return token.NoPos
}

func (p *ParenExpr) Pos() token.Pos {
	return token.NoPos
}

func (i *Ident) Pos() token.Pos {
	return token.NoPos
}

func (p *File) Pos() token.Pos {
	return token.NoPos
}

func (p *PackageSpec) Pos() token.Pos {
	return token.NoPos
}

func (f *FuncDecl) Pos() token.Pos {
	return token.NoPos
}

func (v *VarSpec) Pos() token.Pos {
	return token.NoPos
}

func (a *AssignStmt) Pos() token.Pos {
	return token.NoPos
}

func (i *IfStmt) Pos() token.Pos {
	return token.NoPos
}

func (f *ForStmt) Pos() token.Pos {
	return token.NoPos
}

func (r *ReturnStmt) Pos() token.Pos {
	return token.NoPos
}

func (i *ImportSpec) Pos() token.Pos {
	return token.NoPos
}

func (s *SelectorExpr) Pos() token.Pos {
	return token.NoPos
}
func (b BranchStmt) Pos() token.Pos {
	return token.NoPos
}

func (f *Float) Pos() token.Pos {
	return token.NoPos
}

func (l LabeledStmt) Pos() token.Pos {
	return token.NoPos
}

func (c *Char) Pos() token.Pos {
	return token.NoPos
}

func (b *BlockStmt) End() token.Pos {
	return token.NoPos
}

func (e *ExprStmt) End() token.Pos {
	return token.NoPos
}

func (c *CallExpr) End() token.Pos {
	return token.NoPos
}

func (b *BinaryExpr) End() token.Pos {
	return token.NoPos
}

func (n *Int) End() token.Pos {
	return token.NoPos
}

func (u *UnaryExpr) End() token.Pos {
	return token.NoPos
}

func (p *ParenExpr) End() token.Pos {
	return token.NoPos
}

func (i *Ident) End() token.Pos {
	return token.NoPos
}

func (p *File) End() token.Pos {
	return token.NoPos
}

func (p *PackageSpec) End() token.Pos {
	return token.NoPos
}

func (f *FuncDecl) End() token.Pos {
	return token.NoPos
}

func (v *VarSpec) End() token.Pos {
	return token.NoPos
}

func (a *AssignStmt) End() token.Pos {
	return token.NoPos
}

func (i *IfStmt) End() token.Pos {
	return token.NoPos
}

func (f *ForStmt) End() token.Pos {
	return token.NoPos
}

func (r *ReturnStmt) End() token.Pos {
	return token.NoPos
}

func (i *ImportSpec) End() token.Pos {
	return token.NoPos
}

func (s *SelectorExpr) End() token.Pos {
	return token.NoPos
}

func (f *Float) End() token.Pos {
	return token.NoPos
}

func (b BranchStmt) End() token.Pos {
	return token.NoPos
}

func (l LabeledStmt) End() token.Pos {
	return token.NoPos
}

func (c *Char) End() token.Pos {
	return token.NoPos
}

func (b *BlockStmt) stmtType() {

}

func (e *ExprStmt) stmtType() {

}

func (e *ExprStmt) nodeType() {

}

func (c *CallExpr) exprType() {

}

func (b *BinaryExpr) exprType() {

}

func (n *Int) exprType() {

}

func (u *UnaryExpr) exprType() {

}

func (p *ParenExpr) exprType() {

}

func (i *Ident) exprType() {

}

func (i *Ident) nodeType() {

}

func (p *File) nodeType() {

}

func (p *PackageSpec) nodeType() {

}

func (f *FuncDecl) nodeType() {

}

func (v *VarSpec) stmtType() {

}

func (v *VarSpec) nodeType() {

}

func (a *AssignStmt) stmtType() {

}

func (i *IfStmt) stmtType() {

}

func (f *ForStmt) stmtType() {

}

func (r *ReturnStmt) stmtType() {

}

func (i *ImportSpec) nodeType() {

}

func (s *SelectorExpr) exprType() {

}

func (f *Float) exprType() {

}

func (b BranchStmt) stmtType() {

}

func (l LabeledStmt) stmtType() {

}

func (c *Char) exprType() {

}
