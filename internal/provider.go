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

func (p provider) typeTokenFor(vt *typeToken) (*typeToken, bool) {
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

						if vt.sameType(ret.Type().String()) {
							token = newTypeToken("", ret.Type().String(), vt.varName())
						} else {
							var varName = ""
							if ret.Type().String() == "error" {
								// TODO: fix - error name is always err
								varName = "err"
							}
							token = newTypeToken("", ret.Type().String(), varName)
						}

						if kt, ok := p.known[token.signature]; ok {
							rets = append(rets, kt)
						} else {
							p.known[token.signature] = token
							rets = append(rets, token)
						}
					}

					// adjust parameter names for known tokens
					for k := 0; k < len(args); k++ {
						arg := args[k]
						if kt, ok := p.known[arg.signature]; ok {
							args[k] = kt
						}
					}

					// log.Println("found provider for", vt.signature)
					// log.Println("known", p.known)
					// log.Println("rets", rets)
					// log.Println("args", args)
					return &funcToken{
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

func (p provider) provideWithReturns(vt *typeToken, retNames []string) *funcToken {
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

			if sig.Params().Len() != len(args) {
				// not a function with parameter match
				continue
			}

			for i := 0; i < sig.Results().Len(); i++ {
				result := sig.Results().At(i)
				if vt.sameType(result.Type().String()) {

					rets := []*typeToken{}
					if sig.Results().Len() != len(retNames) {
						panic("expected and provided return value count don't match")
					}

					for j := 0; j < sig.Results().Len(); j++ {
						ret := sig.Results().At(j)
						var token *typeToken
						token = newTypeToken("", ret.Type().String(), retNames[j])
						p.known[token.signature] = token
						rets = append(rets, token)
					}

					// log.Println("found provider for", vt.signature)
					// log.Println("known", p.known)
					// log.Println("rets", rets)
					// log.Println("args", args)
					return &funcToken{
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
