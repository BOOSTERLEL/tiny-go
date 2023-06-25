package parser

import (
	"tiny-go/ast"
	"tiny-go/token"
)

func (p *Parser) parseStmtVar() *ast.VarSpec {
	tokVar := p.MustAcceptToken(token.VAR)
	tokIdent := p.MustAcceptToken(token.IDENT)

	var varSpec = &ast.VarSpec{
		VarPos: tokVar.Pos,
	}

	varSpec.Name = &ast.Ident{
		NamePos: tokIdent.Pos,
		Name:    tokIdent.Literal,
	}

	// var name type?
	if typ, ok := p.AcceptToken(token.IDENT); ok {
		tokType := LoopUp(typ.Literal)
		varSpec.Type = &ast.Ident{
			NamePos: typ.Pos,
			Name:    typ.Literal,
			Type:    tokType,
		}
		varSpec.Name.Type = tokType
	}

	// var name =
	if _, ok := p.AcceptToken(token.ASSIGN); ok {
		varSpec.Value = p.parseExpr()
	}

	p.AcceptTokenList(token.SEMICOLON)
	return varSpec
}
