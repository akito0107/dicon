package dicon

import (
	"fmt"
	"log"
	"strings"
)

type ContainerGenerator struct {
	Generator
}

func NewContainerGenerator() *ContainerGenerator {
	return &ContainerGenerator{}
}

func (g *ContainerGenerator) Generate(it *InterfaceType, fs []FuncType) error {
	g.PackageName = it.PackageName
	g.AppendHeader(it)
	g.AppendStructDefs(it)
	g.AppendMethod(fs, "")
	return nil
}

func (g *ContainerGenerator) AppendStructDefs(it *InterfaceType) {
	g.Printf("type dicontainer struct {\n")
	g.Printf("store map[string]interface{}\n")
	g.Printf("}\n")
	g.Printf("func NewDIContainer() %s {\n", it.Name)
	g.Printf("return &dicontainer{\n")
	g.Printf("store: map[string]interface{}{},\n")
	g.Printf("}\n")
	g.Printf("}\n")
	g.Printf("\n")
}

func (g *ContainerGenerator) AppendMethod(funcs []FuncType, _ string) {
	for _, f := range funcs {
		g.Printf("func (d *dicontainer) %s()", f.Name)
		if len(f.ReturnTypes) != 2 {
			log.Fatalf("Must be (instance, error) signature but %v", f.ReturnTypes)
		}

		returnType := f.ReturnTypes[0]
		g.Printf("(%s, error) {\n", returnType.ConvertName(g.PackageName))

		g.Printf("if i, ok := d.store[\"%s\"]; ok {\n", f.Name)
		g.Printf("instance, ok := i.(%s)\n", returnType.ConvertName(g.PackageName))
		g.Printf("if !ok {\n")
		g.Printf("return nil, fmt.Errorf(\"invalid instance is cached %%v\", instance)\n")
		g.Printf("}\n")
		g.Printf("return instance, nil\n")
		g.Printf("}\n")

		dep := make([]string, 0, len(f.ArgumentTypes))
		for i, a := range f.ArgumentTypes {
			g.Printf("dep%d, err := d.%s()\n", i, a.SimpleName())
			g.Printf("if err != nil {\n")
			g.Printf("return nil, errors.Wrap(err, \"resolve %s failed at DICON\")\n", a.SimpleName())
			g.Printf("}\n")
			dep = append(dep, fmt.Sprintf("dep%d", i))
		}

		g.Printf("instance, err := %sNew%s(%s)\n", g.RelativePackageName(f.PackageName), f.Name, strings.Join(dep, ", "))
		g.Printf("if err != nil {\n")
		g.Printf("return nil, errors.Wrap(err, \"creation %s failed at DICON\")\n", f.Name)
		g.Printf("}\n")
		g.Printf("d.store[\"%s\"] = instance\n", f.Name)
		g.Printf("return instance, nil\n")
		g.Printf("}\n")
	}
}

func (g *ContainerGenerator) RelativePackageName(packageName string) string {
	if g.PackageName == packageName {
		return ""
	}
	return packageName + "."
}
