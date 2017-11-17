package internal

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type PackageParser struct {
	PackageName string
}

type comments []comment
type comment string

type InterfaceType struct {
	PackageName string
	Comments    comments
	Name        string
	Funcs       []FuncType
}

type FuncType struct {
	ArgumentTypes []parameterType
	ReturnTypes   []parameterType
	PackageName   string
	Comments      comments
	Name          string
}

type parameterType struct {
	Type string
}

func NewPackageParser(pack string) *PackageParser {
	return &PackageParser{
		PackageName: pack,
	}
}

func (p *PackageParser) FindDicon(filenames []string) (*InterfaceType, error) {
	result := make([]InterfaceType, 0)
	for _, filename := range filenames {
		its, err := findDicon(p.PackageName+"/"+filename, nil, "+DICON")
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

	return &(result[0]), nil
}

func (p *PackageParser) FindConstructors(filenames []string, funcnames []string) ([]FuncType, error) {
	var result []FuncType

	for _, f := range filenames {
		r, err := findConstructors(p.PackageName, f, nil, funcnames)
		if err != nil {
			return nil, err
		}
		result = append(result, r...)
	}

	return result, nil
}

func findConstructors(packageName string, from string, src interface{}, funcnames []string) ([]FuncType, error) {
	f, err := parser.ParseFile(token.NewFileSet(), packageName+"/"+from, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	funcs := make([]FuncType, 0)

	ast.Inspect(f, func(n ast.Node) bool {
		fun, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}
		if fun.Type.Results == nil || len(fun.Type.Results.List) != 1 {
			return true
		}
		resultType := fun.Type.Results.List[0]
		for _, name := range funcnames {
			s := fmt.Sprintf("New%s", name)
			if s == fun.Name.Name && name == fmt.Sprintf("%v", resultType.Type) {
				args := make([]parameterType, 0, len(fun.Type.Params.List))
				for _, p := range fun.Type.Params.List {
					args = append(args, parameterType{
						Type: parseTypeExpr(p.Type),
					})
				}
				rt := parameterType{
					Type: parseTypeExpr(resultType.Type),
				}
				returns := []parameterType{rt}

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

func findDicon(from string, src interface{}, annotation string) ([]InterfaceType, error) {
	f, err := parser.ParseFile(token.NewFileSet(), from, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	its := make([]InterfaceType, 0)

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
		it, ok := findInterface(g.Specs)
		if !ok {
			return true
		}
		it.Comments = comments
		it.PackageName = pkg
		its = append(its, *it)

		return true
	})

	return its, nil
}

func findComments(cs *ast.CommentGroup) comments {
	res := comments{}
	if cs == nil {
		return res
	}
	for _, c := range cs.List {
		t := strings.TrimSpace(strings.TrimLeft(c.Text, "//"))
		res = append(res, comment(t))
	}
	return res
}

func isAnnotated(cs comments, annotation string) bool {
	for _, c := range cs {
		if strings.HasPrefix(string(c), annotation) {
			return true
		}
	}
	return false
}

func findInterface(specs []ast.Spec) (*InterfaceType, bool) {
	it := &InterfaceType{}
	funcs := make([]FuncType, 0)

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

			params := make([]parameterType, 0)
			for _, p := range f.Params.List {
				params = append(params, parameterType{
					Type: parseTypeExpr(p.Type),
				})
			}
			ft.ArgumentTypes = params

			returns := make([]parameterType, 0)
			for _, r := range f.Results.List {
				returns = append(returns, parameterType{
					Type: parseTypeExpr(r.Type),
				})
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

func parseTypeExpr(expr ast.Expr) string {
	ret := ""
	if idx, ok := expr.(*ast.SelectorExpr); ok {
		ret = fmt.Sprintf("%v", idx.Sel)
	} else {
		ret = fmt.Sprintf("%v", expr)
	}
	return ret
}
