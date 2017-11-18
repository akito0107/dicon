package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(pkgs []string, filename string, dry bool) error {
	var it *internal.InterfaceType

	for _, pkg := range pkgs {
		files, err := ioutil.ReadDir(filepath.Join(".", pkg))
		if err != nil {
			return err
		}
		filenames := make([]string, 0, len(files))
		for _, f := range files {
			filenames = append(filenames, f.Name())
		}
		pparser := internal.NewPackageParser(pkg)
		res, err := pparser.FindDicon(filenames)
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

	funcnames := make([]string, 0, len(it.Funcs))
	for _, fn := range it.Funcs {
		funcnames = append(funcnames, fn.Name)
	}

	var funcs []internal.FuncType
	for _, pkg := range pkgs {
		files, err := ioutil.ReadDir(filepath.Join(".", pkg))
		if err != nil {
			return err
		}
		filenames := make([]string, 0, len(files))
		for _, f := range files {
			filenames = append(filenames, f.Name())
		}
		pparser := internal.NewPackageParser(pkg)
		ft, err := pparser.FindConstructors(filenames, funcnames)
		if err != nil {
			return err
		}
		funcs = append(funcs, ft...)
	}

	g := internal.NewGenerator()

	if err := g.Generate(it, funcs); err != nil {
		return err
	}

	name := filepath.Join(targetPkg, filename+".go")
	var w io.Writer
	if dry {
		w = os.Stdout
	} else {
		if _, err := os.Stat(name); !os.IsNotExist(err) {
			os.Remove(name)
		}
		out, err := os.Create(name)
		if err != nil {
			return err
		}
		defer out.Close()
		w = out
	}
	if err := g.Out(w, name); err != nil {
		return err
	}

	return nil
}
