package compiler

import (
	"fmt"
	"tiny-go/ast"
	"tiny-go/token"
)

func (p *Compiler) posLine(pos token.Pos) int {
	if p.file != nil && p.file.Source != "" {
		line := pos.Position(p.file.FileName, p.file.Source).Line
		return line
	}
	return 0
}

// getType 用于获取表达式类型
func getType(expr ast.Expr) (string, bool) {
	switch expr.(type) {
	case *ast.BinaryExpr:
		return getType(expr.(*ast.BinaryExpr).X)
	case *ast.UnaryExpr:
		return getType(expr.(*ast.UnaryExpr).X)
	case *ast.CallExpr:
		return "i32", true
	case *ast.Int:
		return "i32", true
	case *ast.Char:
		return "i8", true
	case *ast.Float:
		return "float", true
	case *ast.Ident:
		return expr.(*ast.Ident).Name, false
	default:
		panic(fmt.Sprintf("unknown: %[1]T, %[1]v", expr))
	}
}

// opType 用于获取表达式操作指令
func opType(op token.TokenType, typ string) string {
	switch op {
	case token.ADD:
		switch typ {
		case "float":
			return "fadd"
		default:
			return "add"
		}
	case token.SUB:
		switch typ {
		case "float":
			return "fsub"
		default:
			return "sub"
		}
	case token.MUL:
		switch typ {
		case "float":
			return "fmul"
		default:
			return "mul"
		}
	case token.DIV:
		switch typ {
		case "float":
			return "fdiv"
		default:
			return "div"
		}
	case token.MOD:
		return "srem"
	case token.EQL:
		switch typ {
		case "float":
			return "fcmp oeq"
		default:
			return "icmp eq"
		}
	case token.NEQ:
		switch typ {
		case "float":
			return "fcmp une"
		default:
			return "icmp ne"
		}
	case token.GTR:
		switch typ {
		case "float":
			return "fcmp ogt"
		default:
			return "icmp sgt"
		}
	case token.GEQ:
		switch typ {
		case "float":
			return "fcmp oge"
		default:
			return "icmp sge"
		}
	case token.LSS:
		switch typ {
		case "float":
			return "fcmp olt"
		default:
			return "icmp slt"
		}
	case token.LEQ:
		switch typ {
		case "float":
			return "fcmp ole"
		default:
			return "icmp sle"
		}
	}
	return ""
}
