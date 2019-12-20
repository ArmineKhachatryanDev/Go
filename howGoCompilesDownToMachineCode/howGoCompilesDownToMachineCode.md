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


  We will create our source code string and initialize the scanner.Scanner struct which will scan our source code. We call Scan() as many times as we can and print the token’s position, type, and literal string until we reach the End of File (EOF) marker.


  When we run the program, it will print the following:

 
1:1   package "package"
1:9   IDENT   "main"
1:13  ;       "\n"
2:1   import  "import"
2:8   STRING  "\"fmt\""
2:13  ;       "\n"
3:1   func    "func"
3:6   IDENT   "main"
3:10  (       ""
3:11  )       ""
3:13  {       ""
4:3   IDENT   "fmt"
4:6   .       ""
4:7   IDENT   "Println"
4:14  (       ""
4:15  STRING  "\"Hello, world!\""
4:30  )       ""
4:31  ;       "\n"
5:1   }       ""
5:2   ;       "\n"
5:3   EOF     ""