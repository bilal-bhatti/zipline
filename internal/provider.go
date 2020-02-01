package internal

import (
	"go/types"

	"golang.org/x/tools/go/packages"
)

type provider struct {
	pkgs  []*packages.Package
	cache map[string]*types.Func
}

func newProvider(pkgs []*packages.Package) *provider {
	return &provider{
		pkgs:  pkgs,
		cache: make(map[string]*types.Func),
	}
}

func (p provider) provide(vt *varToken) *funcToken {
	if pf, ok := p.cache[vt.signature]; ok {
		return &funcToken{cpkg: pf.Pkg().Name(), signature: pf.FullName()}
	}

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

			// match signature --- should be improved
			// look for func that takes an http.Request type
			// and returns the desired type
			// TODO: handle error return
			if sig.Params().Len() == 1 && sig.Results().Len() == 1 {
				if vt.signature == sig.Results().At(0).Type().String() && hreq.signature == sig.Params().At(0).Type().String() {
					// cache
					p.cache[vt.signature] = pf

					return &funcToken{cpkg: pf.Pkg().Name(), signature: pf.FullName()}
				}
			}
		}
	}

	return nil
}
