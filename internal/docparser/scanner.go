package docparser

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"os"
	"strings"

	"github.com/bilal-bhatti/zipline/internal/schema"
	"golang.org/x/tools/go/packages"
)

// find a struct type by a name (used for looking up {models.ErrorResponse} in comments)
func FindStruct(pkgs []*packages.Package, tname string) {
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				if gd, ok := n.(*ast.GenDecl); ok {
					for _, spec := range gd.Specs {
						if specType, ok := spec.(*ast.TypeSpec); ok {
							if _, ok := specType.Type.(*ast.StructType); ok {
								if pkg.Name+"."+specType.Name.Name == tname {
									obj := pkg.TypesInfo.Defs[specType.Name]
									schema, err := schema.FieldNew("--", obj.Type())
									if err != nil {
										fmt.Println("error generating schema", err)
									}

									if err := json.NewEncoder(os.Stdout).Encode(schema); err != nil {
										fmt.Println("error generating schema", err)
									}

									if len(strings.TrimSpace(specType.Doc.Text())) > 0 {
										fmt.Println("docs", specType.Doc.Text())
									} else {
										fmt.Println("docs", gd.Doc.Text())
									}
									return true
								}
							}
						}
					}
				}

				return true
			})
		}
	}
}
