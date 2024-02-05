package schema

import (
	"fmt"
	"go/types"
	"log"
	"reflect"

	"github.com/bilal-bhatti/zipline/internal/tokens"
	"github.com/fatih/structtag"
	"github.com/go-openapi/spec"
)

func Field(name string, t types.Type) (*spec.Schema, error) {
	switch tt := t.(type) {
	case *types.Pointer:
		return Field(name, tt.Elem())
	case *types.Struct:
		skema, err := Skema("object")
		if err != nil {
			return nil, err
		}

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
				tags, err = structtag.Parse(tag)
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
				var fs *spec.Schema

				if fieldTypeToken.Signature == "time.Time" {
					fs = spec.DateTimeProperty()
				} else {
					fs, err = Field(f.Name(), f.Type())
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
						skema.Properties[k] = v
					}
				} else {
					var named = jsonName
					if jsonName == "" {
						named = f.Name()
					}

					skema.Properties[named] = *fs
				}
			}
		}
		return skema, nil
	case *types.Named:
		return Field(name, tt.Underlying())
	case *types.Slice:
		obj, err := Field(name, tt.Elem())
		if err != nil {
			return nil, err
		}
		return spec.ArrayProperty(obj), nil
	case *types.Basic:
		return Skema(tt.String())
	case *types.Map:
		return &spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type:                 spec.StringOrArray{"object"},
				AdditionalProperties: &spec.SchemaOrBool{Allows: true},
			},
		}, nil
	default:
		return nil, fmt.Errorf("oh noos unknown type: %v", reflect.TypeOf(tt))
	}
}

func Skema(t string) (*spec.Schema, error) {
	switch t {
	case "object":
		return &spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type:                 spec.StringOrArray{"object"},
				Properties:           make(map[string]spec.Schema),
				AdditionalProperties: &spec.SchemaOrBool{Allows: false},
			},
		}, nil
	case "int", "uint":
		return &spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"integer"},
			},
		}, nil
	case "int8", "uint8":
		return spec.Int8Property(), nil
	case "int16", "uint16":
		return spec.Int16Property(), nil
	case "int32", "uint32":
		return spec.Int32Property(), nil
	case "int64", "uint64":
		return spec.Int64Property(), nil
	case "float32":
		return spec.Float32Property(), nil
	case "float64":
		return spec.Float64Property(), nil
	case "string":
		return spec.StringProperty(), nil
	case "bool":
		return spec.BoolProperty(), nil
	case "time.Time":
		return spec.DateTimeProperty(), nil
	case "byte":
		// TODO: Swagger 2.0 doesn't support binary data type
		// map byte to string, but it should really be removed
		// fail at this point? maybe a strict mode vs lax mode?
		log.Println("Swagger 2.0 doesn't support byte type, generated spec will not be valid")
		return &spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"byte"},
			},
		}, nil
	default:
		return nil, fmt.Errorf("oh noos, what's this: %s", t)
	}
}
