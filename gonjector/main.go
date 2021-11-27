package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

func main() {
	//load the code we want to inject first

	fset := token.NewFileSet()

	// parsing with comments will break go/ast , this is a known issue https://github.com/golang/go/issues/20744

	pkgs, err := parser.ParseDir(fset, "./target/", nil, parser.AllErrors)
	if err != nil {
		panic(err)
	}
	for _, pkg := range pkgs {
		for fileName, file := range pkg.Files {
			log.Println("Instrumenting:", fileName)

			//add log lib <= this step can also be automated
			astutil.AddImport(fset, file, "log")

			ast.Inspect(file, Instrument)

			//save the new generated code
			//ioutil.WriteFile(fileName, buf.Bytes(), 0664)
			printer.Fprint(os.Stdout, fset, file)

		}
	}
}

func Instrument(n ast.Node) bool {

	// Find Function Declaration Statements
	fn, ok := n.(*ast.FuncDecl)

	if !ok {
		return true
	}

	codeToInject := GenAst(fn.Name.Name)

	if fn.Body != nil {
		fn.Body.List = append(codeToInject, fn.Body.List...)
	}

	return true
}

func GenAstFromFile(fnName string) []ast.Stmt {
	fset := token.NewFileSet()

	injectable_ast, err := parser.ParseFile(fset, "inject.go", nil, parser.AllErrors)
	if err != nil {
		panic(err)
	}

	fn := injectable_ast.Decls[1].(*ast.FuncDecl)
	ast.Inspect(fn.Body, CleanAstPos)

	ast.Fprint(os.Stdout, fset, fn.Body.List, nil)
	return fn.Body.List
}

func CleanAstPos(n ast.Node) bool {
	switch n.(type) {
	case *ast.CallExpr:
		node := n.(*ast.CallExpr)
		node.Lparen = token.NoPos
		node.Rparen = token.NoPos
	case *ast.Ident:
		node := n.(*ast.Ident)
		node.NamePos = token.NoPos
	case *ast.BasicLit:
		node := n.(*ast.BasicLit)
		node.ValuePos = token.NoPos
	}
	return true
}

func GenAst(fnName string) []ast.Stmt {
	//for now we are gnerating the AST ourselfs , since the AST generated from source code is position dependent

	fnName = fmt.Sprintf("\"Calling: %s\"", fnName)
	return []ast.Stmt{
		&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "log",
					},
					Sel: &ast.Ident{
						Name: "Println",
					},
				},

				Args: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: fnName,
					},
				},
			},
		},
	}
}
