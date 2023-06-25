package parser

import (
	"tiny-go/ast"
	"tiny-go/token"
)

func (p *Parser) parseStmtFor() *ast.ForStmt {
	tokFor := p.MustAcceptToken(token.FOR)

	forStmt := &ast.ForStmt{
		For: tokFor.Pos,
	}

	// for {}
	if _, ok := p.AcceptToken(token.LBRACE); ok {
		p.UnreadToken()
		forStmt.Body = p.parseStmtBlock()
		return forStmt
	}

	// for Cond {}
	// for Init?; Cond?; Post? {}

	// for ; ...
	if _, ok := p.AcceptToken(token.SEMICOLON); ok {
		forStmt.Init = nil

		// for ;; ...
		if _, ok := p.AcceptToken(token.SEMICOLON); ok {
			if _, ok := p.AcceptToken(token.LBRACE); ok {
				// for ;; {}
				p.UnreadToken()
				forStmt.Body = p.parseStmtBlock()
				return forStmt
			} else {
				// for ;; postStmt {}
				forStmt.Post = p.parseStmt()
				forStmt.Body = p.parseStmtBlock()
				return forStmt
			}
		} else {
			// for ; cond ; ... {}
			forStmt.Cond = p.parseExpr()
			p.MustAcceptToken(token.SEMICOLON)
			if _, ok := p.AcceptToken(token.LBRACE); ok {
				// for ; cond ; {}
				p.UnreadToken()
				forStmt.Body = p.parseStmtBlock()
				return forStmt
			} else {
				// for ; cond ; postStmt {}
				forStmt.Post = p.parseStmt()
				forStmt.Body = p.parseStmtBlock()
				return forStmt
			}
		}
	} else {
		// for expr ... {}
		stmt := p.parseStmt()

		if _, ok := p.AcceptToken(token.LBRACE); ok {
			// for cond {}
			p.UnreadToken()
			if expr, ok := stmt.(ast.Expr); ok {
				forStmt.Cond = expr
			}
			forStmt.Body = p.parseStmtBlock()
			return forStmt
		} else {
			// for init ...
			p.MustAcceptToken(token.SEMICOLON)
			forStmt.Init = stmt

			// for ;; ...
			if _, ok := p.AcceptToken(token.SEMICOLON); ok {
				if _, ok := p.AcceptToken(token.LBRACE); ok {
					// for ;; {}
					p.UnreadToken()
					forStmt.Body = p.parseStmtBlock()
					return forStmt
				} else {
					// for ;; postStmt {}
					forStmt.Post = p.parseStmt()
					forStmt.Body = p.parseStmtBlock()
					return forStmt
				}
			} else {
				// for ; cond ; ... {}
				forStmt.Cond = p.parseExpr()
				p.MustAcceptToken(token.SEMICOLON)
				if _, ok := p.AcceptToken(token.LBRACE); ok {
					// for ; cond ; {}
					p.UnreadToken()
					forStmt.Body = p.parseStmtBlock()
					return forStmt
				} else {
					// for ; cond ; postStmt {}
					forStmt.Post = p.parseStmt()
					forStmt.Body = p.parseStmtBlock()
					return forStmt
				}
			}
		}
	}
}

func (p *Parser) parseStmtBreak() *ast.BranchStmt {
	tokBreak := p.MustAcceptToken(token.BREAK)

	return &ast.BranchStmt{
		TokPos:  tokBreak.Pos,
		TokType: token.BREAK,
	}
}

func (p *Parser) parseStmtContinue() *ast.BranchStmt {
	tokContinue := p.MustAcceptToken(token.CONTINUE)

	return &ast.BranchStmt{
		TokPos:  tokContinue.Pos,
		TokType: token.CONTINUE,
	}
}
