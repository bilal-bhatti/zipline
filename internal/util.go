package internal

import (
	"context"
	"go/ast"
	"go/types"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

var tr = false

func trace(f string, a ...interface{}) {
	if tr {
		log.Print(f, a)
	}
}

func load() []*packages.Package {
	wd, err := os.Getwd()
	if err != nil {
		log.Println("Failed to get working directory: ", err)
		panic(err)
	}

	cfg := &packages.Config{
		Context:    context.Background(),
		Mode:       packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo,
		Dir:        wd,
		Env:        os.Environ(),
		BuildFlags: []string{"-tags=ziplinegen"},
	}

	// Package pattern to search
	ps := []string{"./..."}
	// ps := pkgs(nil)

	pkgs, err := packages.Load(cfg, ps...)
	if err != nil {
		panic(err)
	}
	var errs []error
	for _, p := range pkgs {
		for _, e := range p.Errors {
			errs = append(errs, e)
		}
	}
	if len(errs) > 0 {
		log.Println(errs)
		panic(errs)
	}

	return pkgs
}

func qualifiedIdentObject(info *types.Info, expr ast.Expr) types.Object {
	switch expr := expr.(type) {
	case *ast.Ident:
		return info.ObjectOf(expr)
	case *ast.SelectorExpr:
		pkgName, ok := expr.X.(*ast.Ident)
		if !ok {
			return nil
		}
		if _, ok := info.ObjectOf(pkgName).(*types.PkgName); !ok {
			return nil
		}
		return info.ObjectOf(expr.Sel)
	default:
		return nil
	}
}

// Useful to do a broad search for a zipline spec or a template
func isZiplineNode(info *types.Info, fn ast.Node) bool {
	foundit := false
	ast.Inspect(fn, func(n ast.Node) bool {
		// ****************
		// returning true here, means we don't want to inspect the node
		// any further and want to move on to the next node
		// ****************

		// check the token ident is the ZiplineTemplate type
		id, ok := n.(*ast.Ident)
		if !ok {
			return true
		}

		ido := qualifiedIdentObject(info, id)

		if ido == nil {
			return true
		}

		// ensure receiver var type is ZiplineTemplate
		if !strings.HasSuffix(ido.Type().String(), ZiplineTemplate) {
			return true
		}

		// if we arrived here, then all previous checks passed and
		// found what we were looking for
		foundit = true
		if foundit {
			// returning false, signals stopping any further inspection
			return false
		}

		// keep inspecting
		return true
	})
	return foundit
}

func goSrcRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "Failed to get working directory")
	}

	return strings.TrimPrefix(wd, os.Getenv("GOPATH")+"/src/"), nil
}
