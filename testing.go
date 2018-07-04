package dicon

import (
	"testing"

	"go/format"

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
