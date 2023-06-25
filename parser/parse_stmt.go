package parser

import (
	"tiny-go/ast"
	"tiny-go/token"
)

func (p *Parser) parseStmt() ast.Stmt {
	switch tok := p.PeekToken(); tok.Type {
	case token.EOF:
		return nil
	case token.ERROR:
		p.errorf(tok.Pos, "invalid token: %s", tok.Literal)
	case token.SEMICOLON:
		p.AcceptTokenList(token.SEMICOLON)
		return nil
	case token.LBRACE: // {
		return p.parseStmtBlock()
	case token.RETURN:
		return p.parseStmtReturn()
	case token.VAR:
		return p.parseStmtVar()
	default:
		return p.parseStmtExprOrAssign()
	}
	panic("unreachable")
}

func (p *Parser) parseStmtExprOrAssign() ast.Stmt {
	// exprList ;
	// exprList := exprList;
	// exprList = exprList;
	exprList := p.parseExprList()
	switch tok := p.PeekToken(); tok.Type {
	case token.SEMICOLON, token.LBRACE:
		if len(exprList) != 1 {
			p.errorf(tok.Pos, "unknown token: %v", tok.Type)
		}
		return &ast.ExprStmt{
			X: exprList[0],
		}
	case token.DEFINE, token.ASSIGN:
		p.ReadToken()
		exprValueList := p.parseExprList()
		if len(exprList) != len(exprValueList) {
			p.errorf(tok.Pos, "unknown token: %v", tok)
		}
		var assignStmt = &ast.AssignStmt{
			Target: make([]*ast.Ident, len(exprList)),
			OpPos:  tok.Pos,
			Op:     tok.Type,
			Value:  make([]ast.Expr, len(exprList)),
		}
		for i, target := range exprList {
			assignStmt.Target[i] = target.(*ast.Ident)
			assignStmt.Value[i] = exprValueList[i]
			var typ string
			switch exprValueList[i].(type) {
			case *ast.Float:
				typ = "float"
			case *ast.Int:
				typ = "i32"
			case *ast.Char:
				typ = "i8"
			case *ast.Ident:
				typ = exprValueList[i].(*ast.Ident).Type
			case *ast.BinaryExpr:
				switch exprValueList[i].(*ast.BinaryExpr).X.(type) {
				case *ast.Float:
					typ = "float"
				case *ast.Int:
					typ = "i32"
				case *ast.Char:
					typ = "i8"
				case *ast.Ident:
					typ = exprValueList[i].(*ast.BinaryExpr).X.(*ast.Ident).Type
				default:
					p.errorf(tok.Pos, "unknown token: %v", tok)
				}
			}
			assignStmt.Target[i].Type = typ
		}
		return assignStmt
	default:
		p.errorf(tok.Pos, "unknown token: %v", tok)
	}
	panic("unreachable")
}

func (p *Parser) parseStmtBlock() *ast.BlockStmt {
	block := &ast.BlockStmt{}

	tokBegin := p.MustAcceptToken(token.LBRACE) // {

Loop:
	for {
		switch tok := p.PeekToken(); tok.Type {
		case token.EOF:
			break Loop
		case token.ERROR:
			p.errorf(tok.Pos, "invalid token: %s", tok.Literal)
		case token.SEMICOLON:
			p.AcceptTokenList(token.SEMICOLON)
		case token.LBRACE: // {
			block.List = append(block.List, p.parseStmtBlock())
		case token.RBRACE: // }
			break Loop
		case token.VAR:
			block.List = append(block.List, p.parseStmtVar())
		case token.RETURN:
			block.List = append(block.List, p.parseStmtReturn())
		case token.IF:
			block.List = append(block.List, p.parseStmtIf())
		case token.FOR:
			block.List = append(block.List, p.parseStmtFor())
		case token.BREAK:
			block.List = append(block.List, p.parseStmtBreak())
		case token.CONTINUE:
			block.List = append(block.List, p.parseStmtContinue())
		case token.GOTO:
			block.List = append(block.List, p.parseStmtGoto())
		default:
			p.ReadToken()
			tok = p.PeekToken()
			p.UnreadToken()
			if tok.Type == token.COLON {
				block.List = append(block.List, p.parseStmtLabeled())
			} else {
				block.List = append(block.List, p.parseStmtExprOrAssign())
			}
		}
	}

	tokEnd := p.MustAcceptToken(token.RBRACE) // }

	block.Lbrace = tokBegin.Pos
	block.Rbrace = tokEnd.Pos

	return block
}

func (p *Parser) parseStmtExpr() *ast.ExprStmt {
	return &ast.ExprStmt{
		X: p.parseExpr(),
	}
}

func (p *Parser) parseStmtGoto() *ast.BranchStmt {
	tokGoto := p.MustAcceptToken(token.GOTO)

	return &ast.BranchStmt{
		TokPos:  tokGoto.Pos,
		TokType: token.GOTO,
	}
}

func (p *Parser) parseStmtLabeled() *ast.LabeledStmt {
	tokLabel := p.MustAcceptToken(token.IDENT)
	tokColon := p.MustAcceptToken(token.COLON)

	return &ast.LabeledStmt{
		Label: &ast.Ident{
			NamePos: tokLabel.Pos,
			Name:    tokLabel.Literal,
		},
		Colon: tokColon.Pos,
	}
}
