package docparser

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/bilal-bhatti/zipline/internal/schema"
	"github.com/bilal-bhatti/zipline/internal/util"
	"golang.org/x/tools/go/packages"
)

// find a struct type by a name (used for looking up {models.ErrorResponse} in comments)
func FindStruct(pkgs []*packages.Package, tname string) (map[string]interface{}, error) {
	var obj types.Object
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				if gd, ok := n.(*ast.GenDecl); ok {
					for _, spec := range gd.Specs {
						if specType, ok := spec.(*ast.TypeSpec); ok {
							if _, ok := specType.Type.(*ast.StructType); ok {
								if pkg.Name+"."+specType.Name.Name == tname {
									obj = pkg.TypesInfo.Defs[specType.Name]
									// fmt.Println("found object ", obj)
									// schema, err := schema.FieldNew("--", obj.Type())
									// if err != nil {
									// 	fmt.Println("error generating schema", err)
									// }

									// if err := json.NewEncoder(os.Stdout).Encode(schema); err != nil {
									// 	fmt.Println("error generating schema", err)
									// }

									// if len(strings.TrimSpace(specType.Doc.Text())) > 0 {
									// 	fmt.Println("docs", specType.Doc.Text())
									// } else {
									// 	fmt.Println("docs", gd.Doc.Text())
									// }
									return false
								}
							}
						}
					}
				}

				return true
			})
		}
	}

	if obj == nil {
		return nil, fmt.Errorf("struct not found: %s", tname)
	}

	sch, err := schema.FieldNew("--", obj.Type())
	if err != nil {
		return nil, err
	}

	smap, err := util.StructToMap(sch)
	if err != nil {
		return nil, err
	}

	return smap, nil
}
