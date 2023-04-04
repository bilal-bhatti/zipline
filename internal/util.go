package internal

import (
	"context"
	"fmt"
	"go/ast"
	"go/types"
	"log"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/bilal-bhatti/zipline/internal/debug"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

func load(ps []string) ([]*packages.Package, error) {
	wd, err := os.Getwd()
	if err != nil {
		log.Println("Failed to get working directory: ", err)
		return nil, err
	}

	cfg := &packages.Config{
		Context:    context.Background(),
		Mode:       packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo,
		Dir:        wd,
		Env:        os.Environ(),
		BuildFlags: []string{"-tags=ziplinegen"},
	}

	pkgs, err := packages.Load(cfg, ps...)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load packages")
	}

	var errs []error
	for _, p := range pkgs {
		for _, e := range p.Errors {
			errs = append(errs, e)
		}
	}

	if len(errs) > 0 {
		var erm string
		for _, e := range errs {
			erm = fmt.Sprintf("%s\n", e.Error())
		}
		return nil, errors.New(erm)
	}

	// load imported package with same root path as well
	// if already not loaded
	var pkgSet = make(map[string]*packages.Package)
	var importSet = make(map[string]*packages.Package)
	for _, pkg := range pkgs {
		pkgSet[pkg.PkgPath] = pkg
	}

	for _, pkg := range pkgs {
		path := strings.Split(pkg.PkgPath, "/")
		for _, ipkg := range pkg.Imports {
			// already loaded ?
			if _, ok := pkgSet[ipkg.PkgPath]; ok {
				// skip it
				continue
			}
			if len(path) > 2 && strings.HasPrefix(ipkg.PkgPath, strings.Join(path[:2], "/")) {
				importSet[ipkg.PkgPath] = ipkg
			}
		}
	}

	var importList []string
	for key, _ := range importSet {
		importList = append(importList, key)
	}

	ipkgs, err := packages.Load(cfg, importList...)
	if err != nil {
		return nil, err
	}

	return append(pkgs, ipkgs...), nil
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

func findZiplineNodeAs(info *types.Info, fn ast.Node, targetType reflect.Type) (ast.Node, bool) {
	var found ast.Node
	var target ast.Node

	ast.Inspect(fn, func(n ast.Node) bool {
		if reflect.TypeOf(n) == targetType {
			target = n
		}

		if id, ok := n.(*ast.Ident); ok {
			ido := qualifiedIdentObject(info, id)

			if ido != nil && strings.HasSuffix(ido.Type().String(), ZiplineTemplate) {
				// search complete
				found = n
				return false
			}
		}

		// continue to inspect the tree
		return true
	})

	if found != nil {
		return target, true
	}

	return nil, false
}

// Useful to do a broad search for a zipline spec or a template
func isZiplineNode(info *types.Info, fn ast.Node) bool {
	var found ast.Node

	ast.Inspect(fn, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok {
			ido := qualifiedIdentObject(info, id)

			if ido != nil && strings.HasSuffix(ido.Type().String(), ZiplineTemplate) {
				// search complete
				found = n
				return false
			}
		}

		// continue to inspect the tree
		return true
	})
	return found != nil
}

func detectOutDir(pkg *packages.Package) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "failed to get working directory")
	}

	if len(pkg.GoFiles) == 0 {
		return "", errors.New("no files to derive output directory from package " + pkg.PkgPath)
	}

	out := path.Dir(pkg.GoFiles[0])

	debug.Trace("* calculating output package location *")
	debug.Trace("cwd: %s", cwd)
	debug.Trace("package path: %s", pkg.PkgPath)
	debug.Trace("package file: %s", out)

	return out, nil
}

func reverse(ss []string) {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}
