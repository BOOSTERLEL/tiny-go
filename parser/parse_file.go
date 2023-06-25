package parser

import (
	"strconv"
	"tiny-go/ast"
	"tiny-go/token"
)

func (p *Parser) parseFile() {
	p.file = &ast.File{}

	// package xxx
	p.file.Pkg = p.parsePackage()
	p.file.FileName = p.fileName
	p.file.Source = p.src

	for {
		switch tok := p.PeekToken(); tok.Type {
		case token.EOF:
			return
		case token.ERROR:
			panic(tok)
		case token.SEMICOLON:
			p.AcceptTokenList(token.SEMICOLON)
		case token.IMPORT:
			p.file.Imports = append(p.file.Imports, p.parseImport()...)
		case token.VAR:
			p.file.Globals = append(p.file.Globals, p.parseStmtVar())
		case token.FUNC:
			p.file.Funcs = append(p.file.Funcs, p.parseFunc())
		default:
			p.errorf(tok.Pos, "unknown token: %v", tok)
		}
	}
}

func (p *Parser) parsePackage() *ast.PackageSpec {
	tokPkg := p.MustAcceptToken(token.PACKAGE)
	tokPkgIdent := p.MustAcceptToken(token.IDENT)

	return &ast.PackageSpec{
		PkgPos:  tokPkg.Pos,
		NamePos: tokPkgIdent.Pos,
		Name:    tokPkgIdent.Literal,
	}
}

// parseImport parse:
// import "path/to/pkg"
// import name "path/to/pkg"
// import (...)
func (p *Parser) parseImport() []*ast.ImportSpec {
	tokImport := p.MustAcceptToken(token.IMPORT)
	var importSpecList []*ast.ImportSpec

	_, flag := p.AcceptToken(token.LPAREN)
	for {
		var importSpec = &ast.ImportSpec{
			ImportPos: tokImport.Pos,
		}
		asName, ok := p.AcceptToken(token.IDENT)
		if ok {
			importSpec.Name = &ast.Ident{
				NamePos: asName.Pos,
				Name:    asName.Literal,
			}
		}

		if pkgPath, ok := p.AcceptToken(token.STRING); ok {
			path, _ := strconv.Unquote(pkgPath.Literal)
			importSpec.Path = path
		}
		importSpecList = append(importSpecList, importSpec)

		if _, ok := p.AcceptToken(token.RPAREN); ok || !flag {
			break
		}
	}

	return importSpecList
}
