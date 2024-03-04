package internal

import (
	"errors"
	"go/ast"
	"go/types"

	"github.com/bilal-bhatti/zipline/internal/debug"

	"github.com/bilal-bhatti/zipline/internal/tokens"
	"golang.org/x/tools/go/packages"
)

type provider struct {
	pkgs  []*packages.Package
	known map[string]*tokens.TypeToken
	// tmpl  types.Object
}

func newProvider(pkgs []*packages.Package) *provider {
	known := make(map[string]*tokens.TypeToken)

	return &provider{
		pkgs:  pkgs,
		known: known,
	}
}

func (p *provider) memorize(id *ast.Ident) {
	obj := p.qualifiedIdentObject(id)
	tt := tokens.NewTypeToken(obj.Type().String(), obj.Name())
	p.known[tt.Signature] = tt
}

func (p *provider) memorizeTypeToken(t *tokens.TypeToken) {
	p.known[t.Signature] = t
}

func (p provider) typeTokenFor(vt *tokens.TypeToken) (*tokens.TypeToken, bool) {
	v, ok := p.known[vt.Signature]
	return v, ok
}

func (p provider) provide(vt *tokens.TypeToken, retNames []string) (*tokens.FuncToken, error) {
	debug.Trace("scanning packages for %s", vt.FullSignature)

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

			// return counts must match
			if sig.Results().Len() != len(retNames) {
				continue
			}

			args := []*tokens.TypeToken{}

			for i := 0; i < sig.Params().Len(); i++ {
				param := tokens.NewTypeToken(sig.Params().At(i).Type().String(), "")
				if arg, ok := p.typeTokenFor(param); !ok {
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
				if vt.SameType(result.Type(), false) {

					rets := []*tokens.TypeToken{}

					for j := 0; j < sig.Results().Len(); j++ {
						ret := sig.Results().At(j)
						token := tokens.NewTypeToken(ret.Type().String(), retNames[j])
						p.memorizeTypeToken(token)
						rets = append(rets, token)
					}

					debug.Trace("resolved %s with %s: %s", vt.FullSignature, pf.Name(), sig)
					return &tokens.FuncToken{
							Signature: pf.FullName(),
							Args:      args,
							Rets:      rets,
						},
						nil
				}
			}
		}
	}

	return nil, errors.New("no provider function found for " + vt.Signature)
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
