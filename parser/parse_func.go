package parser

import (
	"tiny-go/ast"
	"tiny-go/token"
)

func (p *Parser) parseFunc() *ast.FuncDecl {
	// func main()
	tokFunc := p.MustAcceptToken(token.FUNC)
	tokFuncIdent := p.MustAcceptToken(token.IDENT)

	fn := &ast.FuncDecl{
		FuncPos: tokFunc.Pos,
		NamePos: tokFuncIdent.Pos,
		Name:    tokFuncIdent.Literal,
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
		},
	}
	// parse params
	p.MustAcceptToken(token.LPAREN) // (
	for {
		// )
		if _, ok := p.AcceptToken(token.RPAREN); ok {
			break
		}

		// arg type, ...
		tokArg := p.MustAcceptToken(token.IDENT)
		tokTyp := p.MustAcceptToken(token.IDENT)

		fn.Type.Params.List = append(fn.Type.Params.List, &ast.Field{
			Name: &ast.Ident{
				NamePos: tokArg.Pos,
				Name:    tokArg.Literal,
				Type:    LoopUp(tokTyp.Literal),
			},
			Type: &ast.Ident{
				NamePos: tokTyp.Pos,
				Name:    tokTyp.Literal,
				Type:    LoopUp(tokTyp.Literal),
			},
		})
	}

	// result type
	if _, ok := p.AcceptToken(token.LBRACE, token.SEMICOLON); ok {
		p.UnreadToken()
	} else {
		tok := p.MustAcceptToken(token.IDENT)
		fn.Type.Result = &ast.Ident{
			NamePos: tok.Pos,
			Name:    tok.Literal,
			Type:    LoopUp(tok.Literal),
		}
	}

	// body: {}
	if _, ok := p.AcceptToken(token.LBRACE); ok {
		p.UnreadToken()
		fn.Body = p.parseStmtBlock()
	}

	return fn
}
