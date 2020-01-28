package internal

import (
	"bytes"
	"context"
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"os"
	"reflect"
	"strings"

	"golang.org/x/tools/go/packages"
)

func HasSpecBindings(fn *ast.FuncDecl) bool {
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

func isRestifyImport(path string) bool {
	const vendorPart = "vendor/"
	if i := strings.LastIndex(path, vendorPart); i != -1 && (i == 0 || path[i-1] == '/') {
		path = path[i+len(vendorPart):]
	}
	return path == "echo/zipline"
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

		if buildObj == nil || buildObj.Pkg() == nil || !isRestifyImport(buildObj.Pkg().Path()) || !isSupportedMethod(buildObj) {
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
