# tiny-go

A small compiler for a Go-like language, implemented in Go.

`tiny-go` demonstrates a complete compiler pipeline: lexical analysis, parsing, AST construction, LLVM IR generation, and native executable generation through Clang. It also includes experimental WebAssembly build support through LLVM wasm tools.

## Features

- Hand-written lexer for tokenising `.tgo` source files
- Parser for a small Go-like language subset
- AST representation with text and JSON output
- LLVM IR generation for functions, variables, expressions, and control flow
- Native build and run support through Clang
- Experimental WebAssembly output support
- CLI commands for inspecting each compilation stage

## Language subset

The current implementation focuses on a small but useful subset of Go-style syntax:

- `package` declarations
- `import` declarations
- function declarations
- global and local variables
- assignment with `=` and short declaration with `:=`
- integer, float, and character literals
- arithmetic and logical expressions
- `if / else` statements
- `for` loops
- `break`, `continue`, and `return`
- simple built-in function calls, such as `builtin.println(...)`

## Project structure

```text
.
├── ast/          # AST node definitions and printing utilities
├── build/        # Build context: lex, parse, compile, run, and build orchestration
├── builtin/      # Built-in runtime support and embedded LLVM IR
├── compiler/     # AST-to-LLVM IR compiler
├── lexer/        # Tokeniser for tGo source code
├── parser/       # Parser for files, expressions, functions, and statements
├── token/        # Token and source-position definitions
├── main.go       # CLI entry point
├── hello.tgo     # Example tGo source file
└── run_wasm.js   # Helper script for running wasm output
```

## Compilation pipeline

```text
.tgo source
    ↓
Lexer
    ↓
Tokens
    ↓
Parser
    ↓
AST
    ↓
Compiler
    ↓
LLVM IR
    ↓
Clang / LLVM wasm tools
    ↓
Executable / WebAssembly output
```

## Requirements

- Go 1.19 or later
- Clang, for native executable generation
- Optional: LLVM wasm tools, such as `wasm-llc` and `wasm-ld`, for WebAssembly output

## Quick start

Clone the repository:

```bash
git clone https://github.com/BOOSTERLEL/tiny-go.git
cd tiny-go
```

Download Go dependencies:

```bash
go mod download
```

Check the CLI:

```bash
go run . --help
```

Inspect tokens:

```bash
go run . lex hello.tgo
```

Print the AST:

```bash
go run . ast hello.tgo
```

Print the AST as JSON:

```bash
go run . ast --json hello.tgo
```

Generate LLVM IR:

```bash
go run . asm hello.tgo
```

Compile and run a tGo program:

```bash
go run . run hello.tgo
```

Build an executable:

```bash
go run . build hello.tgo
```

The current build command writes the output executable as `a.out.exe`.

## Example tGo program

```go
package main

import "builtin"

func main() {
    builtin.println(42)
}
```

Save the file as `hello.tgo`, then run:

```bash
go run . run hello.tgo
```

Expected output:

```text
42
```

## CLI commands

```bash
tgo run <file>          # Compile and run a tGo program
tgo build <file>        # Compile a tGo source file
tgo lex <file>          # Print the token list and comments
tgo ast <file>          # Parse source code and print the AST
tgo ast --json <file>   # Print the AST in JSON format
tgo asm <file>          # Generate and print LLVM IR
```

Global options:

```bash
--goos       Target operating system
--goarch     Target architecture
--clang      Path to clang
--wasm-llc   Path to wasm llc tool
--wasm-ld    Path to wasm linker
--debug, -d  Keep intermediate build files
```

Example with a custom Clang path:

```bash
go run . --clang /usr/bin/clang run hello.tgo
```

## Implementation notes

The compiler is organised around a clear staged design:

1. `lexer` scans source code into tokens.
2. `parser` converts tokens into AST nodes.
3. `ast` defines the intermediate tree representation.
4. `compiler` lowers the AST into LLVM IR.
5. `build` writes intermediate LLVM files and invokes Clang or wasm tools.
6. `builtin` provides the small runtime layer used by generated programs.

This makes the project useful for learning how a compiler frontend and a simple LLVM-based backend can be connected in Go.

## Development

Run tests:

```bash
go test ./...
```

Format code:

```bash
go fmt ./...
```
