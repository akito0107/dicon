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

		if f.Name == "Exec" {
			if len(f.ArgumentTypes) != 0 {
				t.Errorf("Exec ArgumentType must be blank but: %d", len(f.ArgumentTypes))
			}
			if got := f.ReturnTypes[0].Type; got != "error" {
				t.Errorf("Exec return type must be error, but : %s", got)
			}
		}

		// test argument type
		if f.Name == "Exec2" {
			if f.ArgumentTypes[0].Type != "int" {
				t.Errorf("Exec2 argument type must be int, but : %s", f.Name)
			}
			if got := f.ReturnTypes[0].Type; got != "string" {
				t.Errorf("Exec2 return type must be string, but : %s", got)
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
