package main

import (
	"fmt"
	"go/importer"
	"go/types"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/bilal-bhatti/zipline/internal"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

func main() {
	zipline := internal.NewZipline()

	zipline.Start()

}

func importPackage(p *packages.Package) {
	pkg, err := importer.For("source", nil).Import(p.PkgPath)
	if err != nil {
		panic(err)
	}
	log.Println("imported package for inspection", pkg)
	for _, imp := range pkg.Imports() {
		log.Println("import", imp.Name())
	}
	scope := pkg.Scope()

	for _, name := range scope.Names() {
		obj := scope.Lookup(name)

		if tn, ok := obj.Type().(*types.Named); ok {
			fmt.Printf("%#v\n", tn.NumMethods())
			for i := 0; i < tn.NumMethods(); i++ {
				method := tn.Method(i)
				log.Println("method", method.FullName())
				sig := method.Type().(*types.Signature)
				log.Println("signature", sig.String(), sig.Recv())
			}
		} else {
			log.Println("something else", obj)
		}
	}
}

func isRestifyImport(path string) bool {
	const vendorPart = "vendor/"
	if i := strings.LastIndex(path, vendorPart); i != -1 && (i == 0 || path[i-1] == '/') {
		path = path[i+len(vendorPart):]
	}
	return path == "echo/spec"
}

func goSrcRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "Failed to get working directory")
	}

	return strings.TrimPrefix(wd, os.Getenv("GOPATH")+"/src/"), nil
}

func findPackages() ([]string, error) {
	// cuz no set impl
	dirs := make(map[string]string)

	err := filepath.Walk(".",
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(p, ".go") {

				if !strings.Contains(p, "cmd") && !strings.Contains(p, "handlers") {
					dir := path.Dir(p)
					dirs[dir] = dir
				}
			}
			return nil
		})

	if err != nil {
		return nil, err
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Println("failed to get working directory: ", err)
		panic(err)
	}

	log.Println("working directory", wd)
	goSrc, err := goSrcRoot()
	if err != nil {
		panic(err)
	}

	pkgs := []string{}
	for k := range dirs {
		pkgs = append(pkgs, path.Join(goSrc, k))
	}

	return pkgs, nil
}
