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

## Code
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
```

The main function is composed of three parts: the name, the declaration, and the body. The name is represented as an identifier with the value main. The declaration, specified by the Type field, would contain a list of parameters and return type if we had specified any. The body consists of a list of statements with all lines of our program, in this case only one.
The main function is composed of three parts: the name, the declaration, and the body. The name is represented as an identifier with the value main. The declaration, specified by the Type field, would contain a list of parameters and return type if we had specified any. The body consists of a list of statements with all lines of our program, in this case only one.

Our single fmt.Println statement consists of quite a few parts in the AST. The statement is an ExprStmt, which represents an expression, which can, for example, be a function call, as it is here, or it can be a literal, a binary operation (for example addition and subtraction), a unary operation (for instance negating a number) and many more. Anything that can be used in a function call’s arguments is an expression.

Our ExprStmt contains a CallExpr, which is our actual function call. This again includes several parts, most important of which are Fun and Args. Fun contains a reference to the function call, in this case, it is a SelectorExpr, because we select the Println identifier from the fmt package. However, in the AST it is not yet known to the compiler that fmt is a package, it could also be a variable in the AST.

Args contains a list of expressions which are the arguments to the function. In this case, we have passed a literal string to the function, so it is represented by a BasicLit with type STRING.

It is clear that we are able to deduce a lot from the AST. That means that we can also inspect the AST further and find for example all function calls in the file. To do so, we are going to use the Inspect function from the ast package. 


After the imports have been resolved and the types have been checked, we are certain the program is valid Go code and we can start the process of converting the AST to (pseudo) machine code.

The first step in this process is to convert the AST to a lower-level representation of the program, specifically into a Static Single Assignment (SSA) form. This intermediate representation is not the final machine code, but it does represent the final machine code a lot more. SSA has a set of properties that make it easier to apply optimizations, most important of which that a variable is always defined before it is used and each variable is assigned exactly once.

After the initial version of the SSA has been generated, a number of optimization passes will be applied. These optimizations are applied to certain pieces of code that can be made simpler or faster for the processor to execute. For example, dead code, such as if (false) { fmt.Println(“test”) } can be eliminated because this will never execute. Another example of an optimization is that certain nil checks can be removed because the compiler can prove that these will never false.
```go

package main

import "fmt"

func main() {
  fmt.Println(2)
}

```
As you can see, this program has only one function and one import. It will print 2 when run. However, this sample will suffice for looking at the SSA.

 * **NOTE:**

 Only the SSA for the main function will be shown, as that is the interesting part.

 To show the generated SSA, we will need to set the GOSSAFUNC environment variable to the function we would like to view the SSA of, in this case main. We will also need to pass the -S flag to the compiler, so it will print the code and create an HTML file. We will also compile the file for Linux 64-bit, to make sure the machine code will be equal to what you will be seeing here. So, to compile the file we will run:

 $GOSSAFUNC=main GOOS=linux GOARCH=amd64 go build -gcflags "-S" simple.go



It will print all SSA, but it will also generate a ssa.html file which is interactive so we will use that.
When you open ssa.html, a number of passes will be shown, most of which are collapsed. The start pass is the SSA that is generated from the AST; the lower pass converts the non-machine specific SSA to machine-specific SSA and genssa is the final generated machine code.

The start phase’s code will look like this:
```go

b1:
	v1  = InitMem <mem>
	v2  = SP <uintptr>
	v3  = SB <uintptr>
	v4  = ConstInterface <interface {}>
	v5  = ArrayMake1 <[1]interface {}> v4
	v6  = VarDef <mem> {.autotmp_0} v1
	v7  = LocalAddr <*[1]interface {}> {.autotmp_0} v2 v6
	v8  = Store <mem> {[1]interface {}} v7 v5 v6
	v9  = LocalAddr <*[1]interface {}> {.autotmp_0} v2 v8
	v10 = Addr <*uint8> {type.int} v3
	v11 = Addr <*int> {"".statictmp_0} v3
	v12 = IMake <interface {}> v10 v11
	v13 = NilCheck <void> v9 v8
	v14 = Const64 <int> [0]
	v15 = Const64 <int> [1]
	v16 = PtrIndex <*interface {}> v9 v14
	v17 = Store <mem> {interface {}} v16 v12 v8
	v18 = NilCheck <void> v9 v17
	v19 = IsSliceInBounds <bool> v14 v15
	v24 = OffPtr <*[]interface {}> [0] v2
	v28 = OffPtr <*int> [24] v2
If v19 → b2 b3 (likely) (line 6)

b2: ← b1
	v22 = Sub64 <int> v15 v14
	v23 = SliceMake <[]interface {}> v9 v22 v22
	v25 = Copy <mem> v17
	v26 = Store <mem> {[]interface {}} v24 v23 v25
	v27 = StaticCall <mem> {fmt.Println} [48] v26
	v29 = VarKill <mem> {.autotmp_0} v27
Ret v29 (line 7)

b3: ← b1
	v20 = Copy <mem> v17
	v21 = StaticCall <mem> {runtime.panicslice} v20
Exit v21 (line 6)


```

This simple program already generates quite a lot of SSA (35 lines in total). However, a lot of it is boilerplate and quite a lot can be eliminated (the final SSA version has 28 lines and the final machine code version has 18 lines).