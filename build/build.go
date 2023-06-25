package build

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"tiny-go/ast"
	"tiny-go/builtin"
	"tiny-go/compiler"
	"tiny-go/lexer"
	"tiny-go/parser"
	"tiny-go/token"
)

type Option struct {
	Debug   bool
	GOOS    string
	GOARCH  string
	Clang   string
	WasmLLC string
	WasmLD  string
}

type Context struct {
	opt  Option
	path string
	src  string
}

func NewContext(opt *Option) *Context {
	p := &Context{}
	if opt != nil {
		p.opt = *opt
	}
	if p.opt.Clang == "" {
		if runtime.GOOS == "windows" {
			p.opt.Clang, _ = exec.LookPath("clang.exe")
		} else {
			p.opt.Clang, _ = exec.LookPath("clang")
		}
		if p.opt.Clang == "" {
			p.opt.Clang = "clang"
		}
	}
	if p.opt.GOOS == "" {
		p.opt.GOOS = runtime.GOOS
	}
	if p.opt.GOARCH == "" {
		p.opt.GOARCH = runtime.GOARCH
	}
	return p
}

func (p *Context) Lex(fileName string, src interface{}) (tokens, comments []token.Token, err error) {
	code, err := p.readSource(fileName, src)
	if err != nil {
		return nil, nil, err
	}
	l := lexer.NewLexer(fileName, code)
	tokens = l.Tokens()
	comments = l.Comments()
	return
}

func (p *Context) AST(fileName string, src interface{}) (f *ast.File, err error) {
	code, err := p.readSource(fileName, src)
	if err != nil {
		return nil, err
	}
	f, err = parser.ParseFile(fileName, code)
	if err != nil {
		return nil, err
	}
	return
}

func (p *Context) ASM(fileName string, src interface{}) (ll string, err error) {
	code, err := p.readSource(fileName, src)
	if err != nil {
		return "", err
	}
	f, err := parser.ParseFile(fileName, code)
	if err != nil {
		return "", err
	}
	ll = new(compiler.Compiler).Compile(f)
	return
}

func (p *Context) Build(fileName string, src interface{}, outFIle string) (output []byte, err error) {
	return p.build(fileName, src, outFIle, p.opt.GOOS, p.opt.GOARCH)
}

func (p *Context) build(fileName string, src interface{}, outFile, goos, goarch string) (output []byte, err error) {
	code, err := p.readSource(fileName, src)
	if err != nil {
		return nil, err
	}
	f, err := parser.ParseFile(fileName, code)
	if err != nil {
		return nil, err
	}

	const (
		_a_out_ll         = ".\\builtin\\_a.out.ll"
		_a_out_ll_o       = ".\\builtin\\_a.out.ll.o"
		_a_out_builtin_ll = ".\\builtin\\_a.out.builtin.ll"
	)
	if !p.opt.Debug {
		defer os.Remove(_a_out_ll)
		defer os.Remove(_a_out_ll_o)
		defer os.Remove(_a_out_builtin_ll)
	}

	llBuiltin := builtin.GetBuiltinLL(p.opt.GOOS, p.opt.GOARCH)
	err = os.WriteFile(_a_out_builtin_ll, []byte(llBuiltin), 0666)
	if err != nil {
		return nil, err
	}

	ll := compiler.NewCompiler().Compile(f)
	err = os.WriteFile(_a_out_ll, []byte(ll), 0666)
	if err != nil {
		return nil, err
	}

	if outFile == "" {
		outFile = "a.out"
	}
	if p.opt.GOOS == "wasm" {
		if !strings.HasSuffix(outFile, ".wasm") {
			outFile += ".wasm"
		}

		cmdLLC := exec.Command(p.opt.WasmLLC, "-march=wasm32", "-filetype=obj", "-o", _a_out_ll_o, _a_out_ll)
		if data, err := cmdLLC.CombinedOutput(); err != nil {
			return data, err
		}

		cmdWasmLD := exec.Command(p.opt.WasmLD, "--entry=main", "--allow-undefined",
			"--export-all", _a_out_ll_o, "-o", outFile)
		data, err := cmdWasmLD.CombinedOutput()
		return data, err
	}
	cmd := exec.Command(
		p.opt.Clang, "-Wno-override-module", "-o", outFile,
		_a_out_ll, _a_out_builtin_ll,
	)

	data, err := cmd.CombinedOutput()
	return data, err
}

func (p *Context) Run(fileName string, src interface{}) ([]byte, error) {
	if p.opt.GOOS == "wasm" {
		return nil, fmt.Errorf("donot support run wasm")
	}

	a_out := "./a.out"
	if runtime.GOOS == "windows" {
		a_out = ".\\a.out.exe"
	}
	if !p.opt.Debug {
		defer os.Remove(a_out)
	}

	output, err := p.build(fileName, src, a_out, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return output, err
	}

	output, err = exec.Command(a_out).CombinedOutput()
	if err != nil {
		return output, err
	}
	return output, nil
}

func (p *Context) readSource(fileName string, src interface{}) (string, error) {
	if src != nil {
		switch s := src.(type) {
		case string:
			return s, nil
		case []byte:
			return string(s), nil
		case *bytes.Buffer:
			if s != nil {
				return s.String(), nil
			}
		case io.Reader:
			d, err := io.ReadAll(s)
			return string(d), err
		}
		return "", errors.New("invalid source")
	}
	d, err := os.ReadFile(fileName)
	return string(d), err
}
