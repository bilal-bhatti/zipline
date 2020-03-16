package internal

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/packages"
)

type typeSpecWithPkg struct {
	pkg      *packages.Package
	typeSpec *ast.TypeSpec
}

type scanner struct {
	pkgs []*packages.Package
}

func (s scanner) scan() (map[string]*typeSpecWithPkg, map[string]*template, []*packet) {
	typeSpecs := make(map[string]*typeSpecWithPkg)
	templates := make(map[string]*template)
	packets := []*packet{}

	for _, pkg := range s.pkgs {
		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				switch dt := decl.(type) {
				case *ast.GenDecl:
					for _, spec := range dt.Specs {
						switch specType := spec.(type) {
						case *ast.TypeSpec:
							if _, ok := specType.Type.(*ast.StructType); ok {
								typeSpec := &typeSpecWithPkg{
									pkg:      pkg,
									typeSpec: specType,
								}
								signature := fmt.Sprintf("%s.%s", pkg.PkgPath, specType.Name.String())
								typeSpecs[signature] = typeSpec
							}
						}
					}
				case *ast.FuncDecl:
					funcD, ok := decl.(*ast.FuncDecl)
					if !ok {
						continue
					}

					if !isZiplineNode(pkg.TypesInfo, funcD) {
						continue
					}

					// receiver is ZiplineTemplate, and no additional?
					if funcD.Recv != nil && len(funcD.Recv.List) == 1 {
						// match ZiplineTemplate as receiver
						field := funcD.Recv.List[0]
						if zid, ok := field.Type.(*ast.Ident); ok {
							if zid.String() == ZiplineTemplate {
								// TODO: func must contain a single return statement,
								// returning a func literal, WriteBuffer this check
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
	}

	return typeSpecs, templates, packets
}
