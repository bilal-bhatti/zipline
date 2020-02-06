package internal

import (
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

type provider struct {
	pkgs  []*packages.Package
	known map[string]*varToken
}

func newProvider(pkgs []*packages.Package) *provider {
	known := make(map[string]*varToken)

	known[HREQ.signature] = HREQ
	known[HWRI.signature] = HWRI

	return &provider{
		pkgs:  pkgs,
		known: known,
	}
}

func (p provider) varFor(vt *varToken) (*varToken, bool) {
	v, ok := p.known[vt.signature]
	return v, ok
}

func (p provider) provide(vt *varToken) *funcToken {
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

			args := []*varToken{}

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

					rets := []*varToken{}
					for j := 0; j < sig.Results().Len(); j++ {
						ret := sig.Results().At(j)
						token := newVarToken("", ret.Type().String(), vt.varName())
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
