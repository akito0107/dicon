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
		{in: in{d: "pack", c: "test", s: ""}, out: "pack"},
		{in: in{d: "pack", c: "pack", s: ""}, out: ""},
		{in: in{d: "pack", c: "pack", s: ""}, out: ""},
		{in: in{d: "other", c: "pack", s: "pack"}, out: ""},
		{in: in{d: "pack", c: "test", s: "test"}, out: ""},
		{in: in{d: "pack", c: "test", s: "other"}, out: "other"},
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
		{in: "string", out: "string"},
		{in: "Sample", out: "pack.Sample"},
		{in: "other.Sample", out: "other.Sample"},
		{in: "current.Sample", out: "Sample"},
		{in: "[]string", out: "[]string"},
		{in: "*string", out: "*string"},
		{in: "[]Sample", out: "[]pack.Sample"},
		{in: "[]*Sample", out: "[]*pack.Sample"},
		{in: "*[]Sample", out: "*[]pack.Sample"},
		{in: "map[string]string", out: "map[string]string"},
		{in: "map[current.Sample]string", out: "map[Sample]string"},
		{in: "map[*Sample]string", out: "map[*pack.Sample]string"},
		{in: "map[*Sample]string", out: "map[*pack.Sample]string"},
		{in: "interface{}", out: "interface{}"},
		{in: "struct{}", out: "struct{}"},
		{in: "map[interface{}]struct{}", out: "map[interface{}]struct{}"},
		{in: "chan string", out: "chan string"},
		{in: "chan<- string", out: "chan<- string"},
		{in: "<-chan string", out: "<-chan string"},
		{in: "func() error", out: "func() error"},
		{in: "func(a, b int) error", out: "func(a00 int, a01 int) error"},
		{in: "func(a, b int) (int, int)", out: "func(a00 int, a01 int) (int, int)"},
		{in: "func(a Sample, b test.Sample, c current.Sample) int", out: "func(a0 pack.Sample, a1 test.Sample, a2 Sample) int"},
		{in: "func(fmt string, args ...interface{}) int", out: "func(a0 string, a1 ...interface{}) int"},
		{in: "func(opts ...current.Option) int", out: "func(a0 ...Option) int"},
		{in: "func(opts ...Option) int", out: "func(a0 ...pack.Option) int"},
		{in: "func(opts ...test.Option) int", out: "func(a0 ...test.Option) int"},
	}
	for _, c := range cases {
		ast, e := parser.ParseExpr(c.in)
		if e != nil {
			t.Fatal(e)
		}
		p := NewParameterType("pack", ast)
		if act := p.ConvertName("current"); act != c.out {
			t.Errorf(diff.CharacterDiff(act, c.out))
		}
	}
}
