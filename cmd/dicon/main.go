package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"io"

	"github.com/akito0107/dicon"
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
				return generateContainer(pkgs, filename, d)
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
				return generateMock(distPackage, pkgs, filename, d)
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

func generateContainer(pkgs []string, filename string, dry bool) error {
	it, err := findDicon(pkgs)
	if err != nil {
		return err
	}

	funcnames := it.AggregateFuncName()

	var funcs []dicon.FuncType
	for _, pkg := range pkgs {
		pkgName, filenames, err := readAllFilenames(pkg)
		if err != nil {
			return err
		}
		ft, err := dicon.FindConstructors(pkgName, filenames, funcnames)
		if err != nil {
			return err
		}
		funcs = append(funcs, ft...)
	}

	if err := dicon.DetectCyclicDependency(funcs); err != nil {
		return err
	}

	g := dicon.NewContainerGenerator()
	targetPkg := it.PackageName
	if err := g.Generate(it, funcs); err != nil {
		return err
	}
	return writeFile(g, targetPkg, filename, dry)
}

func generateMock(distPackage string, pkgs []string, filename string, dry bool) error {
	it, err := findDicon(pkgs)
	if err != nil {
		return err
	}

	funcnames := it.AggregateFuncName()

	var mockTargets []dicon.InterfaceType
	for _, pkg := range pkgs {
		pkgName, filenames, err := readAllFilenames(pkg)
		if err != nil {
			return err
		}
		m, err := dicon.FindDependencyInterfaces(pkgName, filenames, funcnames)
		if err != nil {
			return err
		}
		mockTargets = append(mockTargets, m...)
	}

	g := dicon.NewMockGenerator()
	g.PackageName = distPackage
	if err := g.Generate(it, mockTargets); err != nil {
		return err
	}
	return writeFile(g, distPackage, filename, dry)
}

func findDicon(pkgs []string) (*dicon.InterfaceType, error) {
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
		res, err := dicon.FindDicon(pkgName, filenames)
		if err != nil {
			return nil, err
		}

		if res != nil {
			return res, nil
		}
	}

	return nil, fmt.Errorf("+DICON not found")
}

func readAllFilenames(pkg string) (string, []string, error) {
	pkgDir := filepath.Join(".", filepath.FromSlash(pkg))
	files, err := ioutil.ReadDir(pkgDir)
	if err != nil {
		return "", nil, err
	}
	filenames := make([]string, 0, len(files))
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".go") {
			filenames = append(filenames, filepath.Join(pkgDir, f.Name()))
		}
	}
	_, pkgName := filepath.Split(pkgDir)
	return pkgName, filenames, nil
}

func writeFile(g dicon.Outer, targetPkg string, filename string, dry bool) error {
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
