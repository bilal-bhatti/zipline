package internal

import (
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

func (p provider) varFor(vt *typeToken) (*typeToken, bool) {
	v, ok := p.known[vt.signature]
	return v, ok
}

func (p provider) provide(vt *typeToken) *funcToken {
	for _, pkg := range p.pkgs {
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

			// TODO: handle error return

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

			for i := 0; i < sig.Results().Len(); i++ {
				result := sig.Results().At(i)
				if vt.sameType(result.Type().String()) {

					rets := []*typeToken{}
					for j := 0; j < sig.Results().Len(); j++ {
						ret := sig.Results().At(j)
						token := newTypeToken("", ret.Type().String(), vt.varName())
						p.known[token.signature] = token
						rets = append(rets, token)
					}

					return &funcToken{
						cpkg:      pf.Pkg().Name(),
						signature: pf.FullName(),
						args:      args,
						rets:      rets,
					}
				}
			}
		}
	}

	return nil
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
