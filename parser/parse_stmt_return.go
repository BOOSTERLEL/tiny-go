package parser

import (
	"tiny-go/ast"
	"tiny-go/token"
)

func (p *Parser) parseStmtReturn() *ast.ReturnStmt {
	tokReturn := p.MustAcceptToken(token.RETURN)

	retStmt := &ast.ReturnStmt{
		Return: tokReturn.Pos,
	}
	if _, ok := p.AcceptToken(
		token.SEMICOLON, // ;
		token.LBRACE,    // {
		token.RBRACE,    // }
	); !ok {
		retStmt.Result = p.parseExpr()
	} else {
		p.UnreadToken()
	}
	return retStmt
}
