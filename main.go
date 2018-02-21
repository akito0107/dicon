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

var (
	version  = "master"
	revision string
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s version=%s revision=%s\n", c.App.Name, c.App.Version, revision)
	}

	app := cli.NewApp()

	app.Name = "dicon"
	app.Version = version
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
				return runGenerate(pkgs, filename, d)
			},
			Flags: []cli.Flag{
				cli.StringFlag{Name: "pkg, p", Value: "", Usage: "target package(s)."},
				cli.StringFlag{Name: "out, o", Value: "dicon_gen", Usage: "output file name"},
				cli.BoolFlag{Name: "dry-run"},
			},
		},
		{
			Name:    "generate-mock",
			Aliases: []string{"m"},
			Usage:   "generate dicon_mock file",
			Action: func(c *cli.Context) error {
				pkgs := strings.Split(c.String("pkg"), ",")
				filename := c.String("out")
				distPackage := c.String("dist")
				d := c.Bool("dry-run")
				return runGenerateMock(distPackage, pkgs, filename, d)
			},
			Flags: []cli.Flag{
				cli.StringFlag{Name: "pkg, p", Value: "", Usage: "target package(s)."},
				cli.StringFlag{Name: "out, o", Value: "dicon_mock", Usage: "output file name"},
				cli.StringFlag{Name: "dist, d", Value: "mock", Usage: "output package name"},
				cli.BoolFlag{Name: "dry-run"},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runGenerate(pkgs []string, filename string, dry bool) error {
	it, err := findDicon(pkgs)
	if err != nil {
		return err
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
		pkgDir := filepath.Join(".", filepath.FromSlash(pkg))
		files, err := ioutil.ReadDir(pkgDir)
		if err != nil {
			return err
		}
		filenames := make([]string, 0, len(files))
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".go") {
				filenames = append(filenames, filepath.Join(pkgDir, f.Name()))
			}
		}

		_, pkgName := filepath.Split(pkgDir)
		pparser := internal.NewPackageParser(pkgName)
		ft, err := pparser.FindConstructors(filenames, funcnames)
		if err != nil {
			return err
		}
		funcs = append(funcs, ft...)
	}

	if err := internal.DetectCyclicDependency(funcs); err != nil {
		return err
	}

	g := internal.NewGenerator()

	if err := g.Generate(it, funcs); err != nil {
		return err
	}

	return writeFile(g, targetPkg, filename, dry)
}

func runGenerateMock(distPackage string, pkgs []string, filename string, dry bool) error {
	it, err := findDicon(pkgs)
	if err != nil {
		return err
	}
	if it == nil {
		return fmt.Errorf("+DICON not found")
	}

	funcnames := make([]string, 0, len(it.Funcs))
	for _, fn := range it.Funcs {
		funcnames = append(funcnames, fn.Name)
	}

	var mockTargets []internal.InterfaceType
	for _, pkg := range pkgs {
		pkgDir := filepath.Join(".", filepath.FromSlash(pkg))
		files, err := ioutil.ReadDir(pkgDir)
		if err != nil {
			return err
		}
		filenames := make([]string, 0, len(files))
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".go") {
				filenames = append(filenames, filepath.Join(pkgDir, f.Name()))
			}
		}

		_, pkgName := filepath.Split(pkgDir)
		pparser := internal.NewPackageParser(pkgName)
		m, err := pparser.FindDependencyInterfaces(filenames, funcnames)
		if err != nil {
			return err
		}
		mockTargets = append(mockTargets, m...)
	}

	g := internal.NewGenerator()
	g.PackageName = distPackage
	if err := g.GenerateMock(it, mockTargets); err != nil {
		return err
	}
	return writeFile(g, distPackage, filename, dry)
}

func findDicon(pkgs []string) (*internal.InterfaceType, error) {
	var it *internal.InterfaceType
	for _, pkg := range pkgs {
		pkgDir := filepath.Join(".", filepath.FromSlash(pkg))

		files, err := ioutil.ReadDir(pkgDir)
		if err != nil {
			return nil, err
		}
		filenames := make([]string, 0, len(files))
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".go") {
				filenames = append(filenames, filepath.Join(pkgDir, f.Name()))
			}
		}

		_, pkgName := filepath.Split(pkgDir)
		pparser := internal.NewPackageParser(pkgName)
		res, err := pparser.FindDicon(filenames)
		if err != nil {
			return nil, err
		}

		if res != nil {
			it = res
			break
		}
	}
	return it, nil
}

func writeFile(g *internal.Generator, targetPkg string, filename string, dry bool) error {
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
