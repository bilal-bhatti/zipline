package internal

import (
	"go/ast"

	"golang.org/x/tools/go/packages"
)

type scanner struct {
	pkgs []*packages.Package
}

func (s scanner) scan() (map[string]*template, []*packet) {
	templates := make(map[string]*template)
	packets := []*packet{}

	for _, pkg := range s.pkgs {
		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				funcD, ok := decl.(*ast.FuncDecl)
				if !ok {
					continue
				}

				if !isZiplineNode(pkg.TypesInfo, funcD) {
					continue
				}

				if funcD.Recv != nil && len(funcD.Recv.List) == 1 {
					// match ZiplineTemplate as receiver
					field := funcD.Recv.List[0]
					if zid, ok := field.Type.(*ast.Ident); ok {
						if zid.String() == ZiplineTemplate {
							// TODO: func must contain a single return statement,
							// returning a func literal, add this check
							templates[funcD.Name.String()] = newTemplate(funcD)

							// this is a ZiplineTemplate method, used as a template
							// we know what this is, so no further inspection
							continue
						}
					}
				}

				// else this is a func that declares zipline bindings
				packets = append(packets, &packet{
					pkg:      pkg,
					funcDecl: funcD,
				})
			}
		}
	}

	return templates, packets
}
