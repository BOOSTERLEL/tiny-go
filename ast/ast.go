package ast

import (
	"tiny-go/token"
)

// Node 表示AST中全部结点
type Node interface {
	Pos() token.Pos
	End() token.Pos
	nodeType()
}

// File 表示 µGo 文件对应的语法树.
type File struct {
	FileName string // 文件名
	Source   string // 源代码

	Pkg     *PackageSpec  // 包信息
	Imports []*ImportSpec // 导入包信息
	Globals []*VarSpec    // 全局变量
	Funcs   []*FuncDecl   // 函数列表
}

type PackageSpec struct {
	PkgPos  token.Pos
	NamePos token.Pos
	Name    string
}

// ImportSpec 表示一个导入包
type ImportSpec struct {
	ImportPos token.Pos
	Name      *Ident
	Path      string
}

// VarSpec 变量信息
type VarSpec struct {
	VarPos token.Pos // var 关键字位置
	Name   *Ident    // 变量名字
	Type   *Ident    // 变量类型
	Value  Expr      // 变量表达式
}

// FuncDecl 函数信息
type FuncDecl struct {
	FuncPos token.Pos
	NamePos token.Pos
	Name    string
	Type    *FuncType
	Body    *BlockStmt
}

// FuncLit 闭包函数
type FuncLit struct {
	Type *FuncType
	Body *BlockStmt
}

// FuncType 函数类型
type FuncType struct {
	Func   token.Pos
	Params *FieldList
	Result *Ident
}

// FieldList 参数/属性 列表
type FieldList struct {
	Opening token.Pos
	List    []*Field
	Closing token.Pos
}

// Field 参数/属性
type Field struct {
	Name *Ident
	Type *Ident
}

// DeferStmt defer 语句
type DeferStmt struct {
	DeferPos token.Pos
	Call     *CallExpr
}

// ReturnStmt return 语句
type ReturnStmt struct {
	Return token.Pos
	Result Expr
}

// BranchStmt 分支语句
type BranchStmt struct {
	TokPos  token.Pos
	TokType token.TokenType // BREAK, CONTINUE, GOTO
	Label   *Ident
}

// LabeledStmt 标号语句
type LabeledStmt struct {
	Label *Ident
	Colon token.Pos // 冒号 ":" 位置
	Stmt  Stmt
}

// BlockStmt 块语句
type BlockStmt struct {
	Lbrace token.Pos // '{'
	List   []Stmt
	Rbrace token.Pos // '}'
}

type Stmt interface {
	Pos() token.Pos
	End() token.Pos
	stmtType()
}

type ExprStmt struct {
	X Expr
}

// AssignStmt 表示一个赋值语句节点.
type AssignStmt struct {
	Target []*Ident        // 要赋值的目标
	OpPos  token.Pos       // Op 的位置
	Op     token.TokenType // '=' or ':='
	Value  []Expr          // 值
}

// IfStmt 表示一个 if 语句节点.
type IfStmt struct {
	If   token.Pos  // if 关键字的位置
	Init Stmt       // 初始化语句
	Cond Expr       // if 条件, *BinaryExpr
	Body *BlockStmt // if 为真时对应的语句列表
	Else Stmt       // else 对应的语句
}

// ForStmt 表示一个 for 语句节点.
type ForStmt struct {
	For  token.Pos  // for 关键字的位置
	Init Stmt       // 初始化语句
	Cond Expr       // 条件表达式
	Post Stmt       // 迭代语句
	Body *BlockStmt // 循环对应的语句列表
}

type Expr interface {
	Pos() token.Pos
	End() token.Pos
	exprType()
}

// Ident 标识符
type Ident struct {
	NamePos token.Pos
	Name    string
	Type    string
}

// Int 整型
type Int struct {
	ValuePos token.Pos
	ValueEnd token.Pos
	Value    int
}

// Float 浮点数
type Float struct {
	ValuePos token.Pos
	ValueEnd token.Pos
	Value    float64
}

// Char 字符
type Char struct {
	ValuePos token.Pos
	ValueEnd token.Pos
	Value    int
}

// BinaryExpr 二元表达式
type BinaryExpr struct {
	OpPos token.Pos       // 运算符位置
	Op    token.TokenType // 运算符类型
	X     Expr            // 左边的运算对象
	Y     Expr            // 右边的运算对象
}

// UnaryExpr 一元表达式
type UnaryExpr struct {
	OpPos token.Pos       // 运算符位置
	Op    token.TokenType // 运算符类型
	X     Expr            // 运算对象
}

// ParenExpr 表示一个圆括弧表达式.
type ParenExpr struct {
	Lparen token.Pos // "(" 的位置
	X      Expr      // 圆括弧内的表达式对象
	Rparen token.Pos // ")" 的位置
}

// CallExpr 表示一个函数调用
type CallExpr struct {
	Pkg      *Ident    // 对应的包
	FuncName *Ident    // 函数名字
	Lparen   token.Pos // '(' 位置
	Args     []Expr    // 调用参数列表
	Rparen   token.Pos // ')' 位置
}

// SelectorExpr 表示 x.Name 属性选择表达式
type SelectorExpr struct {
	X   Expr
	Sel *Ident
}
