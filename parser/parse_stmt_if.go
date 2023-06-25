package parser

import (
	"tiny-go/ast"
	"tiny-go/token"
)

func (p *Parser) parseStmtIf() *ast.IfStmt {
	tokIf := p.MustAcceptToken(token.IF)

	ifStmt := &ast.IfStmt{
		If: tokIf.Pos,
	}

	stmt := p.parseStmt()
	if _, ok := p.AcceptToken(token.SEMICOLON); ok {
		ifStmt.Init = stmt
		ifStmt.Cond = p.parseExpr()
		ifStmt.Body = p.parseStmtBlock()
	} else {
		ifStmt.Init = nil
		if cond, ok := stmt.(*ast.ExprStmt); ok {
			ifStmt.Cond = cond.X
		} else {
			p.errorf(tokIf.Pos, "if cond expect expr: %#v", stmt)
		}
		ifStmt.Body = p.parseStmtBlock()
	}

	if _, ok := p.AcceptToken(token.ELSE); ok {
		switch p.PeekToken().Type {
		case token.IF: // else if
			ifStmt.Else = p.parseStmtIf()
		default:
			ifStmt.Else = p.parseStmtBlock()
		}
	}

	return ifStmt
}
