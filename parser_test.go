package dicon

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

var TEST_FILE1 = `
package main

type Ex1 interface {
	Exec() error
}

// +DICON
type Ex2 interface {
	Exec() error
	Exec2(i int) (string, error)
}
`

func TestPackageParser_findDicon(t *testing.T) {
	its, err := findDicon("main", "/tmp/tmp.go", TEST_FILE1, "+DICON")
	if err != nil {
		t.Error(err)
	}
	if len(its) != 1 {
		t.Errorf("must 1 interface but %d\n", len(its))
	}
	funcs := its[0].Funcs
	if len(funcs) != 2 {
		t.Errorf("must 2 functions but %d\n", len(funcs))
	}
	for _, it := range its {

		f := it.Funcs[0]
		// test function names
		if f.Name != "Exec" && f.Name != "Exec2" {
			t.Errorf("function name must be Exec, but : %s", f.Name)
		}

		// test functions
		if f.Name == "Exec" {
			if len(f.ArgumentTypes) != 0 {
				t.Errorf("Exec ArgumentType must be blank but: %d", len(f.ArgumentTypes))
			}
			if got := f.ReturnTypes[0].SimpleName(); got != "error" {
				t.Errorf("Exec return type must be error, but : %s", got)
			}
		}

		if f.Name == "Exec2" {
			if f.ArgumentTypes[0].SimpleName() != "int" {
				t.Errorf("Exec2 argument type must be int, but : %s", f.Name)
			}
			if got := f.ReturnTypes[0].SimpleName(); got != "string" {
				t.Errorf("Exec2 return type must be string, but : %s", got)
			}
			if got := f.ReturnTypes[1].SimpleName(); got != "error" {
				t.Errorf("Exec2 return type must be error, but : %s", got)
			}
		}
	}
}

var TEST_COMPONENT = `
package di

type SampleComponent interface {
	Exec() error
}

type sampleComponent struct {
	dep Dependency
}

func NewSampleComponent(dep Dependency) SampleComponent {
	return &sampleComponent {
		dep: dep,
	}
}

func (s *sampleComponent) Exec() error {
	return nil
}
`

func TestPackageParser_FindConstructors(t *testing.T) {
	fs, _ := findConstructors("test", "/tmp/tmp.go", TEST_COMPONENT, []string{"SampleComponent"})
	if len(fs) != 1 {
		t.Fatalf("must be 1")
	}

	fun := fs[0]

	if len(fun.ReturnTypes) != 1 || fun.ReturnTypes[0].ConvertName("test") != "SampleComponent" {
		t.Errorf("return type: %v wrong", fun.ReturnTypes)
	}

	if len(fun.ArgumentTypes) != 1 || fun.ArgumentTypes[0].ConvertName("test") != "Dependency" {
		t.Errorf("arg type: %v wrong", fun.ArgumentTypes)
	}

	if fun.Name != "SampleComponent" {
		t.Errorf("func name is SampleComponent but %s", fun.Name)
	}

	if fun.PackageName != "test" {
		t.Errorf("package name is test but %s", fun.PackageName)
	}
}

var TEST_COMPONENT_ERRORS = `
package di

type SampleComponent interface {
	Exec() error
}

type sampleComponent struct {
	dep Dependency
}

func NewSampleComponent(dep Dependency) (SampleComponent, error) {
	return &sampleComponent {
		dep: dep,
	}, nil
}

func (s *sampleComponent) Exec() error {
	return nil
}
`

func TestPackageParser_FindConstructorsErrors(t *testing.T) {
	fs, _ := findConstructors("test", "/tmp/tmp.go", TEST_COMPONENT_ERRORS, []string{"SampleComponent"})
	if len(fs) != 1 {
		t.Fatalf("must be 1, but %d", len(fs))
	}

	fun := fs[0]

	if len(fun.ReturnTypes) != 2 || fun.ReturnTypes[0].SimpleName() != "SampleComponent" || fun.ReturnTypes[1].SimpleName() != "error" {
		t.Errorf("return type: %s, %s wrong", fun.ReturnTypes[0].SimpleName(), fun.ReturnTypes[1].SimpleName())
	}

	if len(fun.ArgumentTypes) != 1 || fun.ArgumentTypes[0].SimpleName() != "Dependency" {
		t.Errorf("arg type: %v wrong", fun.ArgumentTypes)
	}

	if fun.Name != "SampleComponent" {
		t.Errorf("func name is SampleComponent but %s", fun.Name)
	}

	if fun.PackageName != "test" {
		t.Errorf("package name is test but %s", fun.PackageName)
	}
}

var TEST_DEPENDENCY = `
package di

type Dependency interface {
	Run() error
}
type dependency struct {}

func NewDependency() Dependency {
	return &dependency{}
}

func (*dependency) Run() error {
	return nil
}
`

func TestPackageParer_parseDependencyFuncs(t *testing.T) {
	ds, _ := parseDependencyFuncs("test", []string{"Dependency"}, "/tmp/tmp.go", TEST_DEPENDENCY)
	if len(ds) != 1 {
		t.Fatalf("dependency function length myst be 1 but %d", len(ds))
	}

	if ds[0].Name != "Dependency" {
		t.Errorf("Dependency name must be Dependency but : %s", ds[0].Name)
	}

	if len(ds[0].Funcs) != 1 {
		t.Fatalf("Dependency func has only Run method")
	}

	if ds[0].Funcs[0].ReturnTypes[0].ConvertName("test") != "error" {
		t.Errorf("Return type must be error")
	}
}

func TestPackageParser_findInterface(t *testing.T) {
	parseSpecs := func(src string) []ast.Spec {
		f, err := parser.ParseFile(token.NewFileSet(), "", "package test\n"+src, parser.AllErrors)
		if err != nil {
			t.Fatal(err)
		}
		return f.Decls[0].(*ast.GenDecl).Specs
	}

	ts := []struct {
		specs       []ast.Spec
		packageName string

		expected *InterfaceType
	}{
		{
			specs: parseSpecs(`
type A interface {
	F(a, b int) (int, error)
}
`),
			packageName: "test",

			expected: &InterfaceType{
				Name: "A",
				Funcs: []FuncType{
					{
						Name: "F",
						ArgumentTypes: []ParameterType{
							{"", ast.NewIdent("int")},
							{"", ast.NewIdent("int")},
						},
						ReturnTypes: []ParameterType{
							{"", ast.NewIdent("int")},
							{"", ast.NewIdent("error")},
						},
					},
				},
			},
		},
		{
			specs: parseSpecs(`
type B interface {
	F(a int) error
}
`),
			packageName: "test",

			expected: &InterfaceType{
				Name: "B",
				Funcs: []FuncType{
					{
						Name: "F",
						ArgumentTypes: []ParameterType{
							{"", ast.NewIdent("int")},
						},
						ReturnTypes: []ParameterType{
							{"", ast.NewIdent("error")},
						},
					},
				},
			},
		},
		{
			specs: parseSpecs(`
type C interface {
	F() (w, h int)
}
`),
			packageName: "test",

			expected: &InterfaceType{
				Name: "C",
				Funcs: []FuncType{
					{
						Name:          "F",
						ArgumentTypes: []ParameterType{},
						ReturnTypes: []ParameterType{
							{"", ast.NewIdent("int")},
							{"", ast.NewIdent("int")},
						},
					},
				},
			},
		},
	}

	for _, tc := range ts {
		got, ok := findInterface(tc.packageName, tc.specs)
		if ok != (tc.expected != nil) {
			t.Errorf("unexpected result. expected: %v, but got: %v", tc.expected, got)
			continue
		}
		if !ok {
			continue
		}

		if got.Name != tc.expected.Name {
			t.Errorf("unexpected name. expected: %v, but got: %v", tc.expected.Name, got.Name)
		}
		if len(got.Funcs) != len(tc.expected.Funcs) {
			t.Errorf("unexpected len(Funcs). expected: %v, but got: %v", len(tc.expected.Funcs), len(got.Funcs))
		} else {
			for i := range got.Funcs {
				if got.Funcs[i].Name != tc.expected.Funcs[i].Name {
					t.Errorf("unexpected Funcs[%d].Name. expected: %v, but got: %v", i,
						tc.expected.Funcs[i].Name, got.Funcs[i].Name)
				}

				if len(got.Funcs[i].ArgumentTypes) != len(tc.expected.Funcs[i].ArgumentTypes) {
					t.Errorf("unexpected len(Funcs[%d].ArgumentTypes). expected: %v, but got: %v", i,
						len(tc.expected.Funcs[i].ArgumentTypes), len(got.Funcs[i].ArgumentTypes))
				} else {
					for j := range got.Funcs[i].ArgumentTypes {
						if got.Funcs[i].ArgumentTypes[j].SimpleName() != tc.expected.Funcs[i].ArgumentTypes[j].SimpleName() {
							t.Errorf("unexpected Funcs[%d].ArgumentTypes[%d]. expected: %v, but got: %v", i, j,
								tc.expected.Funcs[i].ArgumentTypes[j].SimpleName(),
								got.Funcs[i].ArgumentTypes[j].SimpleName())
						}
					}
				}
				if len(got.Funcs[i].ReturnTypes) != len(tc.expected.Funcs[i].ReturnTypes) {
					t.Errorf("unexpected len(Funcs[%d].ReturnTypes). expected: %v, but got: %v", i,
						len(tc.expected.Funcs[i].ReturnTypes), len(got.Funcs[i].ReturnTypes))
				} else {
					for j := range got.Funcs[i].ReturnTypes {
						if got.Funcs[i].ReturnTypes[j].SimpleName() != tc.expected.Funcs[i].ReturnTypes[j].SimpleName() {
							t.Errorf("unexpected Funcs[%d].ReturnTypes[%d]. expected: %v, but got: %v", i, j,
								tc.expected.Funcs[i].ReturnTypes[j].SimpleName(),
								got.Funcs[i].ReturnTypes[j].SimpleName())
						}
					}
				}
			}
		}
	}
}
