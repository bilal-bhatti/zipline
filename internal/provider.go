package internal

import (
	"errors"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

type provider struct {
	pkgs  []*packages.Package
	known map[string]*typeToken
	tmpl  types.Object
}

func newProvider(pkgs []*packages.Package) *provider {
	known := make(map[string]*typeToken)

	return &provider{
		pkgs:  pkgs,
		known: known,
	}
}

func (p provider) typeTokenFor(vt *typeToken) (*typeToken, bool) {
	v, ok := p.known[vt.signature]
	return v, ok
}

func (p provider) provideWithReturns(vt *typeToken, retNames []string) (*funcToken, error) {
	trace("find provider for %s", vt.signature)

	var pkgsToScan []*packages.Package
	var pkgSet = make(map[string]*packages.Package)

	for _, pkg := range p.pkgs {
		pkgSet[pkg.PkgPath] = pkg

		// only scan packages that share the same root as binding package
		path := strings.Split(pkg.PkgPath, "/")
		if len(path) >= 2 {
			trace("filtering imported packages by path %s", strings.Join(path[:2], "/"))
		}
		for _, ipkg := range pkg.Imports {
			if len(path) > 2 && strings.HasPrefix(ipkg.PkgPath, strings.Join(path[:2], "/")) {
				pkgSet[ipkg.PkgPath] = ipkg
			}
		}
	}

	for _, p := range pkgSet {
		pkgsToScan = append(pkgsToScan, p)
	}

	for _, pkg := range pkgsToScan {
		info := pkg.TypesInfo
		for _, v := range info.Defs {
			pf, ok := v.(*types.Func)
			if !ok || !v.Exported() {
				continue
			}

			sig, ok := v.Type().(*types.Signature)
			if !ok {
				continue
			}

			// return counts must match
			if sig.Results().Len() != len(retNames) {
				continue
			}

			args := []*typeToken{}

			for i := 0; i < sig.Params().Len(); i++ {
				param := sig.Params().At(i)
				if arg, ok := p.known[strings.TrimPrefix(param.Type().String(), "*")]; !ok {
					args = args[:0]
					continue
				} else {
					args = append(args, arg)
				}
			}

			if sig.Params().Len() != len(args) {
				// not a function with parameter match
				continue
			}

			for i := 0; i < sig.Results().Len(); i++ {
				result := sig.Results().At(i)
				if vt.sameType(result.Type().String()) {

					rets := []*typeToken{}

					for j := 0; j < sig.Results().Len(); j++ {
						ret := sig.Results().At(j)
						var token *typeToken
						token = newTypeToken("", ret.Type().String(), retNames[j])
						p.known[token.signature] = token
						rets = append(rets, token)
					}

					trace("found a match for %s with %s : %s", vt.signature, pf.Name(), sig.String())
					return &funcToken{
							signature: pf.FullName(),
							args:      args,
							rets:      rets,
						},
						nil
				}
			}
		}
	}

	return nil, errors.New("no provider function found for " + vt.signature)
}

func (p provider) qualifiedIdentObject(expr ast.Expr) types.Object {
	for _, pkg := range p.pkgs {
		info := pkg.TypesInfo
		switch expr := expr.(type) {
		case *ast.Ident:
			obj := info.ObjectOf(expr)
			if obj != nil {
				return obj
			}
		case *ast.SelectorExpr:
			pkgName, ok := expr.X.(*ast.Ident)
			if !ok {
				continue
			}
			if _, ok := info.ObjectOf(pkgName).(*types.PkgName); !ok {
				continue
			}
			obj := info.ObjectOf(expr.Sel)
			if obj != nil {
				return obj
			}
		}
	}
	return nil
}
