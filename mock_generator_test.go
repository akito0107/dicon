package dicon

import (
	"bytes"
	"testing"

	"go/ast"

	"github.com/andreyvit/diff"
)

func TestMockGenerator_AppendMockStruct(t *testing.T) {
	t.Run("Simple case", func(t *testing.T) {

		ex := pretty(t, []byte(`type TestInterfaceMock struct {
		TestFuncMock func(a0 Arg1, a1 Arg2) Ret1
	}

	func NewTestInterfaceMock() *TestInterfaceMock {
		return &TestInterfaceMock{}
	}

	func (mk *TestInterfaceMock) TestFunc(a0 Arg1, a1 Arg2) Ret1 {
		return mk.TestFuncMock(a0, a1)
	}
`))

		p1 := ParameterType{
			DeclaredPackageName: "test",
			src:                 createAst(t, "Arg1"),
		}
		p2 := ParameterType{
			DeclaredPackageName: "test",
			src:                 createAst(t, "Arg2"),
		}
		r1 := ParameterType{
			DeclaredPackageName: "test",
			src:                 createAst(t, "Ret1"),
		}

		f1 := FuncType{
			Name:          "TestFunc",
			ArgumentTypes: []ParameterType{p1, p2},
			ReturnTypes:   []ParameterType{r1},
		}

		it := &InterfaceType{
			Name:        "TestInterface",
			PackageName: "test",
			Funcs:       []FuncType{f1},
		}

		g := MockGenerator{
			Generator{
				PackageName: "test",
			},
		}
		g.AppendMockStruct(it)
		act := pretty(t, g.buf.Bytes())
		if !bytes.Equal(act, ex) {
			t.Errorf("Not Matched: \n%v", diff.LineDiff(string(ex), string(act)))
		}
	})

	t.Run("Multiple Funcs", func(t *testing.T) {
		ex := pretty(t, []byte(`type TestInterfaceMock struct {
		TestFunc1Mock func(a0 Arg1)
		TestFunc2Mock func(a0 Arg1, a1 Arg2) (Ret1, Ret2)
	}

	func NewTestInterfaceMock() *TestInterfaceMock {
		return &TestInterfaceMock{}
	}

	func (mk *TestInterfaceMock) TestFunc1(a0 Arg1) {
		return mk.TestFunc1Mock(a0)
	}
	func (mk *TestInterfaceMock) TestFunc2(a0 Arg1, a1 Arg2) (Ret1, Ret2) {
		return mk.TestFunc2Mock(a0, a1)
	}
`))

		p1 := ParameterType{
			DeclaredPackageName: "test",
			src:                 createAst(t, "Arg1"),
		}
		p2 := ParameterType{
			DeclaredPackageName: "test",
			src:                 createAst(t, "Arg2"),
		}
		r1 := ParameterType{
			DeclaredPackageName: "test",
			src:                 createAst(t, "Ret1"),
		}
		r2 := ParameterType{
			DeclaredPackageName: "test",
			src:                 createAst(t, "Ret2"),
		}

		f1 := FuncType{
			Name:          "TestFunc1",
			ArgumentTypes: []ParameterType{p1},
			ReturnTypes:   []ParameterType{},
		}

		f2 := FuncType{
			Name:          "TestFunc2",
			ArgumentTypes: []ParameterType{p1, p2},
			ReturnTypes:   []ParameterType{r1, r2},
		}

		it := &InterfaceType{
			Name:        "TestInterface",
			PackageName: "test",
			Funcs:       []FuncType{f1, f2},
		}

		g := MockGenerator{
			Generator{
				PackageName: "test",
			},
		}
		g.AppendMockStruct(it)
		act := pretty(t, g.buf.Bytes())
		if !bytes.Equal(act, ex) {
			t.Errorf("Not Matched: \n%v", diff.LineDiff(string(ex), string(act)))
		}
	})

	t.Run("Mutiple Funcs With Package", func(t *testing.T) {
		ex := pretty(t, []byte(`type TestInterfaceMock struct {
		TestFunc1Mock func(a0 pak1.Arg1)
		TestFunc2Mock func(a0 pak1.Arg1, a1 Arg2) (Ret1, pak2.Ret2)
	}

	func NewTestInterfaceMock() *TestInterfaceMock {
		return &TestInterfaceMock{}
	}

	func (mk *TestInterfaceMock) TestFunc1(a0 pak1.Arg1) {
		return mk.TestFunc1Mock(a0)
	}
	func (mk *TestInterfaceMock) TestFunc2(a0 pak1.Arg1, a1 Arg2) (Ret1, pak2.Ret2) {
		return mk.TestFunc2Mock(a0, a1)
	}
`))

		p1 := ParameterType{
			DeclaredPackageName: "test",
			src:                 createAst(t, "pak1.Arg1"),
		}
		p2 := ParameterType{
			DeclaredPackageName: "test",
			src:                 createAst(t, "Arg2"),
		}
		r1 := ParameterType{
			DeclaredPackageName: "test",
			src:                 createAst(t, "Ret1"),
		}
		r2 := ParameterType{
			DeclaredPackageName: "test",
			src:                 createAst(t, "pak2.Ret2"),
		}

		f1 := FuncType{
			Name:          "TestFunc1",
			ArgumentTypes: []ParameterType{p1},
			ReturnTypes:   []ParameterType{},
		}

		f2 := FuncType{
			Name:          "TestFunc2",
			ArgumentTypes: []ParameterType{p1, p2},
			ReturnTypes:   []ParameterType{r1, r2},
		}

		it := &InterfaceType{
			Name:        "TestInterface",
			PackageName: "test",
			Funcs:       []FuncType{f1, f2},
		}

		g := MockGenerator{
			Generator{
				PackageName: "test",
			},
		}
		g.AppendMockStruct(it)
		act := pretty(t, g.buf.Bytes())
		if !bytes.Equal(act, ex) {
			t.Errorf("Not Matched: \n%v", diff.LineDiff(string(ex), string(act)))
		}
	})

	t.Run("Variadic Arguments", func(t *testing.T) {
		ex := pretty(t, []byte(`type TestInterfaceMock struct {
		TestFunc1Mock func(a0 string, a1 ...interface{}) error
	}

	func NewTestInterfaceMock() *TestInterfaceMock {
		return &TestInterfaceMock{}
	}

	func (mk *TestInterfaceMock) TestFunc1(a0 string, a1 ...interface{}) error {
		return mk.TestFunc1Mock(a0, a1...)
	}
`))

		p1 := ParameterType{
			src: createAst(t, "string"),
		}
		p2 := ParameterType{
			src: &ast.Ellipsis{
				Elt: &ast.InterfaceType{},
			},
		}
		r1 := ParameterType{
			src: createAst(t, "error"),
		}

		f1 := FuncType{
			Name:          "TestFunc1",
			ArgumentTypes: []ParameterType{p1, p2},
			ReturnTypes:   []ParameterType{r1},
		}

		it := &InterfaceType{
			Name:        "TestInterface",
			PackageName: "test",
			Funcs:       []FuncType{f1},
		}

		g := MockGenerator{
			Generator{
				PackageName: "test",
			},
		}
		g.AppendMockStruct(it)
		act := pretty(t, g.buf.Bytes())
		if !bytes.Equal(act, ex) {
			t.Errorf("Not Matched: \n%v", diff.LineDiff(string(ex), string(act)))
		}
	})
}
