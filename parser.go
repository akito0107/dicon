package dicon

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/k0kubun/pp"
)

type InterfaceType struct {
	PackageName    string
	Comments       []string
	Name           string
	Funcs          []FuncType
	DependPackages []Package
}

func (i *InterfaceType) AggregateFuncName() []string {
	funcnames := make([]string, 0, len(i.Funcs))
	for _, fn := range i.Funcs {
		funcnames = append(funcnames, fn.Name)
	}

	return funcnames
}

type FuncType struct {
	ArgumentTypes []ParameterType
	ReturnTypes   []ParameterType
	PackageName   string
	Comments      []string
	Name          string
}

type Package struct {
	Name string
	Path string
}

func FindDicon(packageName string, filenames []string) (*InterfaceType, error) {
	result := make([]InterfaceType, 0, len(filenames))
	for _, filename := range filenames {
		its, err := findDicon(packageName, filename, nil, "+DICON")
		if err != nil {
			return nil, err
		}
		if len(its) > 1 {
			return nil, fmt.Errorf("DICON interface must be single, but %d", len(its))
		} else if len(its) == 1 {
			result = append(result, its[0])
		}
	}

	if len(result) == 0 {
		return nil, nil
	}

	return &result[0], nil
}

func FindConstructors(packageName string, filenames []string, funcnames []string) ([]FuncType, error) {
	var result []FuncType

	for _, f := range filenames {
		r, err := findConstructors(packageName, f, nil, funcnames)
		if err != nil {
			return nil, err
		}
		result = append(result, r...)
	}

	return result, nil
}

func FindDependencyInterfaces(packageName string, filenames []string, targetNames []string) ([]InterfaceType, error) {
	var result []InterfaceType

	for _, f := range filenames {
		r, err := parseDependencyFuncs(packageName, targetNames, f, nil)
		if err != nil {
			return nil, err
		}
		result = append(result, r...)
	}

	return result, nil
}

func findConstructors(packageName string, from string, src interface{}, funcnames []string) ([]FuncType, error) {
	f, err := parser.ParseFile(token.NewFileSet(), from, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	var funcs []FuncType

	ast.Inspect(f, func(n ast.Node) bool {
		fun, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}
		if fun.Type.Results == nil || (len(fun.Type.Results.List) != 1 && len(fun.Type.Results.List) != 2) {
			return true
		}
		resultType := fun.Type.Results.List[0]
		for _, name := range funcnames {
			s := fmt.Sprintf("New%s", name)
			if s == fun.Name.Name && name == fmt.Sprintf("%v", resultType.Type) {
				if len(fun.Type.Results.List) == 2 && fmt.Sprintf("%v", fun.Type.Results.List[1].Type) != "error" {
					pp.Println(fun.Type.Results.List[1])
					return true
				}
				args := make([]ParameterType, 0, len(fun.Type.Params.List))
				for _, p := range fun.Type.Params.List {
					args = append(args, *NewParameterType(packageName, p.Type))
				}
				returns := []ParameterType{*NewParameterType(packageName, resultType.Type)}
				if len(fun.Type.Results.List) == 2 {
					eType := fun.Type.Results.List[1]
					returns = append(returns, ParameterType{DeclaredPackageName: packageName, src: eType.Type})
				}

				funcs = append(funcs, FuncType{
					ArgumentTypes: args,
					ReturnTypes:   returns,
					Name:          name,
					PackageName:   packageName,
				})
			}
		}
		return true
	})

	return funcs, nil
}

func findDicon(packageName string, from string, src interface{}, annotation string) ([]InterfaceType, error) {
	f, err := parser.ParseFile(token.NewFileSet(), from, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	deps := getDependencies(f)

	var its []InterfaceType

	pkg := f.Name.Name
	ast.Inspect(f, func(n ast.Node) bool {
		g, ok := n.(*ast.GenDecl)
		if !ok || g.Tok != token.TYPE {
			return true
		}
		comments := findComments(g.Doc)
		if !isAnnotated(comments, annotation) {
			return true
		}
		it, ok := findInterface(packageName, g.Specs)
		if !ok {
			return true
		}
		it.Comments = comments
		it.PackageName = pkg
		it.DependPackages = deps
		its = append(its, *it)

		return true
	})

	return its, nil
}

func findComments(cs *ast.CommentGroup) []string {
	var res []string
	if cs == nil {
		return res
	}
	for _, c := range cs.List {
		t := strings.TrimSpace(strings.TrimLeft(c.Text, "//"))
		res = append(res, t)
	}
	return res
}

func isAnnotated(cs []string, annotation string) bool {
	for _, c := range cs {
		if strings.HasPrefix(string(c), annotation) {
			return true
		}
	}
	return false
}

func findInterface(packageName string, specs []ast.Spec) (*InterfaceType, bool) {
	it := &InterfaceType{}
	var funcs []FuncType

	for _, spec := range specs {
		t := spec.(*ast.TypeSpec)
		s, ok := t.Type.(*ast.InterfaceType)
		if !ok {
			return it, false
		}
		it.Name = t.Name.Name
		for _, m := range s.Methods.List {
			f, ok := m.Type.(*ast.FuncType)
			if !ok {
				continue
			}
			ft := &FuncType{}

			var params []ParameterType
			if f.Params != nil {
				for _, p := range f.Params.List {
					for range p.Names {
						params = append(params, *NewParameterType(packageName, p.Type))
					}
				}
			}
			ft.ArgumentTypes = params

			var returns []ParameterType
			if f.Results != nil {
				for _, r := range f.Results.List {
					if len(r.Names) == 0 {
						returns = append(returns, *NewParameterType(packageName, r.Type))
					} else {
						for range r.Names {
							returns = append(returns, *NewParameterType(packageName, r.Type))
						}
					}
				}
			}
			ft.ReturnTypes = returns

			for _, n := range m.Names {
				ft.Name = n.Name
			}

			funcs = append(funcs, *ft)
		}
	}
	it.Funcs = funcs
	return it, true
}

func parseDependencyFuncs(packagename string, targetNames []string, from string, src interface{}) ([]InterfaceType, error) {
	var res []InterfaceType
	f, err := parser.ParseFile(token.NewFileSet(), from, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	deps := getDependencies(f)

	ast.Inspect(f, func(n ast.Node) bool {
		g, ok := n.(*ast.GenDecl)
		if !ok || g.Tok != token.TYPE {
			return true
		}
		it, ok := findInterface(packagename, g.Specs)
		if !ok || !contains(it.Name, targetNames) {
			return true
		}
		it.DependPackages = deps
		res = append(res, *it)
		return true
	})
	return res, nil
}

func contains(s string, source []string) bool {
	for _, str := range source {
		if s == str {
			return true
		}
	}
	return false
}

func getDependencies(f *ast.File) []Package {
	deps := make([]Package, 0, len(f.Imports))
	for _, imp := range f.Imports {
		name := ""
		if imp.Name != nil {
			name = imp.Name.Name
		}
		deps = append(deps, Package{
			Name: name,
			Path: imp.Path.Value,
		})
	}
	return deps
}
