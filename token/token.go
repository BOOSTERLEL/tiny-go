package token

import (
	"fmt"
	"strconv"
)

// TokenType 词法记号类型
type TokenType int

// Token 记号值
type Token struct {
	Type    TokenType // 记号的类型
	Pos     Pos       // 记号所在的位置(从1开始)
	Literal string    // 程序中原始的字符串
}

// 记号类型
const (
	EOF TokenType = iota
	ERROR
	COMMENT

	IDENT
	INT
	FLOAT
	CHAR
	STRING // "abc"

	PACKAGE
	IMPORT
	VAR
	FUNC
	RETURN
	IF
	ELSE
	FOR
	BREAK
	CONTINUE
	DEFER
	GOTO

	ADD // +
	SUB // -
	MUL // *
	DIV // /
	MOD // %

	EQL // ==
	NEQ // !=
	LSS // <
	LEQ // <=
	GTR // >
	GEQ // >=

	AND // &&
	OR  // ||
	NOT // !

	ASSIGN // =
	DEFINE // :=

	LPAREN // (
	RPAREN // )
	LBRACK // [
	RBRACK // ]
	LBRACE // {
	RBRACE // }

	COMMA     // ,
	SEMICOLON // ;
	PERIOD    // .
	COLON     // :
)

func (op TokenType) Precedence() int {
	switch op {
	case OR:
		return 1
	case AND:
		return 2
	case EQL, NEQ, LSS, LEQ, GTR, GEQ:
		return 3
	case ADD, SUB:
		return 4
	case MUL, DIV, MOD:
		return 5
	}
	return 0
}

var tokens = [...]string{
	EOF:     "EOF",
	ERROR:   "ERROR",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	CHAR:   "CHAR",
	STRING: "STRING",

	PACKAGE:  "package",
	IMPORT:   "import",
	VAR:      "var",
	FUNC:     "func",
	RETURN:   "return",
	IF:       "if",
	ELSE:     "else",
	FOR:      "for",
	BREAK:    "break",
	CONTINUE: "continue",
	DEFER:    "defer",
	GOTO:     "goto",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",
	MOD: "%",

	EQL: "==",
	NEQ: "!=",
	LSS: "<",
	LEQ: "<=",
	GTR: ">",
	GEQ: ">=",

	AND: "&&",
	OR:  "||",
	NOT: "!",

	ASSIGN: "=",
	DEFINE: ":=",

	LPAREN: "(",
	RPAREN: ")",
	LBRACK: "[",
	RBRACK: "]",
	LBRACE: "{",
	RBRACE: "}",

	COMMA:     ",",
	SEMICOLON: ";",
	PERIOD:    ".",
	COLON:     ":",
}

func (op TokenType) String() string {
	s := ""
	if 0 <= op && op < TokenType(len(tokens)) {
		s = tokens[op]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(op)) + ")"
	}
	return s
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%v : \"%v\")\t", t.Type, t.Literal)
}

func (t Token) IntValue() int {
	x, err := strconv.ParseInt(t.Literal, 10, 64)
	if err != nil {
		panic(err)
	}
	return int(x)
}

var keywords = map[string]TokenType{
	"package":  PACKAGE,
	"return":   RETURN,
	"import":   IMPORT,
	"func":     FUNC,
	"var":      VAR,
	"if":       IF,
	"else":     ELSE,
	"for":      FOR,
	"break":    BREAK,
	"continue": CONTINUE,
	"defer":    DEFER,
	"goto":     GOTO,
}

func LoopUp(ident string) TokenType {
	if tok, isKeyword := keywords[ident]; isKeyword {
		return tok
	}
	return IDENT
}
