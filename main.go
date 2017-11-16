package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"strings"

	"github.com/akito0107/dicon/internal"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "dicon"
	app.Usage = "DICONtainer Generator"

	app.Commands = []cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "generate dicon_gen file",
			Action: func(c *cli.Context) error {
				pkgs := strings.Split(c.String("pkg"), ",")
				filename := c.String("out")
				d := c.Bool("dry-run")
				return run(pkgs, filename, d)
			},
			Flags: []cli.Flag{
				cli.StringFlag{Name: "pkg, p", Value: "", Usage: "target package(s)."},
				cli.StringFlag{Name: "out, o", Value: "dicon_gen", Usage: "output file name"},
				cli.BoolFlag{Name: "dry-run"},
			},
		},
	}

	app.Run(os.Args)
}

func run(pkgs []string, filename string, dry bool) error {
	var it *internal.InterfaceType

	for _, pkg := range pkgs {
		files, err := ioutil.ReadDir("./" + pkg)
		if err != nil {
			return err
		}
		filenames := make([]string, 0)
		for _, f := range files {
			filenames = append(filenames, f.Name())
		}
		pparser := internal.NewPackageParser(pkg)
		res, err := pparser.FindDicon(&filenames)
		if err != nil {
			return err
		}

		if res != nil {
			it = res
			break
		}
	}

	if it == nil {
		return fmt.Errorf("+DICON not found")
	}

	targetPkg := it.PackageName

	funcnames := make([]string, 0)
	for _, fn := range it.Funcs {
		funcnames = append(funcnames, fn.Name)
	}

	funcs := make([]internal.FuncType, 0)
	for _, pkg := range pkgs {
		files, err := ioutil.ReadDir("./" + pkg)
		if err != nil {
			return err
		}
		filenames := make([]string, 0)
		for _, f := range files {
			filenames = append(filenames, f.Name())
		}
		pparser := internal.NewPackageParser(pkg)
		ft, err := pparser.FindConstructors(&filenames, &funcnames)
		if err != nil {
			return err
		}
		funcs = append(funcs, *ft...)
	}

	g := internal.NewGenerator()

	if e := g.Generate(it, &funcs); e != nil {
		return e
	}

	b, e := g.Out()
	if e != nil {
		panic(e)
	}

	if dry {
		fmt.Printf("%s\n", *b)
		return nil
	}

	name := fmt.Sprintf("%s/%s.go", targetPkg, filename)
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		os.Remove(name)
	}

	out, e := os.Create(name)
	if e != nil {
		return e
	}
	defer out.Close()
	if _, e := out.Write(*b); e != nil {
		return e
	}
	return nil
}
