package dicon

import (
	"go/ast"
	"testing"
)

func TestDetectCyclicDependency(t *testing.T) {
	ts := []struct {
		Funcs      []FuncType
		RaiseError bool
	}{
		{
			Funcs: []FuncType{
				{
					ArgumentTypes: []ParameterType{
						{src: ast.NewIdent("B")},
						{src: ast.NewIdent("C")},
					},
					ReturnTypes: []ParameterType{
						{src: ast.NewIdent("A")},
					},
				},
				{
					ArgumentTypes: []ParameterType{
						{src: ast.NewIdent("C")},
					},
					ReturnTypes: []ParameterType{
						{src: ast.NewIdent("B")},
					},
				},
			},
			RaiseError: false,
		},
		{
			Funcs: []FuncType{
				{
					ArgumentTypes: []ParameterType{
						{src: ast.NewIdent("B")},
						{src: ast.NewIdent("C")},
					},
					ReturnTypes: []ParameterType{
						{src: ast.NewIdent("A")},
					},
				},
				{
					ArgumentTypes: []ParameterType{
						{src: ast.NewIdent("C")},
					},
					ReturnTypes: []ParameterType{
						{src: ast.NewIdent("B")},
					},
				},
				{
					ArgumentTypes: []ParameterType{
						{src: ast.NewIdent("A")},
					},
					ReturnTypes: []ParameterType{
						{src: ast.NewIdent("C")},
					},
				},
			},
			RaiseError: true,
		},
	}

	for _, tc := range ts {
		got := DetectCyclicDependency(tc.Funcs)
		if tc.RaiseError != (got != nil) {
			t.Errorf("unexpected error. expected: %v, but got: %v", tc.RaiseError, got)
		}
	}
}
