package internal

import (
	"encoding/json"
	"fmt"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/fatih/structtag"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
)

type swagger struct {
	swag *spec.Swagger
}

func newSwagger() swagger {
	swag := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Paths: &spec.Paths{
				Paths: make(map[string]spec.PathItem),
			},
			Parameters:  make(map[string]spec.Parameter),
			Definitions: make(map[string]spec.Schema),
		},
	}

	return swagger{
		swag: swag,
	}
}

func (s swagger) generate(packets []*packet) {
	for _, packet := range packets {
		for _, b := range packet.bindings {
			var pi spec.PathItem
			pi, found := s.swag.Paths.Paths[b.path]
			if !found {
				pi = spec.PathItem{}
			}

			op := &spec.Operation{
				OperationProps: spec.OperationProps{
					Description: "controler",
					Produces:    []string{"application/json"},
					Consumes:    []string{"application/json"},
				},
			}

			for _, param := range b.handler.params {
				if param.varType.Type().String() == "context.Context" {
					continue
				}

				// TODO: do better!
				if isPathParam(b, param) {
					op.AddParam(&spec.Parameter{
						ParamProps: spec.ParamProps{
							Name:     param.varName(),
							In:       "path",
							Required: true,
						},
						SimpleSchema: spec.SimpleSchema{
							Type: typeMap[param.signature],
						},
					})
				} else {
					skema := field("--", param.varType.Type())

					ref := spec.Schema{
						SchemaProps: spec.SchemaProps{
							Ref: spec.MustCreateRef("#/definitions/" + param.shortSignature()),
						},
					}

					s.swag.Definitions[param.shortSignature()] = skema

					op.AddParam(&spec.Parameter{
						ParamProps: spec.ParamProps{
							Name:     "body",
							In:       "body",
							Required: true,
							Schema:   &ref,
						},
					})
				}
			}

			for _, ret := range b.handler.returns {
				if ret.varType.Type().String() == "error" {
					continue
				}

				skema := field("--", ret.varType.Type())
				ref := spec.Schema{
					SchemaProps: spec.SchemaProps{
						Ref: spec.MustCreateRef("#/definitions/" + ret.shortSignature()),
					},
				}

				s.swag.Definitions[ret.shortSignature()] = skema

				op.Responses = &spec.Responses{
					ResponsesProps: spec.ResponsesProps{
						StatusCodeResponses: map[int]spec.Response{
							200: spec.Response{
								ResponseProps: spec.ResponseProps{
									Schema: &ref,
								},
							},
							// TODO: handle error responses, maybe?
						},
					},
				}
			}
			// TODO: if len(b.handler.returns) == 0
			// ResponseCode should be NoContent

			switch strings.ToUpper(b.template) {
			case "GET":
				pi.Get = op
			case "POST":
				pi.Post = op
			case "DELETE":
				pi.Delete = op
			case "PATCH":
				pi.Patch = op
			case "PUT":
				pi.Put = op
			default:
				panic("Oh noos what's this: " + b.template)
			}

			s.swag.Paths.Paths[b.path] = pi
		}
	}

	bites, err := json.MarshalIndent(s.swag, "", "  ")
	if err != nil {
		log.Fatal("Failed to generate json", err)
	}

	err = ioutil.WriteFile(OpenAPIFile, bites, 0644)
	if err != nil {
		panic(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(errors.Wrap(err, "Failed to get working directory"))
	}

	log.Println("Wrote OpenAPI spec to", path.Join(wd, OpenAPIFile))
}

func field(name string, t types.Type) spec.Schema {

	switch tt := t.(type) {
	case *types.Pointer:
		return field(name, tt.Elem())
	case *types.Struct:
		skema := skema("object")
		for i := 0; i < tt.NumFields(); i++ {
			f := tt.Field(i)

			var tn string
			if tt.Tag(i) != "" {
				tags, err := structtag.Parse(tt.Tag(i))
				if err != nil {
					panic(err)
				}

				jsonTag, err := tags.Get("json")
				if err != nil {
					panic(err)
				}

				trace("name", jsonTag.Name)
				tn = jsonTag.Name
			}

			// TODO: handle recursive refs
			// what to do when Account struct has an Account field?
			if f.Name() != name && tn != "-" {
				fs := field(f.Name(), f.Type())
				if f.Embedded() && tn == "" {
					for k, v := range fs.Properties {
						skema.Properties[k] = v
					}
				} else {
					var named = tn
					if tn == "" {
						named = f.Name()
					}

					skema.Properties[named] = fs
				}
			}
		}
		return skema
	case *types.Named:
		return field(name, tt.Underlying())
	case *types.Slice:
		arr := skema("array")
		sk := field(name, tt.Elem())
		arr.Items = &spec.SchemaOrArray{
			Schema: &sk,
		}
		return arr
	case *types.Basic:
		return skema(tt.String())
	default:
		panic(fmt.Sprint("Oh noos unknown type: ", reflect.TypeOf(tt)))
	}
}

func skema(t string) spec.Schema {
	switch t {
	case "object":
		return spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type:       spec.StringOrArray{"object"},
				Properties: make(map[string]spec.Schema),
			},
		}
	case "int", "int64", "uint", "uint8", "uint64":
		return spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{typeMap[t]},
			},
		}
	case "float64":
		return spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{typeMap[t]},
			},
		}
	case "string", "array":
		return spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{t},
			},
		}
	case "bool":
		return spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{typeMap[t]},
			},
		}
	default:
		panic("Oh noos, what's this: " + t)
	}
}

var typeMap = map[string]string{
	"int":     "integer",
	"int8":    "integer",
	"int32":   "integer",
	"int64":   "integer",
	"uint":    "integer",
	"uint8":   "integer",
	"uint32":  "integer",
	"uint64":  "integer",
	"float64": "number",
	"bool":    "boolean",
}
