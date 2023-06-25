package parser

import (
	"fmt"
	"tiny-go/ast"
	"tiny-go/lexer"
	"tiny-go/token"
)

type Parser struct {
	fileName string
	src      string

	*TokenStream
	file *ast.File
	err  error
}

func (p *Parser) errorf(pos token.Pos, format string, args ...interface{}) {
	p.err = fmt.Errorf("%s: %s", lexer.PosString(p.fileName, p.src, int(pos)), fmt.Sprintf(format, args...))
	panic(p.err)
}

func (p *Parser) ParseFile() (file *ast.File, err error) {
	defer func() {
		if r := recover(); r != p.err {
			panic(r)
		}
		file, err = p.file, p.err
	}()

	tokens, comments := lexer.Lex(p.fileName, p.src)
	for _, tok := range tokens {
		if tok.Type == token.ERROR {
			p.errorf(tok.Pos, "invalid token: %s", tok.Literal)
		}
	}

	p.TokenStream = NewTokenStream(p.fileName, p.src, tokens, comments)
	p.parseFile()
	return
}

func NewParser(fileName, src string) *Parser {
	return &Parser{
		fileName: fileName,
		src:      src,
	}
}

func ParseFile(fileName, src string) (*ast.File, error) {
	p := NewParser(fileName, src)
	return p.ParseFile()
}

var keywords = map[string]string{
	"int":   "i32",
	"float": "float",
	"char":  "i8",
}

func LoopUp(ident string) string {
	if typ, isKeyword := keywords[ident]; isKeyword {
		return typ
	}
	return "ident"
}
