package schema

import (
	"fmt"
	"go/types"
	"reflect"

	"github.com/bilal-bhatti/zipline/internal/tokens"
	"github.com/fatih/structtag"
	"github.com/getkin/kin-openapi/openapi3"
)

func FieldNew(name string, t types.Type) (*openapi3.Schema, error) {
	switch tt := t.(type) {
	case *types.Pointer:
		return FieldNew(name, tt.Elem())
	case *types.Struct:
		schema := openapi3.NewObjectSchema().WithoutAdditionalProperties()

		for i := 0; i < tt.NumFields(); i++ {
			f := tt.Field(i)
			if !f.Exported() {
				continue
			}

			fieldTypeToken := tokens.NewTypeToken(f.Type().String(), f.Name())

			var jsonName string
			var tag = tt.Tag(i)
			var tags *structtag.Tags
			if tag != "" {
				tags, err := structtag.Parse(tag)
				if err != nil {
					return nil, err
				}

				jsonTag, err := tags.Get("json")
				if err != nil {
					return nil, err
				}
				jsonName = jsonTag.Name
			}

			// TODO: handle recursive refs
			// what to do when Account struct has an Account field?
			if f.Name() != name && jsonName != "-" {
				var fs *openapi3.Schema
				var err error

				if fieldTypeToken.Signature == "time.Time" {
					fs = openapi3.NewDateTimeSchema()
				} else {
					fs, err = FieldNew(f.Name(), f.Type())
					if err != nil {
						return nil, err
					}
				}

				if tags != nil {
					timeFmt, err := tags.Get("format")
					if err == nil {
						fs.Format = timeFmt.Value()
					}
				}

				if f.Embedded() && jsonName == "" {
					for k, v := range fs.Properties {
						// schema.Properties[k] = v
						schema.WithPropertyRef(k, v)
					}
				} else {
					var named = jsonName
					if jsonName == "" {
						named = f.Name()
					}

					schema.Properties[named] = fs.NewRef()
				}
			}
		}
		return schema, nil
	case *types.Named:
		return FieldNew(name, tt.Underlying())
	case *types.Slice:
		obj, err := FieldNew(name, tt.Elem())
		if err != nil {
			return nil, err
		}

		return openapi3.NewArraySchema().WithItems(obj), nil
	case *types.Basic:
		return SchemaNew(tt.String())
	case *types.Map:
		return openapi3.NewObjectSchema(), nil
	default:
		return nil, fmt.Errorf("oh noos unknown type: %v", reflect.TypeOf(tt))
	}
}

func SchemaNew(t string) (*openapi3.Schema, error) {
	switch t {
	case "int", "int8", "int16", "uint", "uint8", "uint16":
		return openapi3.NewIntegerSchema(), nil
	case "int32", "uint32":
		return openapi3.NewInt32Schema(), nil
	case "int64", "uint64":
		return openapi3.NewInt64Schema(), nil
	case "float32", "float64":
		return openapi3.NewFloat64Schema(), nil
	case "string":
		return openapi3.NewStringSchema(), nil
	case "bool":
		return openapi3.NewBoolSchema(), nil
	case "time.Time":
		return openapi3.NewDateTimeSchema(), nil
	case "byte":
		return openapi3.NewBytesSchema(), nil
	default:
		return nil, fmt.Errorf("oh noos, what's this: %s", t)
	}
}
