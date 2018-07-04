package dicon

import (
	"testing"

	"go/ast"
	"go/format"
	"go/parser"

	"fmt"

	"golang.org/x/tools/imports"
)

func pretty(t *testing.T, src []byte) []byte {
	dist, err := format.Source(src)

	if err != nil {
		fmt.Printf("%s\n", src)
		t.Fatal(err)
	}
	return dist
}

func fixImports(t *testing.T, src []byte) []byte {
	dist, err := imports.Process("/tmp/tmp.go", src, &imports.Options{Comments: true})

	if err != nil {
		fmt.Printf("%s\n", src)
		t.Fatal(err)
	}
	return dist
}

func createAst(t *testing.T, expr string) ast.Expr {
	t.Helper()
	ex, err := parser.ParseExpr(expr)
	if err != nil {
		t.Fatal(ex)
		return nil
	}
	return ex
}
