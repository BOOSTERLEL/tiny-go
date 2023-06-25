package lexer

import (
	"fmt"
	gotoken "go/token"
	"tiny-go/token"
)

func PosString(fileName string, src string, pos int) string {
	fSet := gotoken.NewFileSet()
	fSet.AddFile(fileName, 1, len(src)).SetLinesForContent([]byte(src))
	return fmt.Sprintf("%v", fSet.Position(gotoken.Pos(pos+1)))
}

type Lexer struct {
	src      *SourceStream
	tokens   []token.Token
	comments []token.Token
}

func NewLexer(name, input string) *Lexer {
	p := &Lexer{src: NewSourceStream(name, input)}
	p.run()
	return p
}

func (p *Lexer) Tokens() []token.Token {
	return p.tokens
}

func (p *Lexer) Comments() []token.Token {
	return p.comments
}

func (p *Lexer) emit(typ token.TokenType) {
	lit, pos := p.src.EmitToken()
	if typ == token.IDENT {
		typ = token.LoopUp(lit)
	}
	p.tokens = append(p.tokens, token.Token{
		Type:    typ,
		Literal: lit,
		Pos:     token.Pos(pos + 1),
	})
}

func (p *Lexer) emitComment() {
	lit, pos := p.src.EmitToken()
	p.comments = append(p.comments, token.Token{
		Type:    token.COMMENT,
		Literal: lit,
		Pos:     token.Pos(pos + 1),
	})
}

func (p *Lexer) errorf(format string, args ...interface{}) {
	tok := token.Token{
		Type:    token.ERROR,
		Literal: fmt.Sprintf(format, args...),
		Pos:     token.Pos(p.src.pos),
	}
	p.tokens = append(p.tokens, tok)
	panic(tok)
}

func (p *Lexer) run() (tokens []token.Token) {
	defer func() {
		tokens = p.tokens
		if r := recover(); r != nil {
			if _, ok := r.(token.Token); !ok {
				panic(r)
			}
		}
	}()

	for {
		r := p.src.Read()
		if r == rune(token.EOF) {
			p.emit(token.EOF)
			return
		}

		switch {
		case r == '\n':
			p.src.IgnoreToken()
			if len(p.tokens) > 0 {
				switch p.tokens[len(p.tokens)-1].Type {
				case token.RPAREN, token.IDENT, token.INT, token.RETURN, token.FLOAT:
					p.emit(token.SEMICOLON)
				}
			}
		case isSpace(r):
			p.src.IgnoreToken()
		case isAlpha(r):
			p.src.Unread()
			for {
				if r := p.src.Read(); !isAlphaNumberic(r) {
					p.src.Unread()
					p.emit(token.IDENT)
					break
				}
			}
		case '0' <= r && r <= '9': // 123, 1.0
			p.src.Unread()

			digits := "0123456789"
			typ := p.src.AcceptRun(digits)
			p.emit(typ)
		case r == '+': // +, +=, ++
			p.emit(token.ADD)
		case r == '-': // -, -=, --
			p.emit(token.SUB)
		case r == '*': // *, *=
			p.emit(token.MUL)
		case r == '/': // /, //, /*, /=
			peek := p.src.Peek()
			if peek == '/' {
				// line comment
				for {
					t := p.src.Read()
					if t == '\n' {
						p.src.Unread()
						p.emitComment()
						break
					}
					if t == rune(token.EOF) {
						p.emitComment()
						return
					}
				}
			} else if peek == '*' {
				// multiline comment
				for {
					t := p.src.Read()
					if t == '*' && p.src.Peek() == '/' {
						p.src.Read()
						p.emitComment()
						break
					}
					if t == rune(token.EOF) {
						p.errorf("unterminated quoted string")
						return
					}
				}
			} else {
				p.emit(token.DIV)
			}
		case r == '%': // %
			p.emit(token.MOD)
		case r == '=': // =,==
			switch p.src.Read() {
			case '=':
				p.emit(token.EQL)
			default:
				p.src.Unread()
				p.emit(token.ASSIGN)
			}
		case r == '!': // !=
			switch p.src.Read() {
			case '=':
				p.emit(token.NEQ)
			default:
				p.src.Unread()
				//p.errorf("unrecognized character: %#u", r)
				p.emit(token.NOT)
			}
		case r == '<': // <,<=
			switch p.src.Read() {
			case '=':
				p.emit(token.LEQ)
			default:
				p.src.Unread()
				p.emit(token.LSS)
			}
		case r == '>': // >,>=
			switch p.src.Read() {
			case '=':
				p.emit(token.GEQ)
			default:
				p.src.Unread()
				p.emit(token.GTR)
			}
		case r == ':': // :,:=
			switch p.src.Read() {
			case '=':
				p.emit(token.DEFINE)
			default:
				//p.errorf("unrecognized character: %#u", r)
				p.emit(token.COLON)
			}
		case r == '&': // &&
			switch p.src.Read() {
			case '&':
				p.emit(token.AND)
			default:
				p.errorf("unrecognized character: %#u", r)
			}
		case r == '|':
			switch p.src.Read() {
			case '|':
				p.emit(token.OR)
			default:
				p.errorf("unrecognized character: %#u", r)
			}
		case r == '"':
			p.lexQuote()
		case r == '\'':
			if p.src.Read() == '\\' {
				p.src.Read()
			}
			p.src.Read()
			p.emit(token.CHAR)
		case r == '.':
			p.emit(token.PERIOD)
		case r == '(':
			p.emit(token.LPAREN)
		case r == '[':
			p.emit(token.LBRACK)
		case r == '{':
			p.emit(token.LBRACE)
		case r == ')':
			p.emit(token.RPAREN)
		case r == ']':
			p.emit(token.RBRACK)
		case r == '}':
			p.emit(token.RBRACE)
		case r == ',':
			p.emit(token.COMMA)
		case r == ';':
			p.emit(token.SEMICOLON)
		default:
			p.errorf("unrecognized character: %#U", r)
			return
		}
	}
}

func (p *Lexer) lexQuote() {
	for {
		switch p.src.Read() {
		case rune(token.EOF):
			p.errorf("unterminated quoted string")
			return
		case '\\':
			p.src.Read()
		case '"':
			p.emit(token.STRING)
			return
		}
	}
}

func Lex(name, input string) (tokens, comments []token.Token) {
	l := NewLexer(name, input)
	tokens = l.Tokens()
	comments = l.Comments()
	return
}
