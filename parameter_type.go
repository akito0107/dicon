package dicon

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"regexp"
	"strings"
)

type ParameterType struct {
	DeclaredPackageName string
	src                 ast.Expr
}

func NewParameterType(packageName string, expr ast.Expr) *ParameterType {
	return &ParameterType{
		DeclaredPackageName: packageName,
		src:                 expr,
	}
}

func (p *ParameterType) ConvertName(packageName string) string {
	return convertName(p.DeclaredPackageName, packageName, p.src)
}

func (p *ParameterType) SimpleName() string {
	switch n := p.src.(type) {
	case *ast.SelectorExpr:
		return n.Sel.Name
	case *ast.Ident:
		return n.Name
	}
	panic("unreachable")
}

func convertName(declared, packageName string, expr ast.Expr) string {
	switch ex := expr.(type) {
	case *ast.Ident:
		name := ex.Name
		if isPrimitive(name) {
			return name
		}
		selector := relativeSelectorName(declared, packageName, "")
		return buildTypeName(selector, name)
	case *ast.SelectorExpr:
		selector := relativeSelectorName(declared, packageName, fmt.Sprintf("%v", ex.X))
		typ := ex.Sel.Name
		return buildTypeName(selector, typ)
	case *ast.ArrayType:
		return "[]" + convertName(declared, packageName, ex.Elt)
	case *ast.StarExpr:
		return "*" + convertName(declared, packageName, ex.X)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", convertName(declared, packageName, ex.Key), convertName(declared, packageName, ex.Value))
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.StructType:
		return "struct{}"
	case *ast.ChanType:
		var ch string
		if token.Pos(ex.Arrow) == token.NoPos {
			ch = "chan "
		} else if ex.Dir == ast.SEND {
			ch = "chan<- "
		} else if ex.Dir == ast.RECV {
			ch = "<-chan "
		} else {
			log.Fatalf("unknown chan type %+v", ex)
		}
		return ch + convertName(declared, packageName, ex.Value)
	case *ast.FuncType:
		var args []string
		for i, a := range ex.Params.List {
			if len(a.Names) == 1 {
				args = append(args, fmt.Sprintf("a%d %s", i, convertName(declared, packageName, a.Type)))
			} else {
				ty := convertName(declared, packageName, a.Type)
				for j := range a.Names {
					args = append(args, fmt.Sprintf("a%d%d %s", i, j, ty))
				}
			}
		}
		var rets []string
		for _, r := range ex.Results.List {
			rets = append(rets, convertName(declared, packageName, r.Type))
		}
		arg := "(" + strings.Join(args, ", ") + ")"
		var ret string
		if len(rets) > 1 {
			ret = "(" + strings.Join(rets, ", ") + ")"
		} else {
			ret = strings.Join(rets, "")
		}
		return "func" + arg + " " + ret
	case *ast.Ellipsis:
		return "..." + convertName(declared, packageName, ex.Elt)
	default:
		log.Fatalf("unsupported type %+v", expr)
	}
	return ""
}

var reg = regexp.MustCompile("^[a-z].*")

func isPrimitive(in string) bool {
	return reg.MatchString(in)
}

func relativeSelectorName(declared, current, selector string) string {
	if declared == current {
		return selector
	}
	if selector == current {
		return ""
	}
	if declared != current && selector == "" {
		return declared
	}
	if declared != current && selector != "" {
		return selector
	}
	return ""
}

func buildTypeName(selector, name string) string {
	if selector == "" {
		return name
	}
	return selector + "." + name
}
