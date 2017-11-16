package internal

import "testing"

var TEST_FILE1 = `
package main

type Ex1 interface {
	Exec() error
}

// +DICON
type Ex2 interface {
	Exec() error
	Exec2(i int) string
}
`

// must retrieve annotated interface
func TestParse(t *testing.T) {
	its, e := findDicon("/tmp/tmp.go", TEST_FILE1, "+DICON")
	if e != nil {
		t.Error(e)
	}

	if len(its) != 1 {
		t.Errorf("must 1 interface but %d\n", len(its))
	}
	funcs := its[0].Funcs
	if len(funcs) != 2 {
		t.Errorf("must 2 functions but %d\n", len(funcs))
	}
}

func TestParseWithFuncNames(t *testing.T) {
	its, _ := findDicon("/tmp/tmp.go", TEST_FILE1, "+DICON")
	for _, it := range its {
		for _, f := range it.Funcs {
			if f.Name != "Exec" && f.Name != "Exec2" {
				t.Errorf("function name must be Exec, but : %s", f.Name)
			}
		}
	}
}

func TestParseWithFuncParameters(t *testing.T) {
	its, _ := findDicon("/tmp/tmp.go", TEST_FILE1, "+DICON")
	for _, it := range its {
		f := it.Funcs[0]
		if f.Name == "Exec2" {
			if f.ArgumentTypes[0].Type != "int" {
				t.Errorf("argument type must be int, but : %s", f.Name)
			}
		}
	}
}

func TestParseWithReturnParameters(t *testing.T) {
	its, _ := findDicon("/tmp/tmp.go", TEST_FILE1, "+DICON")
	for _, it := range its {
		f := it.Funcs[0]
		if f.Name == "Exec" {
			if got := f.ReturnTypes[0].Type; got != "error" {
				t.Errorf("return type must be error, but : %s", got)
			}
		}
		if f.Name == "Exec2" {
			if got := f.ReturnTypes[0].Type; got != "string" {
				t.Errorf("return type must be string, but : %s", got)
			}
		}
	}
}

var TEST_DICON = `
package di

// +DICON
type DIContainer interface {
	SampleComponent() SampleComponent
}
`

var TEST_DEPENDENCY = `
package di

type Dependency interface {}
type dependency struct {}

func NewDependency() Dependency {
	return &dependency{}
}
`

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

	if len(fun.ReturnTypes) != 1 || fun.ReturnTypes[0].Type != "SampleComponent" {
		t.Errorf("return type: %v wrong", fun.ReturnTypes)
	}

	if len(fun.ArgumentTypes) != 1 || fun.ArgumentTypes[0].Type != "Dependency" {
		t.Errorf("arg type: %v wrong", fun.ArgumentTypes)
	}

	if fun.Name != "SampleComponent" {
		t.Errorf("func name is SampleComponent but %s", fun.Name)
	}

	if fun.PackageName != "test" {
		t.Errorf("package name is test but %s", fun.PackageName)
	}
}
