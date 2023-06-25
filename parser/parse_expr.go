package parser

import (
	"strconv"
	"tiny-go/ast"
	"tiny-go/token"
)

func (p *Parser) parseExpr() ast.Expr {
	return p.parseExprBinary(1)
}

// parseExprList x, y :=
func (p *Parser) parseExprList() (exprs []ast.Expr) {
	for {
		exprs = append(exprs, p.parseExpr())
		if p.PeekToken().Type != token.COMMA {
			break
		}
		p.ReadToken()
	}
	return
}

func (p *Parser) parseExprBinary(prec int) ast.Expr {
	x := p.parseExprUnary()
	for {
		op := p.PeekToken()
		if op.Type.Precedence() < prec {
			return x
		}
		p.MustAcceptToken(op.Type)
		y := p.parseExprBinary(op.Type.Precedence() + 1)
		x = &ast.BinaryExpr{
			OpPos: op.Pos,
			Op:    op.Type,
			X:     x,
			Y:     y,
		}
	}
	return nil
}

func (p *Parser) parseExprUnary() ast.Expr {
	if _, ok := p.AcceptToken(token.ADD); ok {
		return p.parseExprPrimary()
	}
	if tok, ok := p.AcceptToken(token.SUB, token.NOT); ok {
		return &ast.UnaryExpr{
			OpPos: tok.Pos,
			Op:    tok.Type,
			X:     p.parseExprPrimary(),
		}
	}
	return p.parseExprPrimary()
}

func (p *Parser) parseExprPrimary() ast.Expr {
	if _, ok := p.AcceptToken(token.LPAREN); ok {
		expr := p.parseExpr()
		p.MustAcceptToken(token.RPAREN)
		return expr
	}

	switch tok := p.PeekToken(); tok.Type {
	case token.IDENT: // call
		p.ReadToken()
		nextTok := p.PeekToken()
		p.UnreadToken()

		switch nextTok.Type {
		case token.LPAREN:
			return p.parseExprCall()
		case token.PERIOD:
			return p.parseExprSelector()
		default:
			p.MustAcceptToken(token.IDENT)
			return &ast.Ident{
				NamePos: tok.Pos,
				Name:    tok.Literal,
			}
		}
	case token.INT:
		tokInt := p.MustAcceptToken(token.INT)
		value, _ := strconv.Atoi(tokInt.Literal)
		return &ast.Int{
			ValuePos: tokInt.Pos,
			ValueEnd: tokInt.Pos + token.Pos(len(tokInt.Literal)),
			Value:    value,
		}
	case token.FLOAT:
		tokFloat := p.MustAcceptToken(token.FLOAT)
		value, _ := strconv.ParseFloat(tokFloat.Literal, 64)
		return &ast.Float{
			ValuePos: tokFloat.Pos,
			ValueEnd: tokFloat.Pos + token.Pos(len(tokFloat.Literal)),
			Value:    value,
		}
	case token.CHAR:
		tokChar := p.MustAcceptToken(token.CHAR)
		var value int
		if len(tokChar.Literal) == 3 {
			value = int(tokChar.Literal[1])
		} else {
			value, _ =strconv.Atoi(tokChar.Literal[1:2])
		}
		return &ast.Char{
			ValuePos: tokChar.Pos,
			ValueEnd: tokChar.Pos+token.Pos(len(tokChar.Literal)),
			Value: value,
		}
	default:
		p.errorf(tok.Pos, "unknown tok: type=%v, lit=%q", tok.Type, tok.Literal)
		panic("unreachable")
	}
}

func (p *Parser) parseExprCall() *ast.CallExpr {
	tokIdent := p.MustAcceptToken(token.IDENT)
	tokLparen := p.MustAcceptToken(token.LPAREN)
	arg0 := p.parseExpr()
	tokRparen := p.MustAcceptToken(token.RPAREN)

	return &ast.CallExpr{
		FuncName: &ast.Ident{NamePos: tokIdent.Pos, Name: tokIdent.Literal},
		Lparen:   tokLparen.Pos,
		Args:     []ast.Expr{arg0},
		Rparen:   tokRparen.Pos,
	}
}

func (p *Parser) parseExprSelector() ast.Expr {
	tokX := p.MustAcceptToken(token.IDENT)
	_ = p.MustAcceptToken(token.PERIOD)
	tokSel := p.MustAcceptToken(token.IDENT)

	// pkg.fn(...)
	if nextTok := p.PeekToken(); nextTok.Type == token.LPAREN {
		var arg0 ast.Expr
		tokLparen := p.MustAcceptToken(token.LPAREN)
		if tok := p.PeekToken(); tok.Type != token.RPAREN {
			arg0 = p.parseExpr()
		}
		tokRparen := p.MustAcceptToken(token.RPAREN)

		return &ast.CallExpr{
			Pkg: &ast.Ident{
				NamePos: tokX.Pos,
				Name:    tokX.Literal,
			},
			FuncName: &ast.Ident{
				NamePos: tokSel.Pos,
				Name:    tokSel.Literal,
			},
			Lparen: tokLparen.Pos,
			Args:   []ast.Expr{arg0},
			Rparen: tokRparen.Pos,
		}
	}

	return &ast.SelectorExpr{
		X: &ast.Ident{
			NamePos: tokX.Pos,
			Name:    tokX.Literal,
		},
		Sel: &ast.Ident{
			NamePos: tokSel.Pos,
			Name:    tokX.Literal,
		},
	}
}
