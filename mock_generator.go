package dicon

import (
	"fmt"
	"go/ast"
	"strings"
)

type MockGenerator struct {
	Generator
}

func NewMockGenerator() *MockGenerator {
	return &MockGenerator{}
}

func (g *MockGenerator) Generate(it *InterfaceType, targets []InterfaceType) error {
	if g.PackageName == "" {
		g.PackageName = it.PackageName
	}
	g.AppendHeader(it)
	g.AppendImports(targets)

	for _, i := range targets {
		g.AppendMockStruct(&i)
	}
	return nil
}

func (g *MockGenerator) AppendImports(targets []InterfaceType) {
	g.Printf("import (\n")
	defer g.Printf(")\n")

	imported := make(map[string]struct{})
	for _, target := range targets {
		for _, dep := range target.DependPackages {
			if _, ok := imported[dep.Path]; !ok {
				g.Printf("%s %s\n", dep.Name, dep.Path)
				imported[dep.Path] = struct{}{}
			}
		}
	}
}

func (g *MockGenerator) AppendMockStruct(it *InterfaceType) {
	g.Printf("type %sMock struct {\n", it.Name)
	args := map[string][]string{}
	returns := map[string][]string{}

	for _, f := range it.Funcs {
		var ags []string
		for i, a := range f.ArgumentTypes {
			ags = append(ags, fmt.Sprintf("a%d %s", i, a.ConvertName(g.PackageName)))
		}
		args[f.Name] = ags

		var rets []string
		for _, r := range f.ReturnTypes {
			rets = append(rets, r.ConvertName(g.PackageName))
		}
		returns[f.Name] = rets
		g.Printf("%sMock func(%s)", f.Name, strings.Join(ags, ","))
		if len(rets) == 1 {
			g.Printf("%s", strings.Join(rets, ","))
		} else if len(rets) != 0 {
			g.Printf("(%s)", strings.Join(rets, ","))
		}
		g.Printf("\n")
	}

	g.Printf("}\n")
	g.Printf("\n")
	g.Printf("func New%sMock() *%sMock {\n", it.Name, it.Name)
	g.Printf("return &%sMock{}\n", it.Name)
	g.Printf("}\n")
	g.Printf("\n")

	for _, f := range it.Funcs {
		ags := args[f.Name]
		rets := returns[f.Name]

		g.Printf("func (mk *%sMock) %s(%s) ", it.Name, f.Name, strings.Join(ags, ","))
		if len(rets) == 1 {
			g.Printf("%s", rets[0])
		} else if len(rets) != 0 {
			g.Printf("(%s)", strings.Join(rets, ","))
		}
		g.Printf(" {\n")
		var a []string
		for i, at := range f.ArgumentTypes {
			if _, ok := at.src.(*ast.Ellipsis); ok {
				a = append(a, fmt.Sprintf("a%d...", i))
			} else {
				a = append(a, fmt.Sprintf("a%d", i))
			}
		}
		g.Printf("return mk.%sMock(%s)\n", f.Name, strings.Join(a, ","))
		g.Printf("}\n")
	}
}
