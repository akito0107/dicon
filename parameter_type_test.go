package dicon

import (
	"testing"

	"go/parser"

	"github.com/andreyvit/diff"
)

func Test_relativeSelectorName(t *testing.T) {
	type in struct {
		d string
		c string
		s string
	}
	cases := []struct {
		in  in
		out string
	}{
		{in{"pack", "test", ""}, "pack"},
		{in{"pack", "pack", ""}, ""},
		{in{"pack", "pack", ""}, ""},
		{in{"other", "pack", "pack"}, ""},
		{in{"pack", "test", "test"}, ""},
		{in{"pack", "test", "other"}, "other"},
	}

	for _, c := range cases {
		if r := relativeSelectorName(c.in.d, c.in.c, c.in.s); r != c.out {
			t.Errorf(diff.CharacterDiff(r, c.out))
		}
	}
}

func TestParameterTypeConvertName(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{"string", "string"},
		{"Sample", "pack.Sample"},
		{"other.Sample", "other.Sample"},
		{"current.Sample", "Sample"},
		{"[]string", "[]string"},
		{"*string", "*string"},
		{"[]Sample", "[]pack.Sample"},
		{"[]*Sample", "[]*pack.Sample"},
		{"*[]Sample", "*[]pack.Sample"},
		{"map[string]string", "map[string]string"},
		{"map[current.Sample]string", "map[Sample]string"},
		{"map[*Sample]string", "map[*pack.Sample]string"},
		{"map[*Sample]string", "map[*pack.Sample]string"},
		{"interface{}", "interface{}"},
		{"struct{}", "struct{}"},
		{"map[interface{}]struct{}", "map[interface{}]struct{}"},
		{"chan string", "chan string"},
		{"chan<- string", "chan<- string"},
		{"<-chan string", "<-chan string"},
		{"func() error", "func() error"},
		{"func(a, b int) error", "func(a00 int, a01 int) error"},
		{"func(a, b int) (int, int)", "func(a00 int, a01 int) (int, int)"},
		{"func(a Sample, b test.Sample, c current.Sample) int", "func(a0 pack.Sample, a1 test.Sample, a2 Sample) int"},
		{"func(fmt string, args ...interface{}) int", "func(a0 string, a1 ...interface{}) int"},
		{"func(opts ...current.Option) int", "func(a0 ...Option) int"},
		{"func(opts ...Option) int", "func(a0 ...pack.Option) int"},
		{"func(opts ...test.Option) int", "func(a0 ...test.Option) int"},
	}
	for _, c := range cases {
		ast, e := parser.ParseExpr(c.in)
		if e != nil {
			t.Fatal(e)
		}
		p := &ParameterType{DeclaredPackageName: "pack", src: ast}
		if act := p.ConvertName("current"); act != c.out {
			t.Errorf(diff.CharacterDiff(act, c.out))
		}
	}
}
