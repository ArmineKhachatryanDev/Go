## 
Compiled Versus Interpreted Languages
Every program is a set of instructions, whether it’s to add two numbers or send a request over the internet. Compilers and interpreters take human-readable code and convert it to computer-readable machine code. In a compiled language, the target machine directly translates the program. In an interpreted language, the source code is not directly translated by the target machine. Instead, a different program, aka the interpreter, reads and executes the code.

## Go
We will take a look at three phases of the compiler:

* The scanner, which converts the source code into a list of tokens, for use by the parser.
* The parser, which converts the tokens into an Abstract Syntax Tree to be used by code generation.
* The code generation, which converts the Abstract Syntax Tree to machine code.

Note: The packages we are going to be using (go/scanner, go/parser, go/token, go/ast, etc.) are not used by the Go compiler, but are mainly provided for use by tools to operate on Go source code. However, the actual Go compiler has very similar semantics. It does not use these packages because the compiler was once written in C and converted to Go code, so the actual Go compiler is still reminiscent of that structure.

The first step of every compiler is to break up the raw source code text into tokens, which is done by the scanner (also known as lexer). Tokens can be keywords, strings, variable names, function names, etc. Every valid program “word” is represented by a token. In concrete terms for Go, this might mean we have a token “package”, “main”, “func” and so forth.

Each token is represented by its position, type, and raw text in Go. Go even allows us to execute the scanner ourselves in a Go program by using the go/scanner and go/token packages. That means we can inspect what our program looks like to the Go compiler after it has been scanned. To do so, we are going to create a simple program that prints all tokens of a Hello World program.

The program will look like this:


![alt text](https://github.com/ArmineKhachatryanDev/Go/blob/master/howGoCompilesDownToMachineCode/scanner.png)

## Code
```go
package main

import (
  "fmt"
  "go/scanner"
  "go/token"
)

func main() {
  src := []byte(`package main

import "fmt"

func main() {
  fmt.Println("Hello, world!")
}
`)

  var s scanner.Scanner
  fset := token.NewFileSet()
  file := fset.AddFile("", fset.Base(), len(src))
  s.Init(file, src, nil, 0)

  for {
     pos, tok, lit := s.Scan()
     fmt.Printf("%-6s%-8s%q\n", fset.Position(pos), tok, lit)

     if tok == token.EOF {
        break
     }
  }
```

  We will create our source code string and initialize the scanner.Scanner struct which will scan our source code. We call Scan() as many times as we can and print the token’s position, type, and literal string until we reach the End of File (EOF) marker.


  When we run the program, it will print the following:

  ![alt text](https://github.com/ArmineKhachatryanDev/Go/blob/master/howGoCompilesDownToMachineCode/scannerResult.png)



Here we can see what the Go parser uses when it compiles a program. What we can also see is that the scanner adds semicolons where those would usually be placed in other programming languages such as C. This explains why Go does not need semicolons: they are placed intelligently by the scanner.

After the source code has been scanned, it will be passed to the parser. The parser is a phase of the compiler that converts the tokens into an Abstract Syntax Tree (AST). The AST is a structured representation of the source code. In the AST we will be able to see the program structure, such as functions and constant declarations.

Go has again provided us with packages to parse the program and view the AST: go/parser and go/ast. We can use them like this to print the full AST:

![alt text](https://github.com/ArmineKhachatryanDev/Go/blob/master/howGoCompilesDownToMachineCode/Parser.png)

##Code
```go

package main

import (
  "go/ast"
  "go/parser"
  "go/token"
  "log"
)

func main() {
  src := []byte(`package main
import "fmt"
func main() {
  fmt.Println("Hello, world!")
}
`)

  fset := token.NewFileSet()

  file, err := parser.ParseFile(fset, "", src, 0)
  if err != nil {
     log.Fatal(err)
  }

  ast.Print(fset, file)
}

Output:
     0  *ast.File {
     1  .  Package: 1:1
     2  .  Name: *ast.Ident {
     3  .  .  NamePos: 1:9
     4  .  .  Name: "main"
     5  .  }
     6  .  Decls: []ast.Decl (len = 2) {
     7  .  .  0: *ast.GenDecl {
     8  .  .  .  TokPos: 3:1
     9  .  .  .  Tok: import
    10  .  .  .  Lparen: -
    11  .  .  .  Specs: []ast.Spec (len = 1) {
    12  .  .  .  .  0: *ast.ImportSpec {
    13  .  .  .  .  .  Path: *ast.BasicLit {
    14  .  .  .  .  .  .  ValuePos: 3:8
    15  .  .  .  .  .  .  Kind: STRING
    16  .  .  .  .  .  .  Value: "\"fmt\""
    17  .  .  .  .  .  }
    18  .  .  .  .  .  EndPos: -
    19  .  .  .  .  }
    20  .  .  .  }
    21  .  .  .  Rparen: -
    22  .  .  }
    23  .  .  1: *ast.FuncDecl {
    24  .  .  .  Name: *ast.Ident {
    25  .  .  .  .  NamePos: 5:6
    26  .  .  .  .  Name: "main"
    27  .  .  .  .  Obj: *ast.Object {
    28  .  .  .  .  .  Kind: func
    29  .  .  .  .  .  Name: "main"
    30  .  .  .  .  .  Decl: *(obj @ 23)
    31  .  .  .  .  }
    32  .  .  .  }
    33  .  .  .  Type: *ast.FuncType {
    34  .  .  .  .  Func: 5:1
    35  .  .  .  .  Params: *ast.FieldList {
    36  .  .  .  .  .  Opening: 5:10
    37  .  .  .  .  .  Closing: 5:11
    38  .  .  .  .  }
    39  .  .  .  }
    40  .  .  .  Body: *ast.BlockStmt {
    41  .  .  .  .  Lbrace: 5:13
    42  .  .  .  .  List: []ast.Stmt (len = 1) {
    43  .  .  .  .  .  0: *ast.ExprStmt {
    44  .  .  .  .  .  .  X: *ast.CallExpr {
    45  .  .  .  .  .  .  .  Fun: *ast.SelectorExpr {
    46  .  .  .  .  .  .  .  .  X: *ast.Ident {
    47  .  .  .  .  .  .  .  .  .  NamePos: 6:2
    48  .  .  .  .  .  .  .  .  .  Name: "fmt"
    49  .  .  .  .  .  .  .  .  }
    50  .  .  .  .  .  .  .  .  Sel: *ast.Ident {
    51  .  .  .  .  .  .  .  .  .  NamePos: 6:6
    52  .  .  .  .  .  .  .  .  .  Name: "Println"
    53  .  .  .  .  .  .  .  .  }
    54  .  .  .  .  .  .  .  }
    55  .  .  .  .  .  .  .  Lparen: 6:13
    56  .  .  .  .  .  .  .  Args: []ast.Expr (len = 1) {
    57  .  .  .  .  .  .  .  .  0: *ast.BasicLit {
    58  .  .  .  .  .  .  .  .  .  ValuePos: 6:14
    59  .  .  .  .  .  .  .  .  .  Kind: STRING
    60  .  .  .  .  .  .  .  .  .  Value: "\"Hello, world!\""
    61  .  .  .  .  .  .  .  .  }
    62  .  .  .  .  .  .  .  }
    63  .  .  .  .  .  .  .  Ellipsis: -
    64  .  .  .  .  .  .  .  Rparen: 6:29
    65  .  .  .  .  .  .  }
    66  .  .  .  .  .  }
    67  .  .  .  .  }
    68  .  .  .  .  Rbrace: 7:1
    69  .  .  .  }
    70  .  .  }
    71  .  }
    ..  .  .. // Left out for brevity
    83  }
