package internal

import (
	"errors"
	"go/ast"
	"go/types"
	"strings"

	"github.com/bilal-bhatti/zipline/internal/tokens"
	"golang.org/x/tools/go/packages"
)

type provider struct {
	pkgs  []*packages.Package
	known map[string]*tokens.TypeToken
	tmpl  types.Object
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

func (p provider) provideWithReturns(vt *tokens.TypeToken, retNames []string) (*funcToken, error) {
	for _, pkg := range p.pkgs {
		info := pkg.TypesInfo
		trace("scanning %s for type %s", pkg.PkgPath, vt.Signature)

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
				if vt.SameType(result.Type(), false) {

					rets := []*tokens.TypeToken{}

					for j := 0; j < sig.Results().Len(); j++ {
						ret := sig.Results().At(j)
						var token *tokens.TypeToken
						token = tokens.NewTypeToken(ret.Type().String(), retNames[j])
						p.known[token.Signature] = token
						rets = append(rets, token)
					}

					trace("found a match for %s with %s : %s", vt.Signature, pf.Name(), sig.String())
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
