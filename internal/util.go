package internal

import (
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/importer"
	"go/token"
	"go/types"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

func hasSpecBindings(fn *ast.FuncDecl) bool {
	foundit := false
	ast.Inspect(fn, func(n ast.Node) bool {

		expStmt, ok := n.(*ast.ExprStmt)
		if !ok {
			return true
		}
		callExp, ok := expStmt.X.(*ast.CallExpr)
		if !ok {
			return true
		}

		// skip everything other than selector expression
		selExp, ok := callExp.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		// check the receiver is the package spec
		if x, ok := selExp.X.(*ast.Ident); !ok {
			return true
		} else {
			if x.Name != "spec" {
				return true
			}
		}

		if selExp.Sel.Name != "Post" && selExp.Sel.Name != "Get" {
			return true
		}

		// verify arguments
		if len(callExp.Args) != 2 {
			return true
		}

		if bl, ok := callExp.Args[0].(*ast.BasicLit); !ok || bl.Kind != token.STRING {
			return true
		}

		if _, ok := callExp.Args[1].(*ast.SelectorExpr); !ok {
			return true
		}

		// found what we were looking for
		foundit = true
		return false
	})
	return foundit
}

func isSupportedMethod(obj types.Object) bool {
	switch obj.Name() {
	case "Post", "Get":
		return true
	default:
		return false
	}
}

func isZiplineImport(path string) bool {
	const vendorPart = "vendor/"
	if i := strings.LastIndex(path, vendorPart); i != -1 && (i == 0 || path[i-1] == '/') {
		path = path[i+len(vendorPart):]
	}
	return path == "github.com/bilal-bhatti/zipline/pkg"
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

func print(msg string, n ast.Node) {
	var buf bytes.Buffer
	printNode(&buf, n)
	log.Println(msg, string(buf.Bytes()))
}

func printNode(buff *bytes.Buffer, n ast.Node) {
	switch t := n.(type) {
	case *ast.Ident:
		buff.WriteString(t.Name)
	case *ast.AssignStmt:
		for i := 0; i < len(t.Lhs); i++ {
			lhs := t.Lhs[i]
			printNode(buff, lhs)
			if i+1 < len(t.Lhs) {
				buff.WriteString(", ")
			}
		}
		buff.WriteString(" ")
		buff.WriteString(t.Tok.String())
		buff.WriteString(" ")
		for i := 0; i < len(t.Rhs); i++ {
			rhs := t.Rhs[i]
			printNode(buff, rhs)
			if i+1 < len(t.Rhs) {
				buff.WriteString(", ")
			}
		}
	case *ast.CallExpr:
		printNode(buff, t.Fun)
		buff.WriteString("(")
		for i := 0; i < len(t.Args); i++ {
			arg := t.Args[i]
			printNode(buff, arg)
			if i+1 < len(t.Args) {
				buff.WriteString(", ")
			}
		}
		buff.WriteString(")")
	case *ast.ExprStmt:
		printNode(buff, t.X)
	case *ast.SelectorExpr:
		printNode(buff, t.X)
		buff.WriteString(".")
		printNode(buff, t.Sel)
	case *ast.BasicLit:
		buff.WriteString(t.Value)
	case *ast.ReturnStmt:
		buff.WriteString("return ")
		for i := 0; i < len(t.Results); i++ {
			r := t.Results[i]
			printNode(buff, r)
			if i+1 < len(t.Results) {
				buff.WriteString(", ")
			}
		}
	default:
		log.Println("Unhandled node type", reflect.TypeOf(t))
	}
}

func load() []*packages.Package {
	wd, err := os.Getwd()
	if err != nil {
		log.Println("failed to get working directory: ", err)
		panic(err)
	}

	log.Println("Working directory", wd)

	cfg := &packages.Config{
		Context:    context.Background(),
		Mode:       packages.LoadAllSyntax,
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

func isBindingSpecNode(info *types.Info, fn ast.Node) bool {
	foundit := false

	ast.Inspect(fn, func(n ast.Node) bool {
		// skip everything other than call expression

		callExp, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		selExp, ok := callExp.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		// check the receiver is the package spec
		if x, ok := selExp.X.(*ast.Ident); !ok {
			return true
		} else {
			if x.Name != "zipline" {
				return true
			}
		}

		// // verify arguments
		if len(callExp.Args) != 1 {
			return true
		}

		buildObj := qualifiedIdentObject(info, callExp.Fun)

		if buildObj == nil || buildObj.Pkg() == nil || !isZiplineImport(buildObj.Pkg().Path()) || !isSupportedMethod(buildObj) {
			return true
		}

		// found what we were looking for
		foundit = true
		if foundit {
			return false
		}
		return true
	})
	return foundit
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
