package internal

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"os"
	"strings"

	"github.com/bilal-bhatti/zipline/internal/schema"
	"golang.org/x/tools/go/packages"
)

type typeSpecWithPkg struct {
	pkg      *packages.Package
	typeSpec *ast.TypeSpec
	docs     string
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

								signature := fmt.Sprintf("%s.%s", pkg.PkgPath, specType.Name.String())
								var docstring string
								if specType.Doc == nil {
									docstring = strings.TrimSpace(dt.Doc.Text())
								} else {
									docstring = strings.TrimSpace(specType.Doc.Text())
								}
								typeSpec := &typeSpecWithPkg{
									pkg:      pkg,
									typeSpec: specType,
									docs:     docstring,
								}
								typeSpecs[signature] = typeSpec
							}
						}
					}
				case *ast.FuncDecl:
					if !isZiplineNode(pkg.TypesInfo, dt) {
						continue
					}

					// check if method receiver is ZiplineTemplate
					if dt.Recv != nil && len(dt.Recv.List) > 0 {
						if ok := isZiplineNode(pkg.TypesInfo, dt.Recv); ok {
							templates[dt.Name.String()] = newTemplate(dt)

							// this is a ZiplineTemplate method, used as a template
							// we know what this is, so no further inspection
							continue
						}
					}

					// else this is a func that declares zipline bindings
					packets = append(packets, &packet{
						pkg:      pkg,
						funcDecl: dt,
					})
				}
			}
		}
	}

	s.find("models.ErrorResponse")

	return typeSpecs, templates, packets
}

func (s scanner) find(tname string) {
	for _, pkg := range s.pkgs {
		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				if ts, ok := n.(*ast.TypeSpec); ok {
					if struc, ok := ts.Type.(*ast.StructType); ok {
						if pkg.Name+"."+ts.Name.Name == tname {
							obj := pkg.TypesInfo.Defs[ts.Name]
							fmt.Println("obj types.Type", obj.Type(), ts.Doc.Text())
							fmt.Println("struct", struc.Fields)

							// schema, err := schema.Field("--", obj.Type())
							schema, err := schema.FieldNew("--", obj.Type())
							if err != nil {
								fmt.Println("error generating schema", err)
							}

							if err := json.NewEncoder(os.Stdout).Encode(schema); err != nil {
								fmt.Println("error generating schema", err)
							}
							return false
						}
					}
				}
				return true
			})
		}
	}
}
