package compiler

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strings"
	"tiny-go/ast"
	"tiny-go/builtin"
	"tiny-go/token"
)

type Compiler struct {
	file   *ast.File
	scope  *Scope
	nextId int
}

func NewCompiler() *Compiler {
	return &Compiler{
		scope: NewScope(Universe),
	}
}

func (p *Compiler) enterScope() {
	p.scope = NewScope(p.scope)
}

func (p *Compiler) leaveScope() {
	p.scope = p.scope.Outer
}

func (p *Compiler) restoreScope(scope *Scope) {
	p.scope = scope
}

func (p *Compiler) genHeader(w io.Writer, file *ast.File) {
	_, _ = fmt.Fprintf(w, ";package %s\n", file.Pkg.Name)
	_, _ = fmt.Fprintf(w, builtin.Header)
}

func (p *Compiler) genMain(w io.Writer, file *ast.File) {
	if file.Pkg.Name != "main" {
		return
	}
	for _, fn := range file.Funcs {
		if fn.Name == "main" {
			_, _ = fmt.Fprintf(w, builtin.MainMain)
			return
		}
	}
}

func (p *Compiler) genInit(w io.Writer, file *ast.File) {
	_, _ = fmt.Fprintf(w, "define i32 @tiny_go_%s_init() {\n", file.Pkg.Name)

	for _, g := range file.Globals {
		var localName = "0"
		if g.Value != nil {
			localName = p.compileExpr(w, g.Value)
		}

		var varName string
		if _, obj := p.scope.Lookup(g.Name.Name); obj != nil {
			varName = obj.MangledName
		} else {
			panic(fmt.Sprintf("var %v undefined", g))
		}

		_, _ = fmt.Fprintf(w, "\tstore i32 %s, i32* %s\n", localName, varName)
	}
	_, _ = fmt.Fprintln(w, "\tret i32 0")
	_, _ = fmt.Fprintln(w, "}")
}

func (p *Compiler) genLabelId(name string) string {
	id := fmt.Sprintf("%s.%d", name, p.nextId)
	p.nextId++
	return id
}

func (p *Compiler) compileFile(w io.Writer, file *ast.File) {
	defer p.restoreScope(p.scope)
	p.enterScope()

	// import
	for _, x := range file.Imports {
		var mangledName = fmt.Sprintf("@tiny_go_%s", x.Path)
		if x.Name != nil {
			p.scope.Insert(&Object{
				Name:        x.Name.Name,
				MangledName: mangledName,
				Node:        x,
			})
		} else {
			p.scope.Insert(&Object{
				Name:        x.Path,
				MangledName: mangledName,
				Node:        x,
			})
		}
	}

	// global vars
	for _, g := range file.Globals {
		var mangledName = fmt.Sprintf("@tiny_go_%s_%s", file.Pkg.Name, g.Name.Name)
		p.scope.Insert(&Object{
			Name:        g.Name.Name,
			MangledName: mangledName,
			Type:        g.Type.Type,
			Node:        g,
		})
		_, _ = fmt.Fprintf(w, "%s = global i32 0\n", mangledName)
	}

	// global funcs
	for _, fn := range file.Funcs {
		var mangledName = fmt.Sprintf("@tiny_go_%s_%s", file.Pkg.Name, fn.Name)

		// func type
		var typ string
		if fn.Type.Result == nil {
			typ = "i32"
		} else {
			typ = fn.Type.Result.Type
		}
		p.scope.Insert(&Object{
			Name:        fn.Name,
			MangledName: mangledName,
			Type:        typ,
			Node:        fn,
		})
	}
	p.genInit(w, file)
	for _, fn := range file.Funcs {
		p.compileFunc(w, file, fn)
	}
}

func (p *Compiler) compileFunc(w io.Writer, file *ast.File, fn *ast.FuncDecl) {
	defer p.restoreScope(p.scope)
	p.enterScope()

	// args
	var argNameList []string
	var argTypeList []string
	for _, arg := range fn.Type.Params.List {
		var mangledName = fmt.Sprintf("%%local_%s.pos.%d", arg.Name.Name, arg.Name.NamePos)
		var argType = arg.Type.Type
		argNameList = append(argNameList, mangledName)
		argTypeList = append(argTypeList, argType)
	}

	// result type
	var typ string
	if fn.Type.Result == nil {
		typ = "i32"
	} else {
		typ = fn.Type.Result.Type
	}

	if fn.Body == nil {
		_, _ = fmt.Fprintf(w, "declare %s @tiny_go_%s_%s()\n", typ, file.Pkg.Name, fn.Name)
		return
	}
	_, _ = fmt.Fprintf(w, "define %s @tiny_go_%s_%s(", typ, file.Pkg.Name, fn.Name)
	var first = true
	for i, argRegName := range argNameList {
		if first {
			first = false
			_, _ = fmt.Fprintf(w, "%s noundef %s.arg%d", argTypeList[i], argRegName, i)
			continue
		}
		_, _ = fmt.Fprintf(w, ", %s noundef %s.arg%d", argTypeList[i], argRegName, i)
	}
	_, _ = fmt.Fprintf(w, ") {\n")

	// fn body
	func() {
		// args+body scope
		defer p.restoreScope(p.scope)
		p.enterScope()

		// args
		for i, arg := range fn.Type.Params.List {
			var argRegName = fmt.Sprintf("%s.arg%d", argNameList[i], i)
			var mangledName = argNameList[i]
			p.scope.Insert(&Object{
				Name:        arg.Name.Name,
				MangledName: mangledName,
				Type:        arg.Type.Type,
				Node:        fn,
			})

			_, _ = fmt.Fprintf(w, "\t%s = alloca %s, align 4\n", mangledName, arg.Type.Type)
			_, _ = fmt.Fprintf(w, "\tstore %s %s, %s* %s\n", arg.Type.Type, argRegName, arg.Type.Type, mangledName)
		}

		// body
		for _, x := range fn.Body.List {
			p.compileStmt(w, x)
		}
	}()

	if fn.Type.Result == nil {
		_, _ = fmt.Fprintf(w, "\tret %s 0\n", typ)
	}
	_, _ = fmt.Fprintln(w, "}")
}

func (p *Compiler) compileStmt(w io.Writer, stmt ast.Stmt) {
	switch stmt := stmt.(type) {
	case *ast.VarSpec:
		var localName = "0"
		if stmt.Value != nil {
			localName = p.compileExpr(w, stmt.Value)
		}

		var mangledName = fmt.Sprintf("%%local_%s.pos.%d", stmt.Name.Name, stmt.VarPos)
		var width int
		p.scope.Insert(&Object{
			Name:        stmt.Name.Name,
			MangledName: mangledName,
			Type:        stmt.Type.Type,
			Node:        stmt,
		})
		switch stmt.Type.Type {
		case "i8":
			width = 1
		default:
			width = 4
		}

		_, _ = fmt.Fprintf(w, "\t%s = alloca %s, align %d\n", mangledName, stmt.Type.Type, width)
		_, _ = fmt.Fprintf(w, "\tstore %s %s, %s* %s\n", stmt.Type.Type, localName, stmt.Type.Type, mangledName)
	case *ast.AssignStmt:
		p.compileStmtAssign(w, stmt)
	case *ast.ReturnStmt:
		p.compileStmtReturn(w, stmt)
	case *ast.IfStmt:
		p.compileStmtIf(w, stmt)
	case *ast.ForStmt:
		p.compileStmtFor(w, stmt)
	case *ast.BranchStmt:
		p.compileStmtBranch(w, stmt)
	case *ast.BlockStmt:
		defer p.restoreScope(p.scope)
		p.enterScope()
		for _, x := range stmt.List {
			p.compileStmt(w, x)
		}
	case *ast.ExprStmt:
		p.compileExpr(w, stmt.X)
	default:
		panic("unreachable")
	}
}

func (p *Compiler) compileStmtReturn(w io.Writer, stmt *ast.ReturnStmt) {
	if stmt.Result != nil {
		typ, flag := getType(stmt.Result)
		if !flag {
			if _, obj := p.scope.Lookup(typ); obj != nil {
				typ = obj.Type
			} else {
				panic(fmt.Sprintf("var %s undefined", typ))
			}
		}
		_, _ = fmt.Fprintf(w, "\tret %s %v\n", typ, p.compileExpr(w, stmt.Result))
	} else {
		_, _ = fmt.Fprintf(w, "\tret i32 0\n")
	}
}

func (p *Compiler) compileStmtBranch(w io.Writer, stmt *ast.BranchStmt) {
	var branch []string
	if _, obj := p.scope.Lookup("for"); obj != nil {
		branch = strings.Fields(obj.MangledName)
	} else {
		panic(fmt.Sprintf("for undefined"))
	}
	var i int
	switch stmt.TokType {
	case token.BREAK:
		i = 1
	case token.CONTINUE:
		i = 0
	default:
		panic("unreachable")
	}
	_, _ = fmt.Fprintf(w, "\tbr label %%%s\n", branch[i])
}

func (p *Compiler) compileStmtAssign(w io.Writer, stmt *ast.AssignStmt) {
	var varNameList = make([]string, len(stmt.Value))
	for i := range stmt.Target {
		varNameList[i] = p.compileExpr(w, stmt.Value[i])
	}

	if stmt.Op == token.DEFINE {
		for i, target := range stmt.Target {
			if _, obj := p.scope.Lookup(target.Name); obj == nil {
				var mangledName = fmt.Sprintf("%%local_%s.pos.%d", target.Name, target.NamePos)
				var width int
				typ, flag := getType(stmt.Value[i])
				if !flag {
					if _, obj := p.scope.Lookup(typ); obj != nil {
						typ = obj.Type
					} else {
						panic(fmt.Sprintf("var %s undefined", typ))
					}
				}
				p.scope.Insert(&Object{
					Name:        target.Name,
					MangledName: mangledName,
					Type:        typ,
					Node:        target,
				})
				switch typ {
				case "i32":
					width = 4
				case "i8":
					width = 1
				default:
					width = 4
				}

				_, _ = fmt.Fprintf(w, "\t%s = alloca %s, align %d\n", mangledName, typ, width)
			}
		}
	}

	for i := range stmt.Target {
		var varName string
		var typ string
		if _, obj := p.scope.Lookup(stmt.Target[i].Name); obj != nil {
			varName = obj.MangledName
			typ = obj.Type
		} else {
			panic(fmt.Sprintf("var %s undefined", stmt.Target[0].Name))
		}
		_, _ = fmt.Fprintf(w, "\tstore %s %s, %s* %s\n", typ, varNameList[i], typ, varName)
	}
}

func (p *Compiler) compileStmtCond(stmt ast.Expr) []string {
	var condList []string
	switch stmt.(type) {
	case *ast.BinaryExpr:
		switch stmt.(*ast.BinaryExpr).Op {
		case token.AND, token.OR:
			condList = append(condList, p.compileStmtCond(stmt.(*ast.BinaryExpr).X)...)
			condList = append(condList, p.compileStmtCond(stmt.(*ast.BinaryExpr).Y)...)
		default:
			condList = append(condList, "if.cond.line"+string(rune(stmt.(*ast.BinaryExpr).OpPos)))
		}
	case *ast.UnaryExpr:
		p.convertNot(stmt.(*ast.UnaryExpr).X)
		condList = append(condList, p.compileStmtCond(stmt.(*ast.UnaryExpr).X)...)
	}
	return condList
}

func (p *Compiler) convertNot(stmt ast.Expr) {
	switch stmt.(type) {
	case *ast.UnaryExpr:
		p.convertNot(stmt.(*ast.UnaryExpr).X)
	case *ast.BinaryExpr:
		switch stmt.(*ast.BinaryExpr).Op {
		case token.AND:
			stmt.(*ast.BinaryExpr).Op = token.OR
			p.convertNot(stmt.(*ast.BinaryExpr).X)
			p.convertNot(stmt.(*ast.BinaryExpr).Y)
		case token.OR:
			stmt.(*ast.BinaryExpr).Op = token.AND
			p.convertNot(stmt.(*ast.BinaryExpr).X)
			p.convertNot(stmt.(*ast.BinaryExpr).Y)
		case token.EQL:
			stmt.(*ast.BinaryExpr).Op = token.NEQ
		case token.NEQ:
			stmt.(*ast.BinaryExpr).Op = token.EQL
		case token.GTR:
			stmt.(*ast.BinaryExpr).Op = token.LEQ
		case token.GEQ:
			stmt.(*ast.BinaryExpr).Op = token.LSS
		case token.LSS:
			stmt.(*ast.BinaryExpr).Op = token.GEQ
		case token.LEQ:
			stmt.(*ast.BinaryExpr).Op = token.GTR
		}
	}
}

func (p *Compiler) compileStmtCondTree(stmt ast.Expr) []ast.Expr {
	var treeList []ast.Expr
	switch stmt.(type) {
	case *ast.BinaryExpr:
		switch stmt.(*ast.BinaryExpr).Op {
		case token.AND, token.OR:
			treeList = append(treeList, p.compileStmtCondTree(stmt.(*ast.BinaryExpr).X)...)
			treeList = append(treeList, p.compileStmtCondTree(stmt.(*ast.BinaryExpr).Y)...)
		default:
			treeList = append(treeList, stmt)
		}
	case *ast.UnaryExpr:
		treeList = append(treeList, p.compileStmtCondTree(stmt.(*ast.UnaryExpr).X)...)
	}
	return treeList
}

func (p *Compiler) findFLag(stmt ast.Expr, fOp token.TokenType) []bool {
	var flags []bool
	switch stmt.(type) {
	case *ast.BinaryExpr:
		switch stmt.(*ast.BinaryExpr).Op {
		case token.AND, token.OR:
			flags = append(flags, p.findFLag(stmt.(*ast.BinaryExpr).X, stmt.(*ast.BinaryExpr).Op)...)
			flags = append(flags, p.findFLag(stmt.(*ast.BinaryExpr).Y, token.IDENT)...)
		default:
			switch fOp {
			case token.AND:
				flags = append(flags, true)
			case token.OR:
				flags = append(flags, false)
			case token.IDENT:
				flags = append(flags, false)
			}
		}
	case *ast.UnaryExpr:
		flags = append(flags, p.findFLag(stmt.(*ast.UnaryExpr).X, stmt.(*ast.UnaryExpr).Op)...)
	}
	return flags
}

func (p *Compiler) compileStmtIf(w io.Writer, stmt *ast.IfStmt) {
	defer p.restoreScope(p.scope)
	p.enterScope()

	ifPos := fmt.Sprintf("%d", p.posLine(stmt.If))
	ifInit := p.genLabelId("if.init.line" + ifPos)
	ifCond := p.genLabelId("if.cond.line" + ifPos)
	ifBody := p.genLabelId("if.body.line" + ifPos)
	ifElse := p.genLabelId("if.else.line" + ifPos)
	ifEnd := p.genLabelId("if.end.line" + ifPos)

	// br if.init
	_, _ = fmt.Fprintf(w, "\tbr label %%%s\n", ifInit)

	// if.init
	_, _ = fmt.Fprintf(w, "\n%s:\n", ifInit)
	func() {
		defer p.restoreScope(p.scope)
		p.enterScope()

		if stmt.Init != nil {
			p.compileStmt(w, stmt.Init)
		}
		_, _ = fmt.Fprintf(w, "\tbr label %%%s\n", ifCond)

		//if.cond
		{
			_, _ = fmt.Fprintf(w, "\n%s:\n", ifCond)
			condList := p.compileStmtCond(stmt.Cond)
			valueList := p.compileStmtCondTree(stmt.Cond)
			flags := p.findFLag(stmt.Cond, token.IDENT)
			var ifFalse string
			var ifTrue string
			var jumpTo string
			for i := 0; i < len(condList)-1; i++ {
				condValue := p.compileExpr(w, valueList[i])
				if stmt.Else != nil {
					jumpTo = ifElse
				} else {
					jumpTo = ifEnd
				}
				if flags[i] {
					ifTrue = condList[i+1]
					ifFalse = jumpTo
				} else {
					ifTrue = ifBody
					ifFalse = condList[i+1]
				}
				_, _ = fmt.Fprintf(w, "\tbr i1 %s, label %%%s, label %%%s\n", condValue, ifTrue, ifFalse)
				_, _ = fmt.Fprintf(w, "\n%s:\n", condList[i+1])
			}
			condValue := p.compileExpr(w, valueList[len(condList)-1])
			if stmt.Else != nil {
				_, _ = fmt.Fprintf(w, "\tbr i1 %s, label %%%s, label %%%s\n", condValue, ifBody, ifElse)
			} else {
				_, _ = fmt.Fprintf(w, "\tbr i1 %s, label %%%s, label %%%s\n", condValue, ifBody, ifEnd)
			}
			//condValue := p.compileExpr(w, stmt.Cond)
			//if stmt.Else != nil {
			//	_, _ = fmt.Fprintf(w, "\tbr i1 %s, label %%%s, label %%%s\n", condValue, ifBody, ifElse)
			//} else {
			//	_, _ = fmt.Fprintf(w, "\tbr i1 %s, label %%%s, label %%%s\n", condValue, ifBody, ifEnd)
			//}
		}

		// if.body
		func() {
			defer p.restoreScope(p.scope)
			p.enterScope()

			_, _ = fmt.Fprintf(w, "\n%s:\n", ifBody)
			p.compileStmt(w, stmt.Body)
			if stmt.Else != nil {
				_, _ = fmt.Fprintf(w, "\tbr label %%%s\n", ifElse)
			} else {
				_, _ = fmt.Fprintf(w, "\tbr label %%%s\n", ifEnd)
			}
		}()

		// if.else
		func() {
			defer p.restoreScope(p.scope)
			p.enterScope()

			_, _ = fmt.Fprintf(w, "\n%s:\n", ifElse)
			if stmt.Else != nil {
				p.compileStmt(w, stmt.Else)
			}
			_, _ = fmt.Fprintf(w, "\tbr label %%%s\n", ifEnd)
		}()
	}()

	// end
	_, _ = fmt.Fprintf(w, "\n%s:\n", ifEnd)
}

func (p *Compiler) compileStmtFor(w io.Writer, stmt *ast.ForStmt) {
	defer p.restoreScope(p.scope)
	p.enterScope()

	forPos := fmt.Sprintf("%d", p.posLine(stmt.For))
	forInit := p.genLabelId("for.init.line" + forPos)
	forCond := p.genLabelId("for.cond.line" + forPos)
	forPost := p.genLabelId("for.post.line" + forPos)
	forBody := p.genLabelId("for.body.line" + forPos)
	forEnd := p.genLabelId("for.end.line" + forPos)

	p.scope.Insert(&Object{
		Name:        "for",
		MangledName: forPost + " " + forEnd,
		Type:        "for",
	})

	// br for.init
	_, _ = fmt.Fprintf(w, "\tbr label %%%s\n", forInit)

	// for.init
	func() {
		defer p.restoreScope(p.scope)
		p.enterScope()

		_, _ = fmt.Fprintf(w, "\n%s:\n", forInit)
		if stmt.Init != nil {
			p.compileStmt(w, stmt.Init)
		}
		_, _ = fmt.Fprintf(w, "\tbr label %%%s\n", forCond)

		// for.cond
		_, _ = fmt.Fprintf(w, "\n%s:\n", forCond)
		if stmt.Cond != nil {
			condValue := p.compileExpr(w, stmt.Cond)
			_, _ = fmt.Fprintf(w, "\tbr i1 %s, label %%%s, label %%%s\n", condValue, forBody, forEnd)
		} else {
			_, _ = fmt.Fprintf(w, "\tbr label %%%s\n", forBody)
		}

		// for.body
		func() {
			defer p.restoreScope(p.scope)
			p.enterScope()

			_, _ = fmt.Fprintf(w, "\n%s:\n", forBody)
			p.compileStmt(w, stmt.Body)
			_, _ = fmt.Fprintf(w, "\tbr label %%%s\n", forPost)
		}()

		// for.post
		{
			_, _ = fmt.Fprintf(w, "\n%s:\n", forPost)
			if stmt.Post != nil {
				p.compileStmt(w, stmt.Post)
			}
			_, _ = fmt.Fprintf(w, "\tbr label %%%s\n", forCond)
		}
	}()

	//end
	_, _ = fmt.Fprintf(w, "\n%s:\n", forEnd)
}

func (p *Compiler) compileExpr(w io.Writer, expr ast.Expr) (localName string) {
	switch expr := expr.(type) {
	case *ast.Ident:
		var varName string
		var typ string
		var width int
		if _, obj := p.scope.Lookup(expr.Name); obj != nil {
			varName = obj.MangledName
			typ = obj.Type
			switch typ {
			case "i32":
				width = 4
			case "i8":
				width = 1
			default:
				width = 4
			}
		} else {
			panic(fmt.Sprintf("var %s undefined", expr.Name))
		}

		localName = p.genId()
		_, _ = fmt.Fprintf(w, "\t%s = load %s, %s* %s, align %d\n", localName, typ, typ, varName, width)
		return localName

	case *ast.Int:
		localName = p.genId()
		_, _ = fmt.Fprintf(w, "\t%s = %s i32 %v, %v\n", localName,
			"add", `0`, expr.Value)
		return localName

	case *ast.Float:
		localName = p.genId()
		_, _ = fmt.Fprintf(w, "\t%s = %s float %v, 0x%v\n", localName,
			"fadd", `0.000000e+00`, math.Float64bits(expr.Value))
		return localName

	case *ast.Char:
		localName = p.genId()
		_, _ = fmt.Fprintf(w, "\t%s = %s i8 %v, %v\n", localName, "add", `0`, expr.Value)
		return localName
	case *ast.BinaryExpr:
		typ, flag := getType(expr)
		if !flag {
			if _, obj := p.scope.Lookup(typ); obj != nil {
				typ = obj.Type
			} else {
				panic(fmt.Sprintf("var %s undefined", typ))
			}
		}
		if typ == "i8" {
			typ = "i32"
		}
		localName = p.genId()
		op := opType(expr.Op, typ)
		x := p.compileExpr(w, expr.X)
		y := p.compileExpr(w, expr.Y)
		typX := ""
		typY := ""
		switch expr.X.(type) {
		case *ast.Ident:
			if _, obj := p.scope.Lookup(expr.X.(*ast.Ident).Name); obj != nil {
				typX = obj.Type
				x = p.convert(w, x, typX, "i32")
			} else {
				panic(fmt.Sprintf("var %s undefined", typ))
			}
		case *ast.Char:
			x = p.convert(w, x, "i8", "i32")
		}
		switch expr.Y.(type) {
		case *ast.Ident:
			if _, obj := p.scope.Lookup(expr.Y.(*ast.Ident).Name); obj != nil {
				typY = obj.Type
				y = p.convert(w, y, typY, "i32")
			} else {
				panic(fmt.Sprintf("var %s undefined", typ))
			}
		case *ast.Char:
			y = p.convert(w, y, "i8", "i32")
		}
		_, _ = fmt.Fprintf(w, "\t%s = %s %s %v, %v\n", localName, op, typ, x, y)
		if typX != "" {
			localName = p.convert(w, localName, "i32", typX)
		}

		return localName

	case *ast.UnaryExpr:
		typ, _ := getType(expr)
		op := opType(expr.Op, typ)
		var zero string
		switch typ {
		case "i32":
			zero = `0`
		case "float":
			zero = `0.000000e+00`
		}
		if expr.Op == token.SUB {
			localName = p.genId()
			_, _ = fmt.Fprintf(w, "\t%s = %s %s %v, %v\n",
				localName, op, typ, zero, p.compileExpr(w, expr.X))
			return localName
		}
		return p.compileExpr(w, expr.X)

	case *ast.ParenExpr:
		return p.compileExpr(w, expr.X)

	case *ast.CallExpr:
		var fnName string
		var fnType string
		var paramsType []string
		if expr.Pkg != nil {
			if _, obj := p.scope.Lookup(expr.Pkg.Name); obj != nil {
				fnName = obj.MangledName + "_" + expr.FuncName.Name
				fnType = obj.Type
				if obj.Node != nil {
					//node := obj.Node.(*ast.FuncDecl).Type.Params.List
					//for _, typ := range node {
					//	paramsType = append(paramsType, typ.Type.Type)
					//}
					paramsType = append(paramsType, "i32")
				}
			} else {
				panic(fmt.Sprintf("func %s.%s undefined", expr.Pkg.Name, expr.FuncName.Name))
			}
		} else if _, obj := p.scope.Lookup(expr.FuncName.Name); obj != nil {
			fnName = obj.MangledName
			fnType = obj.Type
			if obj.Node != nil {
				node := obj.Node.(*ast.FuncDecl).Type.Params.List
				for _, typ := range node {
					paramsType = append(paramsType, typ.Type.Type)
				}
			}
		} else {
			panic(fmt.Sprintf("func %s undefined", expr.FuncName.Name))
		}

		// test func
		if fnType == "" {
			fnType = "i32"
		}
		if len(paramsType) == 0 {
			paramsType = append(paramsType, "i32")
		}

		localName = p.genId()
		var localNames []string
		for _, arg := range expr.Args {
			localNames = append(localNames, p.compileExpr(w, arg))
		}
		_, _ = fmt.Fprintf(w, "\t%s = call %s(", localName, fnType)
		first := true
		for _, paramsType := range paramsType {
			if first {
				_, _ = fmt.Fprintf(w, "%s", paramsType)
				first = false
				continue
			}
			_, _ = fmt.Fprintf(w, ", %s", paramsType)
		}
		_, _ = fmt.Fprintf(w, ") %s(", fnName)
		first = true
		for i, paramType := range paramsType {
			if first {
				_, _ = fmt.Fprintf(w, "%s noundef %s", paramType, localNames[0])
				first = false
				continue
			}
			_, _ = fmt.Fprintf(w, ",%s noundef %s", paramType, localNames[i])
		}
		_, _ = fmt.Fprintf(w, ")\n")
		//_, _ = fmt.Fprintf(w, "\t%s = call %s(i32) %s(i32 %v)\n",
		//	localName, fnType, fnName, p.compileExpr(w, expr.Args[0]))
		return localName

	default:
		panic(fmt.Sprintf("unknown: %[1]T, %[1]v", expr))
	}
}

func (p *Compiler) genId() string {
	id := fmt.Sprintf("%%t%d", p.nextId)
	p.nextId++
	return id
}

func (p *Compiler) convert(w io.Writer, localName, typ, NewTyp string) string {
	if typ == NewTyp {
		return ""
	}
	emitName := p.genId()
	var op string
	switch NewTyp {
	case "i32":
		op = "sext"
	case "i8":
		op = "trunc"
	}
	_, _ = fmt.Fprintf(w, "\t%s = %s %s %s to %s\n", emitName, op, typ, localName, NewTyp)
	return emitName
}

func (p *Compiler) Compile(f *ast.File) string {
	var buf bytes.Buffer

	p.genHeader(&buf, f)
	p.compileFile(&buf, f)
	p.genMain(&buf, f)

	return buf.String()
}
