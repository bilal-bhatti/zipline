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
				if funcD, ok := decl.(*ast.FuncDecl); ok {
					if isBindingSpecNode(pkg.TypesInfo, funcD) {
						packets = append(packets, &packet{
							pkg:      pkg,
							bindings: funcD,
						})
					} else if funcD.Recv != nil && len(funcD.Recv.List) == 1 {
						// match ZiplineTemplate as receiver
						field := funcD.Recv.List[0]
						if zt, ok := field.Type.(*ast.Ident); ok {
							if zt.String() == ZiplineTemplate {
								// TODO: func must contain a single return statement,
								// returning a func literal
								templates[funcD.Name.String()] = newTemplate(funcD)
							}
						}
					}
				}
			}
		}
	}

	return templates, packets
}
