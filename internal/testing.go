package internal

import (
	"testing"

	"go/format"

	"golang.org/x/tools/imports"
)

func pretty(t *testing.T, src []byte) []byte {
	dist, e := format.Source(src)

	if e != nil {
		t.Fatal(e)
	}
	return dist
}

func fixImports(t *testing.T, src []byte) []byte {
	dist, e := imports.Process("/tmp/tmp.go", src, &imports.Options{Comments: true})

	if e != nil {
		t.Fatal(e)
	}
	return dist
}
